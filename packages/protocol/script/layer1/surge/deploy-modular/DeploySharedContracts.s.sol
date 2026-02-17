// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { console2 } from "forge-std/src/console2.sol";
import { Bridge } from "src/shared/bridge/Bridge.sol";
import { DefaultResolver } from "src/shared/common/DefaultResolver.sol";
import { SignalService } from "src/shared/signal/SignalService.sol";
import { BridgedERC1155 } from "src/shared/vault/BridgedERC1155.sol";
import { BridgedERC20 } from "src/shared/vault/BridgedERC20.sol";
import { BridgedERC721 } from "src/shared/vault/BridgedERC721.sol";
import { ERC1155Vault } from "src/shared/vault/ERC1155Vault.sol";
import { ERC20Vault } from "src/shared/vault/ERC20Vault.sol";
import { ERC721Vault } from "src/shared/vault/ERC721Vault.sol";
import { DeployCapability } from "test/shared/DeployCapability.sol";

/// @title DeploySharedContracts
/// @notice Step 3: Deploys shared infrastructure (Resolver, SignalService, Bridge, Vaults)
///         and registers L2 addresses.
/// @dev Requires SURGE_INBOX and EFFECTIVE_OWNER from Step 1.
///      Can run in parallel with Step 2 (no mutual dependencies).
/// @custom:security-contact security@nethermind.io
contract DeploySharedContracts is DeployCapability {
    uint256 internal constant ADDRESS_LENGTH = 40;

    // Signer configuration
    uint256 internal immutable privateKey = vm.envUint("PRIVATE_KEY");

    // L2 configuration
    uint64 internal immutable l2ChainId = uint64(vm.envUint("L2_CHAIN_ID"));

    // Addresses from Step 1
    address internal immutable surgeInbox = vm.envAddress("SURGE_INBOX");
    address internal immutable effectiveOwner = vm.envAddress("EFFECTIVE_OWNER");

    modifier broadcast() {
        require(privateKey != 0, "invalid private key");
        vm.startBroadcast(privateKey);
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast {
        require(surgeInbox != address(0), "config: SURGE_INBOX");
        require(effectiveOwner != address(0), "config: EFFECTIVE_OWNER");
        require(l2ChainId != block.chainid || l2ChainId != 0, "config: L2_CHAIN_ID");

        console2.log("=====================================");
        console2.log("Step 3: Deploy Shared Contracts");
        console2.log("=====================================");
        console2.log("** Surge Inbox:", surgeInbox);
        console2.log("** Effective owner:", effectiveOwner);

        // Deploy shared resolver
        address sharedResolver = deployProxy({
            name: "shared_resolver",
            impl: address(new DefaultResolver()),
            data: abi.encodeCall(DefaultResolver.init, (address(0)))
        });

        // Deploy signal service
        address signalService = deployProxy({
            name: "signal_service",
            impl: address(new SignalService(surgeInbox, getL2SignalServiceAddress())),
            data: abi.encodeCall(SignalService.init, (effectiveOwner)),
            registerTo: sharedResolver
        });

        // Deploy bridge
        address bridge = deployProxy({
            name: "bridge",
            impl: address(new Bridge(sharedResolver, signalService)),
            data: abi.encodeCall(Bridge.init, (effectiveOwner)),
            registerTo: sharedResolver
        });

        // Deploy vaults
        address erc20Vault = deployProxy({
            name: "erc20_vault",
            impl: address(new ERC20Vault(sharedResolver)),
            data: abi.encodeCall(ERC20Vault.init, (effectiveOwner)),
            registerTo: sharedResolver
        });

        address erc721Vault = deployProxy({
            name: "erc721_vault",
            impl: address(new ERC721Vault(sharedResolver)),
            data: abi.encodeCall(ERC721Vault.init, (effectiveOwner)),
            registerTo: sharedResolver
        });

        address erc1155Vault = deployProxy({
            name: "erc1155_vault",
            impl: address(new ERC1155Vault(sharedResolver)),
            data: abi.encodeCall(ERC1155Vault.init, (effectiveOwner)),
            registerTo: sharedResolver
        });

        // Deploy bridged token implementations (clone pattern)
        address bridgedErc20 = address(new BridgedERC20(erc20Vault));
        register(sharedResolver, "bridged_erc20", bridgedErc20);

        address bridgedErc721 = address(new BridgedERC721(erc721Vault));
        register(sharedResolver, "bridged_erc721", bridgedErc721);

        address bridgedErc1155 = address(new BridgedERC1155(erc1155Vault));
        register(sharedResolver, "bridged_erc1155", bridgedErc1155);

        // Register L2 addresses in the resolver
        register(sharedResolver, "signal_service", getL2SignalServiceAddress(), l2ChainId);
        register(sharedResolver, "bridge", getL2BridgeAddress(), l2ChainId);
        register(sharedResolver, "erc20_vault", getL2Erc20VaultAddress(), l2ChainId);
        register(sharedResolver, "erc721_vault", getL2Erc721VaultAddress(), l2ChainId);
        register(sharedResolver, "erc1155_vault", getL2Erc1155VaultAddress(), l2ChainId);
        console2.log("** L2 addresses registered");

        // Transfer SharedResolver ownership (requires acceptance)
        DefaultResolver(sharedResolver).transferOwnership(effectiveOwner);
        console2.log("** SharedResolver ownership transfer initiated to:", effectiveOwner);

        // Write all addresses to step-specific JSON
        writeJson("shared_resolver", sharedResolver);
        writeJson("signal_service", signalService);
        writeJson("bridge", bridge);
        writeJson("erc20_vault", erc20Vault);
        writeJson("erc721_vault", erc721Vault);
        writeJson("erc1155_vault", erc1155Vault);
        writeJson("bridged_erc20", bridgedErc20);
        writeJson("bridged_erc721", bridgedErc721);
        writeJson("bridged_erc1155", bridgedErc1155);

        console2.log("=====================================");
        console2.log("Step 3 Complete");
        console2.log("=====================================");
    }

    // ---------------------------------------------------------------
    // L2 Address Getters
    // ---------------------------------------------------------------

    function getL2BridgeAddress() internal view returns (address) {
        return getConstantAddress(vm.toString(l2ChainId), "1");
    }

    function getL2Erc20VaultAddress() internal view returns (address) {
        return getConstantAddress(vm.toString(l2ChainId), "2");
    }

    function getL2Erc721VaultAddress() internal view returns (address) {
        return getConstantAddress(vm.toString(l2ChainId), "3");
    }

    function getL2Erc1155VaultAddress() internal view returns (address) {
        return getConstantAddress(vm.toString(l2ChainId), "4");
    }

    function getL2SignalServiceAddress() internal view returns (address) {
        return getConstantAddress(vm.toString(l2ChainId), "5");
    }

    // ---------------------------------------------------------------
    // Utilities
    // ---------------------------------------------------------------

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

        uint256 paddingLength = ADDRESS_LENGTH - prefixBytes.length - suffixBytes.length;
        bytes memory padding = new bytes(paddingLength);
        for (uint256 i = 0; i < paddingLength; i++) {
            padding[i] = "0";
        }

        string memory hexString = string(abi.encodePacked("0x", prefix, string(padding), suffix));

        return vm.parseAddress(hexString);
    }

    function writeJson(string memory name, address addr) internal {
        vm.writeJson(
            vm.serializeAddress("deploy_shared_contracts", name, addr),
            string.concat(vm.projectRoot(), "/deployments/deploy_shared_contracts.json")
        );
    }
}
