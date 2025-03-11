// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/utils/Strings.sol";
import "@risc0/contracts/groth16/RiscZeroGroth16Verifier.sol";
import { SP1Verifier as SuccinctVerifier } from
    "@sp1-contracts/src/v4.0.0-rc.3/SP1VerifierPlonk.sol";
import "@p256-verifier/contracts/P256Verifier.sol";
import "src/shared/common/DefaultResolver.sol";
import "src/shared/libs/LibStrings.sol";
import "src/shared/tokenvault/BridgedERC1155.sol";
import "src/shared/tokenvault/BridgedERC20.sol";
import "src/shared/tokenvault/BridgedERC721.sol";
import "src/layer1/automata-attestation/AutomataDcapV3Attestation.sol";
import "src/layer1/automata-attestation/lib/PEMCertChainLib.sol";
import "src/layer1/automata-attestation/utils/SigVerifyLib.sol";
import "src/layer1/devnet/DevnetInbox.sol";
import "src/layer1/mainnet/MainnetInbox.sol";
import "src/layer1/based/TaikoInbox.sol";
import "src/layer1/fork-router/PacayaForkRouter.sol";
import "src/layer1/forced-inclusion/TaikoWrapper.sol";
import "src/layer1/forced-inclusion/ForcedInclusionStore.sol";
import "src/layer1/mainnet/multirollup/MainnetBridge.sol";
import "src/layer1/mainnet/multirollup/MainnetERC1155Vault.sol";
import "src/layer1/mainnet/multirollup/MainnetERC20Vault.sol";
import "src/layer1/mainnet/multirollup/MainnetERC721Vault.sol";
import "src/layer1/mainnet/multirollup/MainnetSignalService.sol";
import "src/layer1/preconf/impl/PreconfWhitelist.sol";
import "src/layer1/preconf/impl/PreconfRouter.sol";
import "src/layer1/preconf/PreconfInbox.sol";
import "src/layer1/verifiers/Risc0Verifier.sol";
import "src/layer1/verifiers/SP1Verifier.sol";
import "src/layer1/verifiers/SgxVerifier.sol";
import "src/layer1/verifiers/compose/ComposeVerifier.sol";
import "src/layer1/devnet/verifiers/DevnetVerifier.sol";
import "test/shared/helpers/FreeMintERC20Token.sol";
import "test/shared/helpers/FreeMintERC20Token_With50PctgMintAndTransferFailure.sol";
import "test/shared/DeployCapability.sol";

// Surge: surge specific imports
import "src/layer1/surge/SurgeInbox.sol";
import "src/layer1/surge/SurgeTimelockController.sol";
import "src/layer1/verifiers/compose/AnyTwoVerifier.sol";

