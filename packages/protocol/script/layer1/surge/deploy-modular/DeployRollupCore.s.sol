// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { EmptyImpl } from "../common/EmptyImpl.sol";
import { console2 } from "forge-std/src/console2.sol";
import { SurgeTimelockController } from "src/layer1/surge/SurgeTimelockController.sol";
import { SurgeVerifier } from "src/layer1/surge/SurgeVerifier.sol";
import { DeployCapability } from "test/shared/DeployCapability.sol";

/// @title DeployRollupCore
/// @notice Step 1: Deploys core rollup infrastructure (EmptyImpl, Inbox proxy, SurgeVerifier,
/// optional Timelock).
/// @dev Run this before all other modular deployment steps.
/// @custom:security-contact security@nethermind.io
contract DeployRollupCore is DeployCapability {
    // Signer configuration
    uint256 internal immutable privateKey = vm.envUint("PRIVATE_KEY");

    // Owner configuration
    address internal immutable contractOwner = vm.envAddress("CONTRACT_OWNER");

    // L2 configuration
    uint64 internal immutable l2ChainId = uint64(vm.envUint("L2_CHAIN_ID"));

    // SurgeVerifier configuration
    uint8 internal immutable numProofsThreshold = uint8(vm.envUint("NUM_PROOFS_THRESHOLD"));

    // Timelock configuration
    bool internal immutable useTimelock = vm.envBool("USE_TIMELOCK");
    uint256 internal immutable timelockMinDelay = vm.envUint("TIMELOCK_MIN_DELAY");
    uint48 internal immutable timelockMinFinalizationStreak =
        uint48(vm.envUint("TIMELOCK_MIN_FINALIZATION_STREAK"));
    address[] internal timelockProposers = vm.envAddress("TIMELOCK_PROPOSERS", ",");
    address[] internal timelockExecutors = vm.envAddress("TIMELOCK_EXECUTORS", ",");

    modifier broadcast() {
        require(privateKey != 0, "invalid private key");
        vm.startBroadcast(privateKey);
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast {
        require(l2ChainId != block.chainid || l2ChainId != 0, "config: L2_CHAIN_ID");
        require(contractOwner != address(0), "config: CONTRACT_OWNER");

        console2.log("=====================================");
        console2.log("Step 1: Deploy Rollup Core");
        console2.log("=====================================");
        console2.log("** Contract owner:", contractOwner);

        // Deploy empty implementation for temporary proxy use
        address emptyImpl = address(new EmptyImpl());
        console2.log("** Deployed EmptyImpl:", emptyImpl);

        // Deploy inbox proxy with empty implementation
        address inbox = deployProxy({ name: "surge_inbox", impl: emptyImpl, data: "" });
        console2.log("** Deployed Inbox proxy:", inbox);

        // Deploy SurgeVerifier (deployer is initial owner so step 2 can configure it)
        address verifier = address(new SurgeVerifier(inbox, numProofsThreshold, msg.sender));
        console2.log("** Deployed SurgeVerifier:", verifier);

        // Compute effective owner
        address effectiveOwner = contractOwner;
        if (useTimelock) {
            effectiveOwner = _deployTimelock(inbox, verifier);
            console2.log("** Effective owner (timelock):", effectiveOwner);
        } else {
            console2.log("** Effective owner:", effectiveOwner);
        }

        // Write all addresses to step-specific JSON
        writeJson("empty_impl", emptyImpl);
        writeJson("surge_inbox", inbox);
        writeJson("surge_verifier", verifier);
        writeJson("effective_owner", effectiveOwner);

        console2.log("=====================================");
        console2.log("Step 1 Complete");
        console2.log("=====================================");
    }

    function _deployTimelock(
        address _inbox,
        address _proofVerifier
    )
        internal
        returns (address timelock)
    {
        require(timelockMinDelay > 0, "config: TIMELOCK_MIN_DELAY");
        require(timelockMinFinalizationStreak > 0, "config: TIMELOCK_MIN_FINALIZATION_STREAK");
        require(timelockProposers.length > 0, "config: TIMELOCK_PROPOSERS");
        require(timelockExecutors.length > 0, "config: TIMELOCK_EXECUTORS");

        timelock = address(
            new SurgeTimelockController(
                _inbox,
                _proofVerifier,
                timelockMinFinalizationStreak,
                timelockMinDelay,
                timelockProposers,
                timelockExecutors
            )
        );
        writeJson("surge_timelock", timelock);
        console2.log("** Deployed SurgeTimelockController:", timelock);
    }

    function writeJson(string memory name, address addr) internal {
        vm.writeJson(
            vm.serializeAddress("deploy_rollup_core", name, addr),
            string.concat(vm.projectRoot(), "/deployments/deploy_rollup_core.json")
        );
    }
}
