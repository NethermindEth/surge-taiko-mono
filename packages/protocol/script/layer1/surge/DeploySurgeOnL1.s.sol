// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/utils/Strings.sol";
import "@risc0/contracts/groth16/RiscZeroGroth16Verifier.sol";
import { SP1Verifier as SuccinctVerifier } from "@sp1-contracts/src/v4.0.0-rc.3/SP1VerifierPlonk.sol";

// Actually this one is deployed already on mainnet, but we are now deploying our own (non via-ir)
// version. For mainnet, it is easier to go with one of:
// - https://github.com/daimo-eth/p256-verifier
// - https://github.com/rdubois-crypto/FreshCryptoLib
import "@p256-verifier/contracts/P256Verifier.sol";

import "src/shared/common/LibStrings.sol";
import "src/shared/bridge/Bridge.sol";
import "src/shared/signal/SignalService.sol";
import "src/shared/tokenvault/BridgedERC1155.sol";
import "src/shared/tokenvault/BridgedERC20.sol";
import "src/shared/tokenvault/BridgedERC721.sol";
import "src/shared/tokenvault/ERC1155Vault.sol";
import "src/shared/tokenvault/ERC20Vault.sol";
import "src/shared/tokenvault/ERC721Vault.sol";
import "src/layer1/surge/SurgeTimelockController.sol";
import "src/layer1/automata-attestation/AutomataDcapV3Attestation.sol";
import "src/layer1/automata-attestation/lib/PEMCertChainLib.sol";
import "src/layer1/automata-attestation/utils/SigVerifyLib.sol";
import "src/layer1/based/TaikoL1.sol";
import "src/layer1/surge/SurgeTierRouter.sol";
import "src/layer1/surge/SurgeTaikoL1.sol";
import "src/layer1/verifiers/SgxVerifier.sol";
import "src/layer1/verifiers/Risc0Verifier.sol";
import "src/layer1/verifiers/SP1Verifier.sol";
import "src/layer1/verifiers/compose/TwoOfThreeVerifier.sol";
import "src/layer1/verifiers/compose/ComposeVerifier.sol";
import "test/shared/token/FreeMintERC20.sol";
import "test/shared/token/MayFailFreeMintERC20.sol";
import "test/shared/DeployCapability.sol";

