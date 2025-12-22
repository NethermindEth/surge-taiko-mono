// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {
    Ownable2StepUpgradeable
} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { Ownable2Step } from "@openzeppelin/contracts/access/Ownable2Step.sol";
import { Script } from "forge-std/src/Script.sol";
import { console2 } from "forge-std/src/console2.sol";
import { SurgeTimelockController } from "src/layer1/surge/SurgeTimelockController.sol";

/// @title AcceptOwnership
/// @notice Script to accept ownership of multiple contracts that use Ownable2Step pattern.
contract AcceptOwnership is Script {
    uint256 internal immutable privateKey = vm.envUint("PRIVATE_KEY");
    address[] internal contracts = vm.envAddress("CONTRACT_ADDRESSES", ",");
    bool internal useTimelock = vm.envBool("USE_TIMELOCK");
    address internal timelockController = vm.envAddress("TIMELOCK_CONTROLLER");

    modifier broadcast() {
        require(privateKey != 0, "invalid private key");
        vm.startBroadcast(privateKey);
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast {
        require(contracts.length > 0, "No contract addresses provided");

        console2.log("=====================================");
        console2.log("Accepting Ownership for", contracts.length, "contracts");
        console2.log("=====================================");

        if (useTimelock) {
            require(timelockController != address(0), "No timelock controller specified");
            SurgeTimelockController(payable(timelockController)).acceptOwnership(contracts);
        } else {
            for (uint256 i = 0; i < contracts.length; i++) {
                acceptOwnership(contracts[i]);
            }
        }

        console2.log("=====================================");
        console2.log("Ownership Acceptance Complete");
        console2.log("=====================================");
    }

    function acceptOwnership(address _contract) internal {
        console2.log("Processing contract:", _contract);

        // Try Ownable2StepUpgradeable first (proxied contracts)
        try Ownable2StepUpgradeable(_contract).pendingOwner() returns (address pendingOwner) {
            require(pendingOwner == msg.sender, "Caller is not the pending owner");
            Ownable2StepUpgradeable(_contract).acceptOwnership();
            console2.log("  -> Accepted ownership (upgradeable)");
            console2.log("  -> New owner:", Ownable2StepUpgradeable(_contract).owner());
            return;
        } catch {
            // Not an upgradeable contract, try regular Ownable2Step
        }

        // Try regular Ownable2Step (non-proxied contracts)
        try Ownable2Step(_contract).pendingOwner() returns (address pendingOwner) {
            require(pendingOwner == msg.sender, "Caller is not the pending owner");
            Ownable2Step(_contract).acceptOwnership();
            console2.log("  -> Accepted ownership (non-upgradeable)");
            console2.log("  -> New owner:", Ownable2Step(_contract).owner());
            return;
        } catch {
            revert("Contract does not support Ownable2Step pattern");
        }
    }
}
