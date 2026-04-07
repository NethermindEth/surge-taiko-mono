// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {
    Ownable2StepUpgradeable
} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {
    UUPSUpgradeable
} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import { console2 } from "forge-std/src/console2.sol";
import { IInbox } from "src/layer1/core/iface/IInbox.sol";
import { PreconfWhitelist } from "src/layer1/preconf/impl/PreconfWhitelist.sol";
import { SurgeInbox } from "src/layer1/surge/deployments/internal-devnet/SurgeInbox.sol";
import { DeployCapability } from "test/shared/DeployCapability.sol";

/// @title SetupInbox
/// @notice Step 4: Deploys PreconfWhitelist and SurgeInbox implementation, upgrades the inbox
///         proxy, initializes it, and transfers ownership.
/// @dev Requires SURGE_INBOX and SURGE_VERIFIER from Step 1, SIGNAL_SERVICE from Step 3,
///      and EFFECTIVE_OWNER.
/// @custom:security-contact security@nethermind.io
contract SetupInbox is DeployCapability {
    // Signer configuration
    uint256 internal immutable privateKey = vm.envUint("PRIVATE_KEY");

    // Addresses from previous steps
    address internal immutable surgeInbox = vm.envAddress("SURGE_INBOX");
    address internal immutable surgeVerifier = vm.envAddress("SURGE_VERIFIER");
    address internal immutable signalService = vm.envAddress("SIGNAL_SERVICE");
    address internal immutable effectiveOwner = vm.envAddress("EFFECTIVE_OWNER");

    // Bond configuration
    address internal immutable bondToken = vm.envAddress("BOND_TOKEN");
    uint64 internal immutable minBond = uint64(vm.envUint("MIN_BOND"));
    uint64 internal immutable livenessBond = uint64(vm.envUint("LIVENESS_BOND"));
    uint48 internal immutable withdrawalDelay = uint48(vm.envUint("WITHDRAWAL_DELAY"));

    // Inbox configuration
    uint48 internal immutable provingWindow = uint48(vm.envUint("PROVING_WINDOW"));
    uint48 internal immutable permissionlessProvingDelay = 5 days;
    uint48 internal immutable maxProofSubmissionDelay =
        uint48(vm.envUint("MAX_PROOF_SUBMISSION_DELAY"));
    uint256 internal immutable ringBufferSize = vm.envUint("RING_BUFFER_SIZE");
    uint8 internal immutable basefeeSharingPctg = uint8(vm.envUint("BASEFEE_SHARING_PCTG"));
    uint256 internal immutable minForcedInclusionCount = vm.envUint("MIN_FORCED_INCLUSION_COUNT");
    uint16 internal immutable forcedInclusionDelay = uint16(vm.envUint("FORCED_INCLUSION_DELAY"));
    uint64 internal immutable forcedInclusionFeeInGwei =
        uint64(vm.envUint("FORCED_INCLUSION_FEE_IN_GWEI"));
    uint64 internal immutable forcedInclusionFeeDoubleThreshold =
        uint64(vm.envUint("FORCED_INCLUSION_FEE_DOUBLE_THRESHOLD"));
    uint16 internal immutable minCheckpointDelay = uint16(vm.envUint("MIN_CHECKPOINT_DELAY"));
    uint8 internal immutable permissionlessInclusionMultiplier =
        uint8(vm.envUint("PERMISSIONLESS_INCLUSION_MULTIPLIER"));

    // Finalization streak configuration
    uint48 internal immutable maxFinalizationDelayBeforeStreakReset =
        uint48(vm.envUint("MAX_FINALIZATION_DELAY_BEFORE_STREAK_RESET"));

    // Rollback configuration
    uint48 internal immutable maxFinalizationDelayBeforeRollback =
        uint48(vm.envUint("MAX_FINALIZATION_DELAY_BEFORE_ROLLBACK"));

    modifier broadcast() {
        require(privateKey != 0, "invalid private key");
        vm.startBroadcast(privateKey);
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast {
        require(surgeInbox != address(0), "config: SURGE_INBOX");
        require(surgeVerifier != address(0), "config: SURGE_VERIFIER");
        require(signalService != address(0), "config: SIGNAL_SERVICE");
        require(effectiveOwner != address(0), "config: EFFECTIVE_OWNER");

        console2.log("=====================================");
        console2.log("Step 4: Setup Inbox");
        console2.log("=====================================");
        console2.log("** Surge Inbox proxy:", surgeInbox);
        console2.log("** SurgeVerifier:", surgeVerifier);
        console2.log("** SignalService:", signalService);
        console2.log("** Effective owner:", effectiveOwner);

        // Deploy whitelist
        address whitelist = deployProxy({
            name: "preconf_whitelist",
            impl: address(new PreconfWhitelist()),
            data: abi.encodeCall(PreconfWhitelist.init, (effectiveOwner))
        });

        // Build inbox configuration
        IInbox.Config memory config = IInbox.Config({
            proofVerifier: surgeVerifier,
            proposerChecker: whitelist,
            proverWhitelist: address(0),
            signalService: signalService,
            bondToken: bondToken,
            minBond: minBond,
            livenessBond: livenessBond,
            withdrawalDelay: withdrawalDelay,
            provingWindow: provingWindow,
            permissionlessProvingDelay: permissionlessProvingDelay,
            maxProofSubmissionDelay: maxProofSubmissionDelay,
            ringBufferSize: ringBufferSize,
            basefeeSharingPctg: basefeeSharingPctg,
            minForcedInclusionCount: minForcedInclusionCount,
            forcedInclusionDelay: forcedInclusionDelay,
            forcedInclusionFeeInGwei: forcedInclusionFeeInGwei,
            forcedInclusionFeeDoubleThreshold: forcedInclusionFeeDoubleThreshold,
            minCheckpointDelay: minCheckpointDelay,
            permissionlessInclusionMultiplier: permissionlessInclusionMultiplier
        });

        // Deploy inbox implementation
        address inboxImpl = address(
            new SurgeInbox(
                config, maxFinalizationDelayBeforeStreakReset, maxFinalizationDelayBeforeRollback
            )
        );
        console2.log("** Deployed SurgeInbox implementation:", inboxImpl);

        // Upgrade inbox proxy to actual implementation
        UUPSUpgradeable(surgeInbox).upgradeTo(inboxImpl);
        console2.log("** Inbox proxy upgraded");

        // Initialize inbox
        SurgeInbox(payable(surgeInbox)).init(msg.sender);
        console2.log("** SurgeInbox initialized");

        // Transfer inbox ownership (requires acceptance)
        Ownable2StepUpgradeable(surgeInbox).transferOwnership(effectiveOwner);
        console2.log("** Inbox ownership transfer initiated to:", effectiveOwner);

        // Write addresses to step-specific JSON
        writeJson("preconf_whitelist", whitelist);
        writeJson("surge_inbox_impl", inboxImpl);

        console2.log("=====================================");
        console2.log("Step 4 Complete");
        console2.log("=====================================");
    }

    function writeJson(string memory name, address addr) internal {
        vm.writeJson(
            vm.serializeAddress("setup_inbox", name, addr),
            string.concat(vm.projectRoot(), "/deployments/setup_inbox.json")
        );
    }
}