/// @title DeploySurgeOnL1
/// @notice This script deploys Taiko protocol modified for Nethermind's Surge.
contract DeploySurgeOnL1 is DeployCapability {
    uint256 internal immutable ADDRESS_LENGTH = 40;

    // Surge: Configurable values
    uint64 internal l2ChainId = uint64(vm.envUint("L2_CHAINID"));
    uint64 internal maxLivenessDisruptionPeriod =
        uint64(vm.envUint("MAX_LIVENESS_DISRUPTION_PERIOD"));
    uint64 internal minLivenessStreak = uint64(vm.envUint("MIN_LIVENESS_STREAK"));
    uint96 internal livenessBondBase = uint96(vm.envUint("LIVENESS_BOND_BASE"));
    uint96 internal livenessBondPerBlock = uint96(vm.envUint("LIVENESS_BOND_PER_BLOCK"));

    modifier broadcast() {
        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        require(privateKey != 0, "invalid private key");
        vm.startBroadcast();
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast {
        require(l2ChainId != block.chainid || l2ChainId != 0, "L2_CHAIN_ID");
        require(vm.envBytes32("L2_GENESIS_HASH") != 0, "L2_GENESIS_HASH");

        // Surge: Timelocked owner setup
        //---------------------------------------------------------------
        // Timelocked owner setup
        address[] memory executors = vm.envAddress("OWNER_MULTISIG_SIGNERS", ",");

        address ownerMultisig = vm.envAddress("OWNER_MULTISIG");
        addressNotNull(ownerMultisig, "ownerMultisig");

        address[] memory proposers = new address[](1);
        proposers[0] = ownerMultisig;

        uint256 timelockPeriod = uint64(vm.envUint("TIMELOCK_PERIOD"));
        address timelockController = address(
            new SurgeTimelockedController(
                minLivenessStreak, timelockPeriod, proposers, executors, address(0)
            )
        );
        address contractOwner = timelockController;

        console2.log("contractOwner(timelocked): ", contractOwner);

        // ---------------------------------------------------------------
        // Deploy shared contracts
        (address sharedResolver) = deploySharedContracts(contractOwner);
        console2.log("sharedResolver: ", sharedResolver);
        // ---------------------------------------------------------------
        // Deploy rollup contracts
        address rollupResolver = deployRollupContracts(sharedResolver, contractOwner);

        // ---------------------------------------------------------------
        // Signal service need to authorize the new rollup
        address signalServiceAddr = IResolver(sharedResolver).resolve(
            uint64(block.chainid), LibStrings.B_SIGNAL_SERVICE, false
        );
        SignalService signalService = SignalService(signalServiceAddr);

        address taikoInboxAddr =
            IResolver(rollupResolver).resolve(uint64(block.chainid), LibStrings.B_TAIKO, false);

        if (vm.envAddress("SHARED_RESOLVER") == address(0)) {
            SignalService(signalServiceAddr).authorize(taikoInboxAddr, true);
        }

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
            console2.log("signalService.authorize(taikoInboxAddr, bytes32(block.chainid))");
            console2.log("- signalService : ", signalServiceAddr);
            console2.log("- taikoInboxAddr   : ", taikoInboxAddr);
            console2.log("- chainId       : ", block.chainid);
        }

        // Surge: Dynamic L2 address based on passed chainid
        address taikoL2Address = getConstantAddress(vm.toString(l2ChainId), "10001");
        address l2SignalServiceAddress = getConstantAddress(vm.toString(l2ChainId), "5");
        address l2BridgeAddress = getConstantAddress(vm.toString(l2ChainId), "1");
        address l2Erc20VaultAddress = getConstantAddress(vm.toString(l2ChainId), "2");

        // ---------------------------------------------------------------
        // Register L2 addresses
        register(rollupResolver, "taiko", taikoL2Address, l2ChainId);
        register(rollupResolver, "signal_service", l2SignalServiceAddress, l2ChainId);
        register(sharedResolver, "signal_service", l2SignalServiceAddress, l2ChainId);
        register(sharedResolver, "bridge", l2BridgeAddress, l2ChainId);
        register(sharedResolver, "erc20_vault", l2Erc20VaultAddress, l2ChainId);

        // ---------------------------------------------------------------
        // Deploy other contracts
        if (block.chainid != 1) {
            deployAuxContracts();
        }

        // Surge: Removed preconf deployment

        if (DefaultResolver(sharedResolver).owner() == msg.sender) {
            DefaultResolver(sharedResolver).transferOwnership(contractOwner);
            console2.log("** sharedResolver ownership transferred to:", contractOwner);
        }

        DefaultResolver(rollupResolver).transferOwnership(contractOwner);
        console2.log("** rollupResolver ownership transferred to:", contractOwner);

        OwnableUpgradeable(taikoInboxAddr).transferOwnership(contractOwner);
    }

    function deploySharedContracts(address owner) internal returns (address sharedResolver) {
        addressNotNull(owner, "owner");

        sharedResolver = vm.envAddress("SHARED_RESOLVER");
        if (sharedResolver == address(0)) {
            sharedResolver = deployProxy({
                name: "shared_resolver",
                impl: address(new DefaultResolver()),
                data: abi.encodeCall(DefaultResolver.init, (address(0)))
            });
        }

        // Deploy Bridging contracts
        address signalService = deployProxy({
            name: "signal_service",
            impl: address(new MainnetSignalService(address(sharedResolver))),
            data: abi.encodeCall(SignalService.init, (address(0))),
            registerTo: sharedResolver
        });

        address quotaManager = address(0);
        address bridge = deployProxy({
            name: "bridge",
            impl: address(new MainnetBridge(address(sharedResolver), signalService, quotaManager)),
            data: abi.encodeCall(Bridge.init, (address(0))),
            registerTo: sharedResolver
        });

        Bridge(payable(bridge)).transferOwnership(owner);

        console2.log("------------------------------------------");
        console2.log(
            "Warning - you need to register *all* counterparty bridges to enable multi-hop bridging:"
        );
        console2.log(
            "sharedResolver.registerAddress(remoteChainId, 'bridge', address(remoteBridge))"
        );
        console2.log("- sharedResolver : ", sharedResolver);

        // Deploy Vaults
        address erc20Vault = deployProxy({
            name: "erc20_vault",
            impl: address(new MainnetERC20Vault(address(sharedResolver))),
            data: abi.encodeCall(ERC20Vault.init, (owner)),
            registerTo: sharedResolver
        });

        deployProxy({
            name: "erc721_vault",
            impl: address(new MainnetERC721Vault(address(sharedResolver))),
            data: abi.encodeCall(ERC721Vault.init, (owner)),
            registerTo: sharedResolver
        });

        deployProxy({
            name: "erc1155_vault",
            impl: address(new MainnetERC1155Vault(address(sharedResolver))),
            data: abi.encodeCall(ERC1155Vault.init, (owner)),
            registerTo: sharedResolver
        });

        console2.log("------------------------------------------");
        console2.log(
            "Warning - you need to register *all* counterparty vaults to enable multi-hop bridging:"
        );
        console2.log(
            "sharedResolver.registerAddress(remoteChainId, 'erc20_vault', address(remoteERC20Vault))"
        );
        console2.log(
            "sharedResolver.registerAddress(remoteChainId, 'erc721_vault', address(remoteERC721Vault))"
        );
        console2.log(
            "sharedResolver.registerAddress(remoteChainId, 'erc1155_vault', address(remoteERC1155Vault))"
        );
        console2.log("- sharedResolver : ", sharedResolver);

        // Deploy Bridged token implementations
        register(sharedResolver, "bridged_erc20", address(new BridgedERC20(erc20Vault)));
        register(
            sharedResolver, "bridged_erc721", address(new BridgedERC721(address(sharedResolver)))
        );
        register(
            sharedResolver, "bridged_erc1155", address(new BridgedERC1155(address(sharedResolver)))
        );
    }

    function deployRollupContracts(
        address _sharedResolver,
        address owner
    )
        internal
        returns (address rollupResolver)
    {
        addressNotNull(_sharedResolver, "sharedResolver");
        addressNotNull(owner, "owner");

        rollupResolver = deployProxy({
            name: "rollup_address_resolver",
            impl: address(new DefaultResolver()),
            data: abi.encodeCall(DefaultResolver.init, (address(0)))
        });

        // ---------------------------------------------------------------
        // Register shared contracts in the new rollup resolver
        copyRegister(rollupResolver, _sharedResolver, "signal_service");
        copyRegister(rollupResolver, _sharedResolver, "bridge");

        // Proof verifier
        address proofVerifierAddr = deployProxy({
            name: "proof_verifier",
            impl: address(new AnyTwoVerifier(address(0), address(0), address(0), address(0))),
            data: abi.encodeCall(ComposeVerifier.init, (address(0))),
            registerTo: rollupResolver
        });

        // Inbox
        address surgeInboxAddr = deployProxy({
            name: "taiko",
            impl: address(
                new SurgeInbox(
                    SurgeInbox.ConfigParams(0, 0, 0, 0), address(0), address(0), address(0), address(0)
                )
            ),
            data: ""
        });

        //------------------------------------------------------------------
        // Preconfirmations and forced inclusions contracts

        address whitelistAddr = deployProxy({
            name: "preconf_whitelist",
            impl: address(new PreconfWhitelist(rollupResolver)),
            data: abi.encodeCall(PreconfWhitelist.init, (owner)),
            registerTo: rollupResolver
        });

        address routerAddr = deployProxy({
            name: "preconf_router",
            impl: address(new PreconfRouter(address(0), address(0))),
            data: abi.encodeCall(PreconfRouter.init, address(0)),
            registerTo: rollupResolver
        });

        address storeAddr = deployProxy({
            name: "forced_inclusion_store",
            impl: address(new ForcedInclusionStore(0, 0, surgeInboxAddr, address(1))),
            data: abi.encodeCall(ForcedInclusionStore.init, (address(0))),
            registerTo: rollupResolver
        });

        address taikoWrapperAddr = deployProxy({
            name: "taiko_wrapper",
            impl: address(new TaikoWrapper(address(0), address(0), address(0))),
            data: abi.encodeCall(TaikoWrapper.init, address(0)),
            registerTo: rollupResolver
        });

        //------------------------------------------------------------------
        // Upgrades

        {
            address sgxVerifier =
                deploySgxVerifier(owner, rollupResolver, surgeInboxAddr, proofVerifierAddr);

            (address risc0Verifier, address sp1Verifier) = deployZKVerifiers(owner, rollupResolver);

            UUPSUpgradeable(proofVerifierAddr).upgradeTo({
                newImplementation: address(
                    new AnyTwoVerifier(surgeInboxAddr, sgxVerifier, risc0Verifier, sp1Verifier)
                )
            });
        }

        UUPSUpgradeable(surgeInboxAddr).upgradeTo({
            newImplementation: address(
                new SurgeInbox(
                    SurgeInbox.ConfigParams(
                        l2ChainId, maxLivenessDisruptionPeriod, livenessBondBase, livenessBondPerBlock
                    ),
                    taikoWrapperAddr,
                    proofVerifierAddr,
                    address(0),
                    IResolver(_sharedResolver).resolve(uint64(block.chainid), "signal_service", false)
                )
            )
        });

        UUPSUpgradeable(routerAddr).upgradeTo({
            newImplementation: address(new PreconfRouter(taikoWrapperAddr, whitelistAddr))
        });

        UUPSUpgradeable(taikoWrapperAddr).upgradeTo({
            newImplementation: address(new TaikoWrapper(surgeInboxAddr, storeAddr, routerAddr))
        });

        UUPSUpgradeable(storeAddr).upgradeTo(
            address(
                new ForcedInclusionStore(
                    uint8(vm.envUint("INCLUSION_WINDOW")),
                    uint64(vm.envUint("INCLUSION_FEE_IN_GWEI")),
                    surgeInboxAddr,
                    taikoWrapperAddr
                )
            )
        );

        //-------------------------------------------------------------------
        // Ownership transfers

        SurgeInbox(surgeInboxAddr).init(owner, vm.envBytes32("L2_GENESIS_HASH"));

        OwnableUpgradeable(proofVerifierAddr).transferOwnership(owner);
        console2.log("** proof_verifier ownership transferred to:", owner);

        OwnableUpgradeable(routerAddr).transferOwnership(owner);
        console2.log("** router ownership transferred to:", owner);

        OwnableUpgradeable(taikoWrapperAddr).transferOwnership(owner);
        console2.log("** taiko_wrapper ownership transferred to:", owner);

        OwnableUpgradeable(storeAddr).transferOwnership(owner);
        console2.log("** store ownership transferred to:", owner);
    }

    function deploySgxVerifier(
        address owner,
        address rollupResolver,
        address taikoInbox,
        address taikoProofVerifier
    )
        private
        returns (address sgxVerifier)
    {
        // No need to proxy these, because they are 3rd party. If we want to modify, we simply
        // change the registerAddress("automata_dcap_attestation", address(attestation));
        P256Verifier p256Verifier = new P256Verifier();
        SigVerifyLib sigVerifyLib = new SigVerifyLib(address(p256Verifier));
        PEMCertChainLib pemCertChainLib = new PEMCertChainLib();
        address automataDcapV3AttestationImpl = address(new AutomataDcapV3Attestation());

        address automataProxy = deployProxy({
            name: "automata_dcap_attestation",
            impl: automataDcapV3AttestationImpl,
            data: abi.encodeCall(
                AutomataDcapV3Attestation.init, (owner, address(sigVerifyLib), address(pemCertChainLib))
            ),
            registerTo: rollupResolver
        });

        sgxVerifier = deployProxy({
            name: "sgx_verifier",
            impl: address(new SgxVerifier(l2ChainId, taikoInbox, taikoProofVerifier, automataProxy)),
            data: abi.encodeCall(SgxVerifier.init, owner),
            registerTo: rollupResolver
        });

        // Log addresses for the user to register sgx instance
        console2.log("SigVerifyLib", address(sigVerifyLib));
        console2.log("PemCertChainLib", address(pemCertChainLib));
        console2.log("AutomataDcapVaAttestation", automataProxy);
    }

    function deployZKVerifiers(
        address owner,
        address rollupResolver
    )
        private
        returns (address risc0Verifier, address sp1Verifier)
    {
        // Deploy r0 groth16 verifier
        RiscZeroGroth16Verifier verifier =
            new RiscZeroGroth16Verifier(ControlID.CONTROL_ROOT, ControlID.BN254_CONTROL_ID);
        register(rollupResolver, "risc0_groth16_verifier", address(verifier));

        risc0Verifier = deployProxy({
            name: "risc0_verifier",
            impl: address(new Risc0Verifier(l2ChainId, address(verifier))),
            data: abi.encodeCall(Risc0Verifier.init, (owner)),
            registerTo: rollupResolver
        });

        // Deploy sp1 plonk verifier
        SuccinctVerifier succinctVerifier = new SuccinctVerifier();
        register(rollupResolver, "sp1_remote_verifier", address(succinctVerifier));

        sp1Verifier = deployProxy({
            name: "sp1_verifier",
            impl: address(new SP1Verifier(l2ChainId, address(succinctVerifier))),
            data: abi.encodeCall(SP1Verifier.init, (owner)),
            registerTo: rollupResolver
        });
    }

    function deployAuxContracts() private {
        address horseToken = address(new FreeMintERC20Token("Horse Token", "HORSE"));
        console2.log("HorseToken", horseToken);

        address bullToken =
            address(new FreeMintERC20Token_With50PctgMintAndTransferFailure("Bull Token", "BULL"));
        console2.log("BullToken", bullToken);
    }

    function addressNotNull(address addr, string memory err) private pure {
        require(addr != address(0), err);
    }

    function getConstantAddress(
        string memory prefix,
        string memory suffix
    )
        internal
        pure
        returns (address)
    {
        bytes memory prefixBytes = bytes(prefix);
        bytes memory suffixBytes = bytes(suffix);

        require(
            prefixBytes.length + suffixBytes.length <= ADDRESS_LENGTH, "Prefix + suffix too long"
        );

        // Create the middle padding of zeros
        uint256 paddingLength = ADDRESS_LENGTH - prefixBytes.length - suffixBytes.length;
        bytes memory padding = new bytes(paddingLength);
        for (uint256 i = 0; i < paddingLength; i++) {
            padding[i] = "0";
        }

        // Concatenate the parts
        string memory hexString = string(abi.encodePacked("0x", prefix, string(padding), suffix));

        return vm.parseAddress(hexString);
    }
}
