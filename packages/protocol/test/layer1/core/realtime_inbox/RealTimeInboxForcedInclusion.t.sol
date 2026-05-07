// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// forge-config: default.isolate = true

import { RealTimeInboxTestBase } from "./RealTimeInboxTestBase.sol";
import { IForcedInclusionStore } from "src/layer1/core/iface/IForcedInclusionStore.sol";
import { IInbox } from "src/layer1/core/iface/IInbox.sol";
import { IRealTimeInbox } from "src/layer1/core/iface/IRealTimeInbox.sol";
import { RealTimeInbox } from "src/layer1/core/impl/RealTimeInbox.sol";
import { LibBlobs } from "src/layer1/core/libs/LibBlobs.sol";
import { LibForcedInclusion } from "src/layer1/core/libs/LibForcedInclusion.sol";
import { ICheckpointStore } from "src/shared/signal/ICheckpointStore.sol";

/// @notice Tests for RealTimeInbox forced-inclusion enqueue + propose-time consumption.
contract RealTimeInboxForcedInclusionTest is RealTimeInboxTestBase {
    address internal forcer = Carol;

    function setUp() public override {
        super.setUp();
        vm.deal(forcer, 100 ether);
    }

    // ---------------------------------------------------------------
    // saveForcedInclusion()
    // ---------------------------------------------------------------

    function test_saveForcedInclusion_succeeds() public {
        LibBlobs.BlobReference memory ref =
            LibBlobs.BlobReference({ blobStartIndex: 0, numBlobs: 1, offset: 0 });

        uint64 fee = inbox.getCurrentForcedInclusionFee();

        _setBlobHashes(1);
        vm.prank(forcer);
        inbox.saveForcedInclusion{ value: fee * 1 gwei }(ref);

        (uint48 head, uint48 tail) = inbox.getForcedInclusionState();
        assertEq(head, 0, "head untouched");
        assertEq(tail, 1, "tail advanced");

        IForcedInclusionStore.ForcedInclusion[] memory entries = inbox.getForcedInclusions(0, 10);
        assertEq(entries.length, 1, "single entry");
        assertEq(entries[0].feeInGwei, fee, "fee matches");
        assertEq(entries[0].blobSlice.blobHashes.length, 1, "single blob");
        assertEq(entries[0].blobSlice.timestamp, uint48(block.timestamp), "ts is now");
    }

    function test_saveForcedInclusion_RevertWhen_NotActivated() public {
        RealTimeInbox freshInbox = _deployNonActivatedInbox();

        LibBlobs.BlobReference memory ref =
            LibBlobs.BlobReference({ blobStartIndex: 0, numBlobs: 1, offset: 0 });

        _setBlobHashes(1);
        vm.expectRevert(RealTimeInbox.NotActivated.selector);
        vm.prank(forcer);
        freshInbox.saveForcedInclusion{ value: 1 ether }(ref);
    }

    function test_saveForcedInclusion_RevertWhen_InsufficientFee() public {
        LibBlobs.BlobReference memory ref =
            LibBlobs.BlobReference({ blobStartIndex: 0, numBlobs: 1, offset: 0 });

        _setBlobHashes(1);
        vm.expectRevert(LibForcedInclusion.InsufficientFee.selector);
        vm.prank(forcer);
        inbox.saveForcedInclusion{ value: 0 }(ref);
    }

    function test_saveForcedInclusion_RevertWhen_MultipleBlobs() public {
        LibBlobs.BlobReference memory ref =
            LibBlobs.BlobReference({ blobStartIndex: 0, numBlobs: 2, offset: 0 });

        uint64 fee = inbox.getCurrentForcedInclusionFee();
        _setBlobHashes(2);
        vm.expectRevert(LibForcedInclusion.OnlySingleBlobAllowed.selector);
        vm.prank(forcer);
        inbox.saveForcedInclusion{ value: fee * 1 gwei }(ref);
    }

    function test_saveForcedInclusion_refundsExcess() public {
        LibBlobs.BlobReference memory ref =
            LibBlobs.BlobReference({ blobStartIndex: 0, numBlobs: 1, offset: 0 });

        uint64 fee = inbox.getCurrentForcedInclusionFee();
        uint256 paid = uint256(fee) * 1 gwei + 1 ether; // overpay by 1 ether
        uint256 forcerBalanceBefore = forcer.balance;

        _setBlobHashes(1);
        vm.prank(forcer);
        inbox.saveForcedInclusion{ value: paid }(ref);

        // Net cost = fee, refund = 1 ether
        assertEq(
            forcer.balance, forcerBalanceBefore - uint256(fee) * 1 gwei, "refund of excess only"
        );
    }

    function test_getCurrentForcedInclusionFee_doublesAtThreshold() public {
        // threshold=100 (test-base config), so 100 pending → 2× base fee
        uint64 baseFee = inbox.getCurrentForcedInclusionFee();
        assertEq(baseFee, 1_000_000, "empty queue: base fee");

        LibBlobs.BlobReference memory ref =
            LibBlobs.BlobReference({ blobStartIndex: 0, numBlobs: 1, offset: 0 });

        // Enqueue 100 inclusions
        for (uint256 i; i < 100; ++i) {
            _setBlobHashes(1);
            uint64 fee = inbox.getCurrentForcedInclusionFee();
            vm.prank(forcer);
            inbox.saveForcedInclusion{ value: uint256(fee) * 1 gwei }(ref);
        }

        // Fee at queue size 100 should be 2× base
        uint64 doubled = inbox.getCurrentForcedInclusionFee();
        assertEq(doubled, baseFee * 2, "fee doubled at threshold");
    }

    // ---------------------------------------------------------------
    // propose() with forced inclusions
    // ---------------------------------------------------------------

    function test_propose_consumesSingleForcedInclusion() public {
        // Enqueue one FI
        _enqueueOneForcedInclusion();

        // Build proposal that consumes it
        IRealTimeInbox.ProposeInput memory input = _buildDefaultProposeInput();
        input.numForcedInclusions = 1;
        ICheckpointStore.Checkpoint memory checkpoint = _buildCheckpoint();

        _proposeAndGetLogs(input, checkpoint);

        // Queue head advanced
        (uint48 head, uint48 tail) = inbox.getForcedInclusionState();
        assertEq(head, 1, "head advanced");
        assertEq(tail, 1, "tail unchanged");
    }

    function test_propose_consumesMultipleForcedInclusions() public {
        _enqueueOneForcedInclusion();
        _enqueueOneForcedInclusion();
        _enqueueOneForcedInclusion();

        IRealTimeInbox.ProposeInput memory input = _buildDefaultProposeInput();
        input.numForcedInclusions = 3;
        ICheckpointStore.Checkpoint memory checkpoint = _buildCheckpoint();

        _proposeAndGetLogs(input, checkpoint);

        (uint48 head, uint48 tail) = inbox.getForcedInclusionState();
        assertEq(head, 3, "head advanced by 3");
        assertEq(tail, 3, "tail unchanged");
    }

    function test_propose_emptyQueue_numForcedInclusionsClamped() public {
        // Empty queue, but proposer asks for 5 — should clamp to 0 and succeed.
        IRealTimeInbox.ProposeInput memory input = _buildDefaultProposeInput();
        input.numForcedInclusions = 5;
        ICheckpointStore.Checkpoint memory checkpoint = _buildCheckpoint();

        _proposeAndGetLogs(input, checkpoint);

        (uint48 head, uint48 tail) = inbox.getForcedInclusionState();
        assertEq(head, 0, "head untouched");
        assertEq(tail, 0, "tail untouched");
    }

    function test_propose_RevertWhen_OldestForcedInclusionIsDue() public {
        _enqueueOneForcedInclusion();

        // Warp past forcedInclusionDelay (3600 in test base)
        vm.warp(block.timestamp + 3601);
        vm.roll(block.number + 1);

        // Try to propose without consuming the due FI
        IRealTimeInbox.ProposeInput memory input = _buildDefaultProposeInput();
        input.numForcedInclusions = 0;
        input.maxAnchorBlockNumber = uint48(block.number - 1);
        ICheckpointStore.Checkpoint memory checkpoint = _buildCheckpoint();

        bytes memory data = abi.encode(input);
        _setBlobHashes(1);

        vm.expectRevert(RealTimeInbox.UnprocessedForcedInclusionIsDue.selector);
        vm.prank(proposer);
        inbox.propose(data, checkpoint, bytes(""));
    }

    function test_propose_OldestDueIsConsumed_succeeds() public {
        _enqueueOneForcedInclusion();

        vm.warp(block.timestamp + 3601);
        vm.roll(block.number + 1);

        IRealTimeInbox.ProposeInput memory input = _buildDefaultProposeInput();
        input.numForcedInclusions = 1;
        input.maxAnchorBlockNumber = uint48(block.number - 1);
        ICheckpointStore.Checkpoint memory checkpoint = _buildCheckpoint();

        _proposeAndGetLogs(input, checkpoint);

        (uint48 head,) = inbox.getForcedInclusionState();
        assertEq(head, 1, "head advanced");
    }

    function test_propose_forwardsFeesToProposer() public {
        uint64 fee = inbox.getCurrentForcedInclusionFee();
        _enqueueOneForcedInclusion();

        uint256 proposerBalanceBefore = proposer.balance;

        IRealTimeInbox.ProposeInput memory input = _buildDefaultProposeInput();
        input.numForcedInclusions = 1;
        ICheckpointStore.Checkpoint memory checkpoint = _buildCheckpoint();

        _proposeAndGetLogs(input, checkpoint);

        assertEq(
            proposer.balance,
            proposerBalanceBefore + uint256(fee) * 1 gwei,
            "proposer received FI fee"
        );
    }

    // ---------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------

    function _enqueueOneForcedInclusion() internal {
        LibBlobs.BlobReference memory ref =
            LibBlobs.BlobReference({ blobStartIndex: 0, numBlobs: 1, offset: 0 });
        uint64 fee = inbox.getCurrentForcedInclusionFee();

        _setBlobHashes(1);
        vm.prank(forcer);
        inbox.saveForcedInclusion{ value: uint256(fee) * 1 gwei }(ref);
    }
}
