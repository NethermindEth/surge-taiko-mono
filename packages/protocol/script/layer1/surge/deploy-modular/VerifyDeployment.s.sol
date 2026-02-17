// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {
    Ownable2StepUpgradeable
} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { Strings } from "@openzeppelin/contracts/utils/Strings.sol";
import { Script } from "forge-std/src/Script.sol";
import { console2 } from "forge-std/src/console2.sol";
import { DefaultResolver } from "src/shared/common/DefaultResolver.sol";

/// @title VerifyDeployment
/// @notice Step 5: Read-only verification of all L1 registrations, L2 registrations,
///         and ownership state.
/// @dev Requires all contract addresses from Steps 1-4.
///      This script does not broadcast any transactions.
/// @custom:security-contact security@nethermind.io
contract VerifyDeployment is Script {
    uint256 internal constant ADDRESS_LENGTH = 40;

    // L2 configuration
    uint64 internal immutable l2ChainId = uint64(vm.envUint("L2_CHAIN_ID"));

    // Verifier configuration (affects ownership verification)
    bool internal immutable useDummyVerifier = vm.envBool("USE_DUMMY_VERIFIER");

    // Addresses from Step 1
    address internal immutable surgeInbox = vm.envAddress("SURGE_INBOX");
    address internal immutable surgeVerifier = vm.envAddress("SURGE_VERIFIER");
    address internal immutable effectiveOwner = vm.envAddress("EFFECTIVE_OWNER");

    // Addresses from Step 2 (may be zero if not deployed)
    address internal immutable risc0Verifier = vm.envAddress("RISC0_VERIFIER");
    address internal immutable sp1Verifier = vm.envAddress("SP1_VERIFIER");

    // Addresses from Step 3
    address internal immutable sharedResolver = vm.envAddress("SHARED_RESOLVER");
    address internal immutable signalService = vm.envAddress("SIGNAL_SERVICE");
    address internal immutable bridge = vm.envAddress("BRIDGE");
    address internal immutable erc20Vault = vm.envAddress("ERC20_VAULT");
    address internal immutable erc721Vault = vm.envAddress("ERC721_VAULT");
    address internal immutable erc1155Vault = vm.envAddress("ERC1155_VAULT");

    function run() external view {
        console2.log("=====================================");
        console2.log("Step 5: Verify Deployment");
        console2.log("=====================================");
        console2.log("** Effective owner:", effectiveOwner);

        // Verify L1 registrations
        _verifyL1Registrations();

        // Verify L2 registrations
        _verifyL2Registrations();

        // Verify ownership
        _verifyOwnership();

        console2.log("=====================================");
        console2.log("Deployment Verified Successfully");
        console2.log("=====================================");
    }

    function _verifyL1Registrations() internal view {
        bytes32[] memory sharedNames = new bytes32[](8);
        sharedNames[0] = bytes32("signal_service");
        sharedNames[1] = bytes32("bridge");
        sharedNames[2] = bytes32("erc20_vault");
        sharedNames[3] = bytes32("erc721_vault");
        sharedNames[4] = bytes32("erc1155_vault");
        sharedNames[5] = bytes32("bridged_erc20");
        sharedNames[6] = bytes32("bridged_erc721");
        sharedNames[7] = bytes32("bridged_erc1155");

        for (uint256 i = 0; i < sharedNames.length; i++) {
            try DefaultResolver(sharedResolver).resolve(
                block.chainid, sharedNames[i], false
            ) returns (address) { }
            catch {
                revert(
                    string.concat(
                        "verifyL1Registrations: missing registration for ",
                        Strings.toHexString(uint256(sharedNames[i]))
                    )
                );
            }
        }
        console2.log("** L1 registrations verified");
    }

    function _verifyL2Registrations() internal view {
        require(
            DefaultResolver(sharedResolver).resolve(l2ChainId, bytes32("signal_service"), false)
                == getL2SignalServiceAddress(),
            "verifyL2Registrations: signal_service mismatch"
        );
        require(
            DefaultResolver(sharedResolver).resolve(l2ChainId, bytes32("bridge"), false)
                == getL2BridgeAddress(),
            "verifyL2Registrations: bridge mismatch"
        );
        require(
            DefaultResolver(sharedResolver).resolve(l2ChainId, bytes32("erc20_vault"), false)
                == getL2Erc20VaultAddress(),
            "verifyL2Registrations: erc20_vault mismatch"
        );
        require(
            DefaultResolver(sharedResolver).resolve(l2ChainId, bytes32("erc721_vault"), false)
                == getL2Erc721VaultAddress(),
            "verifyL2Registrations: erc721_vault mismatch"
        );
        require(
            DefaultResolver(sharedResolver).resolve(l2ChainId, bytes32("erc1155_vault"), false)
                == getL2Erc1155VaultAddress(),
            "verifyL2Registrations: erc1155_vault mismatch"
        );
        console2.log("** L2 registrations verified");
    }

    function _verifyOwnership() internal view {
        // Contracts with immediate ownership (owner = effective owner)
        address[] memory ownerContracts = new address[](5);
        ownerContracts[0] = signalService;
        ownerContracts[1] = bridge;
        ownerContracts[2] = erc20Vault;
        ownerContracts[3] = erc721Vault;
        ownerContracts[4] = erc1155Vault;

        // Contracts with pending ownership (pendingOwner = effective owner)
        address[] memory pendingOwnerContracts;
        if (useDummyVerifier) {
            pendingOwnerContracts = new address[](3);
            pendingOwnerContracts[0] = surgeVerifier;
            pendingOwnerContracts[1] = surgeInbox;
            pendingOwnerContracts[2] = sharedResolver;
        } else {
            pendingOwnerContracts = new address[](5);
            pendingOwnerContracts[0] = surgeVerifier;
            pendingOwnerContracts[1] = surgeInbox;
            pendingOwnerContracts[2] = sharedResolver;
            pendingOwnerContracts[3] = risc0Verifier; // May be address(0)
            pendingOwnerContracts[4] = sp1Verifier; // May be address(0)
        }

        // Verify current ownership
        for (uint256 i; i < ownerContracts.length; ++i) {
            if (ownerContracts[i] == address(0)) continue;
            address currentOwner = Ownable2StepUpgradeable(ownerContracts[i]).owner();
            require(
                currentOwner == effectiveOwner,
                string.concat(
                    "verifyOwnership: ",
                    Strings.toHexString(uint160(ownerContracts[i]), 20)
                )
            );
        }

        // Verify pending ownership
        for (uint256 i; i < pendingOwnerContracts.length; ++i) {
            if (pendingOwnerContracts[i] == address(0)) continue;
            address pendingOwner =
                Ownable2StepUpgradeable(pendingOwnerContracts[i]).pendingOwner();
            require(
                pendingOwner == effectiveOwner,
                string.concat(
                    "verifyPendingOwnership: ",
                    Strings.toHexString(uint160(pendingOwnerContracts[i]), 20)
                )
            );
        }

        console2.log("** Ownership verified");
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
}