/// @title DeploySurgeOnL1
/// @notice This script deploys Taiko protocol modified for Nethermind's Surge.
/// @dev Modification include:
/// - Using Ether as bonds, thus no reference to Taiko token and prover set is removed
/// - Configurable chainId to spin up Surge networks with ease
/// - Automated setup of fields in L1 address manager
/// - Removed guardian prover
contract DeploySurgeOnL1 is DeployCapability {
    uint256 internal immutable ADDRESS_LENGTH = 40;

    uint64 internal l2ChainId = uint64(vm.envUint("L2_CHAINID"));

    modifier broadcast() {
        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        require(privateKey != 0, "invalid priv key");
        vm.startBroadcast();
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast {
        require(l2ChainId != block.chainid || l2ChainId != 0, "L2_CHAIN_ID");
        require(vm.envBytes32("L2_GENESIS_HASH") != 0, "L2_GENESIS_HASH");

        //---------------------------------------------------------------
        // Timelocked owner setup
        address[] memory executors = vm.envAddress("OWNER_MULTISIG_SIGNERS", ",");

        address ownerMultisig = vm.envAddress("OWNER_MULTISIG");
        addressNotNull(ownerMultisig, "ownerMultisig");

        address[] memory proposers = new address[](1);
        proposers[0] = ownerMultisig;

        uint256 timelockPeriod = uint64(vm.envUint("TIMELOCK_PERIOD"));
        address timelockController = address(
            new SurgeTimelockedController(timelockPeriod, proposers, executors, address(0))
        );
        address contractOwner = timelockController;

        console2.log("contractOwner(timelocked): ", contractOwner);

        // ---------------------------------------------------------------
        // Deploy shared contracts
        (address sharedAddressManager) = deploySharedContracts(contractOwner);
        console2.log("sharedAddressManager: ", sharedAddressManager);

        // ---------------------------------------------------------------
        // Deploy rollup contracts
        address verifierOwner = vm.envAddress("VERIFIER_OWNER");
        addressNotNull(verifierOwner, "verifierOwner");

        address rollupAddressManager = deployRollupContracts(sharedAddressManager, contractOwner, verifierOwner);

        // ---------------------------------------------------------------
        // Signal service need to authorize the new rollup
        address signalServiceAddr = AddressManager(sharedAddressManager).getAddress(
            uint64(block.chainid), LibStrings.B_SIGNAL_SERVICE
        );
        addressNotNull(signalServiceAddr, "signalServiceAddr");
        SignalService signalService = SignalService(signalServiceAddr);

        address taikoL1Addr = AddressManager(rollupAddressManager).getAddress(
            uint64(block.chainid), LibStrings.B_TAIKO
        );
        addressNotNull(taikoL1Addr, "taikoL1Addr");

        SurgeTimelockedController(payable(timelockController)).init(taikoL1Addr);

        SignalService(signalServiceAddr).authorize(taikoL1Addr, true);

        console2.log("------------------------------------------");
        console2.log("msg.sender: ", msg.sender);
        console2.log("address(this): ", address(this));
        console2.log("signalService.owner(): ", signalService.owner());
        console2.log("------------------------------------------");

        if (signalService.owner() == msg.sender) {
            signalService.transferOwnership(contractOwner);
        } else {
            console2.log("------------------------------------------");
            console2.log("Warning - you need to transact manually:");
            console2.log("signalService.authorize(taikoL1Addr, bytes32(block.chainid))");
            console2.log("- signalService : ", signalServiceAddr);
            console2.log("- taikoL1Addr   : ", taikoL1Addr);
            console2.log("- chainId       : ", block.chainid);
        }

        address taikoL2Address = getConstantAddress(vm.toString(l2ChainId), "10001");
        address l2SignalServiceAddress = getConstantAddress(vm.toString(l2ChainId), "5");
        address l2BridgeAddress = getConstantAddress(vm.toString(l2ChainId), "1");
        address l2Erc20VaultAddress = getConstantAddress(vm.toString(l2ChainId), "2");

        // ---------------------------------------------------------------
        // Register L2 addresses
        register(rollupAddressManager, "taiko", taikoL2Address, l2ChainId);
        register(
            rollupAddressManager, "signal_service", l2SignalServiceAddress, l2ChainId
        );
        register(sharedAddressManager, "signal_service", l2SignalServiceAddress, l2ChainId);
        register(sharedAddressManager, "bridge", l2BridgeAddress, l2ChainId);
        register(sharedAddressManager, "erc20_vault", l2Erc20VaultAddress, l2ChainId);


        // ---------------------------------------------------------------
        // Deploy other contracts
        if (block.chainid != 1) {
            deployAuxContracts();
        }

        if (AddressManager(sharedAddressManager).owner() == msg.sender) {
            AddressManager(sharedAddressManager).transferOwnership(contractOwner);
            console2.log("** sharedAddressManager ownership transferred to:", contractOwner);
        }

        AddressManager(rollupAddressManager).transferOwnership(contractOwner);
        console2.log("** rollupAddressManager ownership transferred to:", contractOwner);
    }

    function deploySharedContracts(address owner) internal returns (address sharedAddressManager) {
        addressNotNull(owner, "owner");

        sharedAddressManager = deployProxy({
            name: "shared_address_manager",
            impl: address(new AddressManager()),
            data: abi.encodeCall(AddressManager.init, (address(0)))
        });

        // Deploy Bridging contracts
        deployProxy({
            name: "signal_service",
            impl: address(new SignalService()),
            data: abi.encodeCall(SignalService.init, (address(0), sharedAddressManager)),
            registerTo: sharedAddressManager
        });

        address brdige = deployProxy({
            name: "bridge",
            impl: address(new Bridge()),
            data: abi.encodeCall(Bridge.init, (address(0), sharedAddressManager)),
            registerTo: sharedAddressManager
        });

        Bridge(payable(brdige)).transferOwnership(owner);

        console2.log("------------------------------------------");
        console2.log(
            "Warning - you need to register *all* counterparty bridges to enable multi-hop bridging:"
        );
        console2.log(
            "sharedAddressManager.setAddress(remoteChainId, \"bridge\", address(remoteBridge))"
        );
        console2.log("- sharedAddressManager : ", sharedAddressManager);

        // Deploy Vaults
        deployProxy({
            name: "erc20_vault",
            impl: address(new ERC20Vault()),
            data: abi.encodeCall(ERC20Vault.init, (owner, sharedAddressManager)),
            registerTo: sharedAddressManager
        });
        deployProxy({
            name: "erc721_vault",
            impl: address(new ERC721Vault()),
            data: abi.encodeCall(ERC721Vault.init, (owner, sharedAddressManager)),
            registerTo: sharedAddressManager
        });
        deployProxy({
            name: "erc1155_vault",
            impl: address(new ERC1155Vault()),
            data: abi.encodeCall(ERC1155Vault.init, (owner, sharedAddressManager)),
            registerTo: sharedAddressManager
        });

        console2.log("------------------------------------------");
        console2.log(
            "Warning - you need to register *all* counterparty vaults to enable multi-hop bridging:"
        );
        console2.log(
            "sharedAddressManager.setAddress(remoteChainId, \"erc20_vault\", address(remoteERC20Vault))"
        );
        console2.log(
            "sharedAddressManager.setAddress(remoteChainId, \"erc721_vault\", address(remoteERC721Vault))"
        );
        console2.log(
            "sharedAddressManager.setAddress(remoteChainId, \"erc1155_vault\", address(remoteERC1155Vault))"
        );
        console2.log("- sharedAddressManager : ", sharedAddressManager);

        // Deploy Bridged token implementations
        register(sharedAddressManager, "bridged_erc20", address(new BridgedERC20()));
        register(sharedAddressManager, "bridged_erc721", address(new BridgedERC721()));
        register(sharedAddressManager, "bridged_erc1155", address(new BridgedERC1155()));
    }

    function deployRollupContracts(
        address _sharedAddressManager,
        address owner,
        address verifierOwner
    )
        internal
        returns (address rollupAddressManager)
    {
        addressNotNull(_sharedAddressManager, "sharedAddressManager");
        addressNotNull(owner, "owner");

        rollupAddressManager = deployProxy({
            name: "rollup_address_manager",
            impl: address(new AddressManager()),
            data: abi.encodeCall(AddressManager.init, (address(0)))
        });

        // ---------------------------------------------------------------
        // Register shared contracts in the new rollup
        copyRegister(rollupAddressManager, _sharedAddressManager, "signal_service");
        copyRegister(rollupAddressManager, _sharedAddressManager, "bridge");

        TaikoL1 taikoL1 = TaikoL1(address(new SurgeTaikoL1(l2ChainId)));

        deployProxy({
            name: "taiko",
            impl: address(taikoL1),
            data: abi.encodeCall(
                TaikoL1.init,
                (
                    owner,
                    rollupAddressManager,
                    vm.envBytes32("L2_GENESIS_HASH"),
                    false
                )
            ),
            registerTo: rollupAddressManager
        });

        deployProxy({
            name: "tier_sgx",
            impl: address(new SgxVerifier()),
            data: abi.encodeCall(SgxVerifier.init, (verifierOwner, rollupAddressManager)),
            registerTo: rollupAddressManager
        });

        register(
            rollupAddressManager,
            "tier_router",
            address(new SurgeTierRouter())
        );

        // No need to proxy these, because they are 3rd party. If we want to modify, we simply
        // change the registerAddress("automata_dcap_attestation", address(attestation));
        P256Verifier p256Verifier = new P256Verifier();
        SigVerifyLib sigVerifyLib = new SigVerifyLib(address(p256Verifier));
        PEMCertChainLib pemCertChainLib = new PEMCertChainLib();

        addAddress({
            name: "PemCertChainLib",
            proxy: address(pemCertChainLib)
        });

        address automateDcapV3AttestationImpl = address(new AutomataDcapV3Attestation());

        address automataProxy = deployProxy({
            name: "automata_dcap_attestation",
            impl: automateDcapV3AttestationImpl,
            data: abi.encodeCall(
                AutomataDcapV3Attestation.init, (verifierOwner, address(sigVerifyLib), address(pemCertChainLib))
            ),
            registerTo: rollupAddressManager
        });

        // Log addresses for the user to register sgx instance
        console2.log("SigVerifyLib", address(sigVerifyLib));
        console2.log("PemCertChainLib", address(pemCertChainLib));
        console2.log("AutomataDcapVaAttestation", automataProxy);

        deployZKVerifiers(verifierOwner, rollupAddressManager);

        // Deploy composite verifier
        deployProxy({
            name: "tier_two_of_three",
            impl: address(new TwoOfThreeVerifier()),
            data: abi.encodeCall(ComposeVerifier.init, (verifierOwner, rollupAddressManager)),
            registerTo: rollupAddressManager
        });
    }

    // deploy both sp1 & risc0 verifiers.
    // using function to avoid stack too deep error
    function deployZKVerifiers(address verifierOwner, address rollupAddressManager) private {
        // Deploy r0 groth16 verifier
        RiscZeroGroth16Verifier verifier =
            new RiscZeroGroth16Verifier(ControlID.CONTROL_ROOT, ControlID.BN254_CONTROL_ID);
        register(rollupAddressManager, "risc0_groth16_verifier", address(verifier));

        deployProxy({
            name: "tier_zkvm_risc0",
            impl: address(new Risc0Verifier()),
            data: abi.encodeCall(Risc0Verifier.init, (verifierOwner, rollupAddressManager)),
            registerTo: rollupAddressManager
        });

        // Deploy sp1 plonk verifier
        SuccinctVerifier succinctVerifier = new SuccinctVerifier();
        register(rollupAddressManager, "sp1_remote_verifier", address(succinctVerifier));

        deployProxy({
            name: "tier_zkvm_sp1",
            impl: address(new SP1Verifier()),
            data: abi.encodeCall(SP1Verifier.init, (verifierOwner, rollupAddressManager)),
            registerTo: rollupAddressManager
        });
    }

    function deployAuxContracts() private {
        address horseToken = address(new FreeMintERC20("Horse Token", "HORSE"));
        console2.log("HorseToken", horseToken);

        address bullToken = address(new MayFailFreeMintERC20("Bull Token", "BULL"));
        console2.log("BullToken", bullToken);
    }

    function addressNotNull(address addr, string memory err) private pure {
        require(addr != address(0), err);
    }

    function getConstantAddress(string memory prefix, string memory suffix) internal pure returns (address) {
        bytes memory prefixBytes = bytes(prefix);
        bytes memory suffixBytes = bytes(suffix);

        require(prefixBytes.length + suffixBytes.length <= ADDRESS_LENGTH, "Prefix + suffix too long");

        // Create the middle padding of zeros
        uint256 paddingLength = ADDRESS_LENGTH - prefixBytes.length - suffixBytes.length;
        bytes memory padding = new bytes(paddingLength);
        for(uint i = 0; i < paddingLength; i++) {
            padding[i] = "0";
        }

        // Concatenate the parts
        string memory hexString = string(
            abi.encodePacked(
                "0x",
                prefix,
                string(padding),
                suffix
            )
        );

        return vm.parseAddress(hexString);
    }
}