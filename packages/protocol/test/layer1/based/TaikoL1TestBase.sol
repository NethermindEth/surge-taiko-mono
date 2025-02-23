// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../TaikoL1Test.sol";

abstract contract TaikoL1TestBase is TaikoTest {
    AddressManager public addressManager;
    SignalService public ss;
    TaikoL1 public L1;
    TaikoData.Config conf;
    uint256 internal logCount;
    Risc0Verifier public rv;
    SP1Verifier public sp1;
    SgxVerifier public sv;
    GuardianProver public gp;
    TestTierProvider public cp;
    Bridge public bridge;

    bytes32 public GENESIS_BLOCK_HASH = keccak256("GENESIS_BLOCK_HASH");

    address public L2SS = randAddress();
    address public L2 = randAddress();
    // Bootstrapped SGX instances (by owner)
    address internal SGX_X_0 = vm.addr(0x1000004);
    address internal SGX_X_1 = vm.addr(0x1000005);

    function deployTaikoL1() internal virtual returns (TaikoL1 taikoL1);

    function tierProvider() internal view returns (ITierProvider) {
        ITierRouter tierRouter = ITierRouter(L1.resolve(LibStrings.B_TIER_ROUTER, false));
        return ITierProvider(tierRouter.getProvider(0));
    }

    function setUp() public virtual {
        L1 = deployTaikoL1();
        conf = L1.getConfig();

        addressManager = AddressManager(
            deployProxy({
                name: "address_manager",
                impl: address(new AddressManager()),
                data: abi.encodeCall(AddressManager.init, (address(0)))
            })
        );

        ss = SignalService(
            deployProxy({
                name: "signal_service",
                impl: address(new SignalService()),
                data: abi.encodeCall(SignalService.init, (address(0), address(addressManager)))
            })
        );
        ss.authorize(address(L1), true);

        sv = SgxVerifier(
            deployProxy({
                name: "tier_sgx",
                impl: address(new SgxVerifier()),
                data: abi.encodeCall(SgxVerifier.init, (address(0), address(addressManager)))
            })
        );

        address[] memory initSgxInstances = new address[](1);
        initSgxInstances[0] = SGX_X_0;
        sv.addInstances(initSgxInstances);

        gp = GuardianProver(
            deployProxy({
                name: "guardian_prover",
                impl: address(new GuardianProver()),
                data: abi.encodeCall(GuardianProver.init, (address(0), address(addressManager)))
            })
        );

        setupGuardianProverMultisig();

        cp = new TestTierProvider();

        bridge = Bridge(
            payable(
                deployProxy({
                    name: "bridge",
                    impl: address(new Bridge()),
                    data: abi.encodeCall(Bridge.init, (address(0), address(addressManager))),
                    registerTo: address(addressManager)
                })
            )
        );

        registerAddress("taiko", address(L1));
        registerAddress("tier_sgx", address(sv));
        registerAddress("tier_guardian", address(gp));
        registerAddress("tier_router", address(cp));
        registerAddress("signal_service", address(ss));
        registerL2Address("taiko", address(L2));
        registerL2Address("signal_service", address(L2SS));
        registerL2Address("taiko_l2", address(L2));

        L1.init(address(0), address(addressManager), GENESIS_BLOCK_HASH, false);

        printVariables("init  ");
    }

    function proposeBlock(
        address proposer,
        uint24 txListSize
    )
        internal
        returns (TaikoData.BlockMetadata memory meta, TaikoData.EthDeposit[] memory ethDeposits)
    {
        // For the test not to fail, set the message.value to the highest, the
        // rest will be returned
        // anyways
        uint256 msgValue = 2 ether;

        TaikoData.HookCall[] memory hookcalls = new TaikoData.HookCall[](0);
        vm.prank(proposer, proposer);
        (meta, ethDeposits) = L1.proposeBlock{ value: msgValue }(
            abi.encode(TaikoData.BlockParams(address(0), address(0), 0, 0, hookcalls, "")),
            new bytes(txListSize)
        );
    }

    function proveBlock(
        address prover,
        TaikoData.BlockMetadata memory meta,
        bytes32 parentHash,
        bytes32 blockHash,
        bytes32 stateRoot,
        uint16 tier,
        bytes4 revertReason
    )
        internal
        virtual
    {
        TaikoData.Transition memory tran = TaikoData.Transition({
            parentHash: parentHash,
            blockHash: blockHash,
            stateRoot: stateRoot,
            graffiti: 0x0
        });

        TaikoData.TierProof memory proof;
        proof.tier = tier;
        address newInstance;

        // Keep changing the pub key associated with an instance to avoid
        // attacks,
        // obviously just a mock due to 2 addresses changing all the time.
        (newInstance,) = sv.instances(0);
        if (newInstance == SGX_X_0) {
            newInstance = SGX_X_1;
        } else {
            newInstance = SGX_X_0;
        }

        if (tier == LibTiers.TIER_SGX) {
            bytes memory signature =
                createSgxSignatureProof(tran, newInstance, prover, keccak256(abi.encode(meta)));

            proof.data = bytes.concat(bytes4(0), bytes20(newInstance), signature);
        }

        if (tier == LibTiers.TIER_GUARDIAN) {
            proof.data = "";

            // Grant 2 signatures, 3rd might be a revert
            vm.prank(David, David);
            gp.approve(meta, tran, proof);
            vm.prank(Emma, Emma);
            gp.approve(meta, tran, proof);

            if (revertReason != "") {
                vm.prank(Frank, Frank);
                vm.expectRevert(); // Revert reason is 'wrapped' so will not be
                    // identical to the expectedRevert
                gp.approve(meta, tran, proof);
            } else {
                vm.prank(Frank, Frank);
                gp.approve(meta, tran, proof);
            }
        } else {
            if (revertReason != "") {
                vm.prank(prover);
                vm.expectRevert(revertReason);
                L1.proveBlock(meta.id, abi.encode(meta, tran, proof));
            } else {
                vm.prank(prover);
                L1.proveBlock(meta.id, abi.encode(meta, tran, proof));
            }
        }
    }

    function verifyBlock(uint64 count) internal {
        L1.verifyBlocks(count);
    }

    function setupGuardianProverMultisig() internal {
        address[] memory initMultiSig = new address[](5);
        initMultiSig[0] = David;
        initMultiSig[1] = Emma;
        initMultiSig[2] = Frank;
        initMultiSig[3] = Grace;
        initMultiSig[4] = Henry;

        gp.setGuardians(initMultiSig, 3, true);
    }

    function registerAddress(bytes32 nameHash, address addr) internal {
        addressManager.setAddress(uint64(block.chainid), nameHash, addr);
        console2.log(block.chainid, uint256(nameHash), unicode"→", addr);
    }

    function registerL2Address(bytes32 nameHash, address addr) internal {
        addressManager.setAddress(conf.chainId, nameHash, addr);
        console2.log(conf.chainId, string(abi.encodePacked(nameHash)), unicode"→", addr);
    }

    function createSgxSignatureProof(
        TaikoData.Transition memory tran,
        address newInstance,
        address prover,
        bytes32 metaHash
    )
        internal
        view
        returns (bytes memory signature)
    {
        uint64 chainId = L1.getConfig().chainId;
        bytes32 digest = LibPublicInput.hashPublicInputs(
            tran, address(sv), newInstance, prover, metaHash, chainId
        );

        uint256 signerPrivateKey;

        // In the test suite these are the 3 which acts as provers
        if (SGX_X_0 == newInstance) {
            signerPrivateKey = 0x1000005;
        } else if (SGX_X_1 == newInstance) {
            signerPrivateKey = 0x1000004;
        }

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPrivateKey, digest);
        signature = abi.encodePacked(r, s, v);
    }

    function giveEthAndDepositBond(address to, uint256 bondEth, uint256 proposalsEth) internal {
        vm.deal(to, bondEth + proposalsEth);
        vm.prank(to);
        L1.depositBond{ value: bondEth }();

        console2.log("ETH balance:", to, to.balance);
    }

    function printVariables(string memory comment) internal view {
        (, TaikoData.SlotB memory b) = L1.getStateVariables();

        string memory str = string.concat(
            "---chain [",
            vm.toString(b.lastVerifiedBlockId),
            unicode"→",
            vm.toString(b.numBlocks),
            "] // ",
            comment
        );
        console2.log(str);
    }

    function mine(uint256 counts) internal {
        vm.warp(block.timestamp + 20 * counts);
        vm.roll(block.number + counts);
    }
}
