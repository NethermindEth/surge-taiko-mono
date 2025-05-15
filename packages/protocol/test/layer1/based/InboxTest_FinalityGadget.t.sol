// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.24;

import "./InboxTestBase.sol";

contract InboxTest_FinalityGadget is InboxTestBase {
    function pacayaConfig() internal pure override returns (ITaikoInbox.Config memory) {
        ITaikoInbox.ForkHeights memory forkHeights;

        return ITaikoInbox.Config({
            chainId: LibNetwork.TAIKO_MAINNET,
            maxUnverifiedBatches: 10,
            batchRingBufferSize: 11,
            maxBatchesToVerify: 5,
            blockMaxGasLimit: 240_000_000,
            livenessBondBase: 125e18, // 125 Taiko token per batch
            livenessBondPerBlock: 0, // deprecated
            stateRootSyncInternal: 5,
            maxAnchorHeightOffset: 64,
            baseFeeConfig: LibSharedData.BaseFeeConfig({
                adjustmentQuotient: 8,
                sharingPctg: 75,
                gasIssuancePerSecond: 5_000_000,
                minGasExcess: 1_340_000_000, // correspond to 0.008847185 gwei basefee
                maxGasIssuancePerBlock: 600_000_000 // two minutes: 5_000_000 * 120
             }),
            provingWindow: 1 hours,
            cooldownWindow: 7 days,
            maxSignalsToReceive: 16,
            maxBlocksPerBatch: 768,
            forkHeights: forkHeights,
            // Surge: to prevent compilation errors
            maxVerificationDelay: 0
        });
    }

    function setUpOnEthereum() internal override {
        bondToken = deployBondToken();
        super.setUpOnEthereum();
    }

    function test_inbox_batch_is_finalised_immediately_with_ZK_TEE_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove using ZK_TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK_TEE, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    function test_inbox_batch_is_finalised_when_ZK_proof_is_followed_by_matching_TEE_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to ZK_TEE
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
    }

    function test_inbox_batch_is_finalised_when_TEE_proof_is_followed_by_matching_ZK_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to ZK_TEE
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
    }

    function test_inbox_sender_of_the_matching_proof_becomes_bond_receiver()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Alice proves the batch using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // Bob proves the batch using matching TEE proof type
        vm.startPrank(Bob);
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);
        vm.stopPrank();

        // The batch is now finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, bond receiver is updated to Bob
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.bondReceiver, Bob);
    }

    function test_inbox_skips_reproving_transition_when_existing_proof_type_matches_new_proof_type()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // Proof type is set to ZK
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));

        // Prove using ZK proof type again
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // Proof type is still ZK, signaling that proving was skipped
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
    }

    function test_inbox_a_ZK_proof_can_be_challenged_by_another_ZK_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using ZK proof type again
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
    }

    function test_inbox_batch_is_finalised_when_ZK_proof_is_challenged_by_a_ZK_TEE_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using ZK_TEE proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK_TEE, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, batch is finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    function test_inbox_sender_becomes_bond_receiver_when_ZK_proof_is_challenged_by_a_ZK_TEE_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using ZK_TEE proof type
        vm.startPrank(Bob);
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK_TEE, batchIds);
        vm.stopPrank();

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, bond receiver is updated to Bob
        assertEq(ts.bondReceiver, Bob);
    }

    function test_inbox_a_ZK_proof_cannot_be_challenged_by_a_TEE_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Attempt to challenge using TEE proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The transition should remain unchanged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp - 2);
        assertEq(ts.blockHash, correctBlockhash(1));
    }

    function test_inbox_a_TEE_proof_can_be_challenged_by_another_TEE_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using TEE proof type again
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
    }

    function test_inbox_a_TEE_proof_can_be_challenged_by_a_ZK_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using ZK proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
    }

    function test_inbox_batch_is_finalised_when_TEE_proof_is_challenged_by_a_ZK_TEE_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using ZK_TEE proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK_TEE, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, batch is finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    function test_inbox_sender_becomes_bond_receiver_when_TEE_proof_is_challenged_by_a_ZK_TEE_proof(
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using ZK_TEE proof type
        vm.startPrank(Bob);
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK_TEE, batchIds);
        vm.stopPrank();

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, bond receiver is updated to Bob
        assertEq(ts.bondReceiver, Bob);
    }

    function test_inbox_batch_is_finalised_when_challenged_ZK_proof_gets_matching_TEE_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        vm.warp(block.timestamp + 2);

        // Challenge using ZK proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // but not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove challenged transitions using TEE proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to ZK_TEE
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, challengedBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
    }

    function test_inbox_batch_is_finalised_when_challenged_TEE_proof_gets_matching_ZK_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        vm.warp(block.timestamp + 2);

        // Challenge using TEE proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // but not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove challenged transitions using ZK proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to ZK_TEE
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, challengedBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
    }

    function test_inbox_batch_is_finalised_when_ZK_proof_is_not_challenged_within_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
        assertEq(ts.challenged, false);
        assertEq(ts.blockHash, correctBlockhash(1));

        // Warp time to just before the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow - 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should still not be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + 2);

        // Attempt to finalise again
        inbox.verifyBatches(1);

        // The batch should now be finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type remains ZK as it was not challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
    }

    function test_inbox_batch_is_finalised_when_TEE_proof_is_not_challenged_within_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));
        assertEq(ts.challenged, false);
        assertEq(ts.blockHash, correctBlockhash(1));

        // Warp time to just before the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow - 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should still not be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + 2);

        // Attempt to finalise again
        inbox.verifyBatches(1);

        // The batch should now be finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type remains TEE as it was not challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));
    }

    function test_inbox_dao_receives_liveness_bond_when_ZK_proof_is_finalised_via_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow + 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should now be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type remains ZK as it was not challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK));
        // and, liveness bond is sent to DAO
        assertEq(
            inbox.bondBalanceOf(TaikoInbox(address(inbox)).dao()), pacayaConfig().livenessBondBase
        );
    }

    function test_inbox_dao_receives_liveness_bond_when_TEE_proof_is_finalised_via_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow + 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should now be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type remains TEE as it was not challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.TEE));
        // and, liveness bond is sent to DAO
        assertEq(
            inbox.bondBalanceOf(TaikoInbox(address(inbox)).dao()), pacayaConfig().livenessBondBase
        );
    }

    function test_inbox_verifier_is_upgradeable_when_challenged_ZK_proof_is_finalised()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // Challenge using ZK proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // Prove challenged transitions using TEE proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // The batch is now finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to ZK_TEE
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
        // and, verifier is upgradeable
        assertEq(inbox.getVerifier().upgradeable, true);
    }

    function test_inbox_verifier_is_upgradeable_when_challenged_TEE_proof_is_finalised()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // Challenge using TEE proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.TEE, batchIds);

        // Prove challenged transitions using ZK proof type
        _challengeTransitionWithProofType(ITaikoInbox.ProofType.ZK, batchIds);

        // The batch is now finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to ZK_TEE
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(ITaikoInbox.ProofType.ZK_TEE));
        // and, verifier is upgradeable
        assertEq(inbox.getVerifier().upgradeable, true);
    }
}
