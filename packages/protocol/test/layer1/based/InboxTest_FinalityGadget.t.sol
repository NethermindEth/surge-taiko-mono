// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.24;

import "./InboxTestBase.sol";
import "./helpers/ProofTypeFixtures.sol";

contract InboxTest_FinalityGadget is InboxTestBase, ProofTypeFixtures {
    using LibProofType for LibProofType.ProofType;

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

    // --------------------------------------------------------------------------------------------
    // Happy cases
    // --------------------------------------------------------------------------------------------

    // ZK + TEE
    // --------

    function test_inbox_batch_is_finalised_immediately_with_ZK_TEE_proof(uint256 _zkTeeIndex)
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkTeeIndex = bound(_zkTeeIndex, 0, zkTeeProofTypes.length - 1);
        LibProofType.ProofType zkTeeProofType = zkTeeProofTypes[_zkTeeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove using ZK + TEE proof type
        _proveBatchesWithProofType(zkTeeProofType, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    // ZK followed by TEE
    // --------------------

    function test_inbox_batch_is_finalised_when_ZK_proof_is_followed_by_matching_TEE_proof(
        uint256 _zkIndex,
        uint256 _teeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(zkProofType, batchIds);

        // The batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(zkProofType));

        // Prove using TEE proof type
        _proveBatchesWithProofType(teeProofType, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to ZK + TEE
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(zkProofType.combine(teeProofType)));
    }

    // TEE followed by ZK
    // --------------------

    function test_inbox_batch_is_finalised_when_TEE_proof_is_followed_by_matching_ZK_proof(
        uint256 _teeIndex,
        uint256 _zkIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(teeProofType, batchIds);

        // The batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(teeProofType));

        // Prove using ZK proof type
        _proveBatchesWithProofType(zkProofType, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to ZK + TEE
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(teeProofType.combine(zkProofType)));
    }

    // Misc
    // ----

    function test_inbox_sender_of_the_matching_proof_becomes_bond_receiver(
        uint256 _zkIndex,
        uint256 _teeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Alice proves the batch using ZK proof type
        _proveBatchesWithProofType(zkProofType, batchIds);

        // Bob proves the batch using matching TEE proof type
        vm.startPrank(Bob);
        _proveBatchesWithProofType(teeProofType, batchIds);
        vm.stopPrank();

        // The batch is now finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, bond receiver is updated to Bob
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.bondReceiver, Bob);
    }

    function test_inbox_skips_reproving_transition_when_both_existing_and_new_proof_types_are_ZK(
        uint256 _zkIndex1,
        uint256 _zkIndex2
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex1 = bound(_zkIndex1, 0, zkProofTypes.length - 1);
        _zkIndex2 = bound(_zkIndex2, 0, zkProofTypes.length - 1);

        vm.assume(_zkIndex1 != _zkIndex2);

        LibProofType.ProofType zkProofType1 = zkProofTypes[_zkIndex1];
        LibProofType.ProofType zkProofType2 = zkProofTypes[_zkIndex2];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type 1
        _proveBatchesWithProofType(zkProofType1, batchIds);

        // Proof type is set to ZK proof type 1
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(zkProofType1));

        // Prove using ZK proof type 2
        _proveBatchesWithProofType(zkProofType2, batchIds);

        // Proof type is still ZK proof type 1, signaling that proving was skipped
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(zkProofType1));
    }

    function test_inbox_skips_reproving_transition_when_both_existing_and_new_proof_types_are_TEE(
        uint256 _teeIndex1,
        uint256 _teeIndex2
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _teeIndex1 = bound(_teeIndex1, 0, teeProofTypes.length - 1);
        _teeIndex2 = bound(_teeIndex2, 0, teeProofTypes.length - 1);

        vm.assume(_teeIndex1 != _teeIndex2);

        LibProofType.ProofType teeProofType1 = teeProofTypes[_teeIndex1];
        LibProofType.ProofType teeProofType2 = teeProofTypes[_teeIndex2];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type 1
        _proveBatchesWithProofType(teeProofType1, batchIds);

        // Proof type is set to TEE proof type 1
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(teeProofType1));

        // Prove using TEE proof type 2
        _proveBatchesWithProofType(teeProofType2, batchIds);

        // Proof type is still TEE proof type 1, signaling that proving was skipped
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(teeProofType1));
    }

    // --------------------------------------------------------------------------------------------
    // Conflicting proof cases
    // --------------------------------------------------------------------------------------------

    // Conflicts with existing ZK proof
    // --------------------------------

    function test_inbox_push_conflicting_ZK_proof_for_existing_ZK_proof(
        uint256 _zkIndex1,
        uint256 _zkIndex2
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex1 = bound(_zkIndex1, 0, zkProofTypes.length - 1);
        _zkIndex2 = bound(_zkIndex2, 0, zkProofTypes.length - 1);
        LibProofType.ProofType zkProofType1 = zkProofTypes[_zkIndex1];
        LibProofType.ProofType zkProofType2 = zkProofTypes[_zkIndex2];

        vm.assume(_zkIndex1 != _zkIndex2);

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type 1
        _proveBatchesWithProofType(zkProofType1, batchIds);

        // The transition has no conflicting proofs yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(zkProofType1));
        assertEq(ts.numConflictingProofs, 0);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Push a conflicting ZK proof
        _pushConflictingProof(zkProofType2, batchIds);

        // The transition now has a conflicting proof
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 1);
        assertEq(ts.blockHashes[1], conflictingBlockHash(1));
        assertTrue(ts.proofTypes[1].equals(zkProofType2));
    }

    function test_inbox_batch_is_finalised_when_conflicting_ZK_TEE_proof_is_pushed_for_existing_ZK_proof(
        uint256 _zkIndex,
        uint256 _zkTeeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        _zkTeeIndex = bound(_zkTeeIndex, 0, zkTeeProofTypes.length - 1);
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];
        LibProofType.ProofType zkTeeProofType = zkTeeProofTypes[_zkTeeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(zkProofType, batchIds);

        // The transition has no conflicting proofs yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(zkProofType));
        assertEq(ts.numConflictingProofs, 0);

        vm.warp(block.timestamp + 2);

        // Push a conflicting ZK + TEE proof
        _pushConflictingProof(zkTeeProofType, batchIds);

        // The transition now has a conflicting proof
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 1);
        assertEq(ts.blockHashes[1], conflictingBlockHash(1));
        assertTrue(ts.proofTypes[1].equals(zkTeeProofType));

        // The batch is now finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, finalising proof index is updated to 1
        ITaikoInbox.Batch memory batch = inbox.getBatch(1);
        assertEq(batch.finalisingProofIndex, 1);
    }

    function test_inbox_batch_is_finalised_when_conflicting_ZK_proof_gets_matching_TEE_proof(
        uint256 _zkIndex1,
        uint256 _zkIndex2,
        uint256 _teeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex1 = bound(_zkIndex1, 0, zkProofTypes.length - 1);
        _zkIndex2 = bound(_zkIndex2, 0, zkProofTypes.length - 1);
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        LibProofType.ProofType zkProofType1 = zkProofTypes[_zkIndex1];
        LibProofType.ProofType zkProofType2 = zkProofTypes[_zkIndex2];
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];

        vm.assume(_zkIndex1 != _zkIndex2);

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type 1
        _proveBatchesWithProofType(zkProofType1, batchIds);

        vm.warp(block.timestamp + 2);

        // Push a conflicting ZK proof type 2
        _pushConflictingProof(zkProofType2, batchIds);

        // The transition now has a conflicting proof
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 1);
        assertEq(ts.blockHashes[1], conflictingBlockHash(1));
        assertTrue(ts.proofTypes[1].equals(zkProofType2));
        // But the batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Push a matching TEE proof for the conflicting ZK proof type 2
        _pushConflictingProof(teeProofType, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, finalising proof type is updated to ZK + TEE
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[1].equals(zkProofType2.combine(teeProofType)));
        // and, finalising proof index is updated to 1
        ITaikoInbox.Batch memory batch = inbox.getBatch(1);
        assertEq(batch.finalisingProofIndex, 1);
    }

    function test_inbox_sender_becomes_bond_receiver_when_conflicting_ZK_TEE_proof_is_pushed_for_existing_ZK_proof(
        uint256 _zkIndex,
        uint256 _zkTeeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        _zkTeeIndex = bound(_zkTeeIndex, 0, zkTeeProofTypes.length - 1);
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];
        LibProofType.ProofType zkTeeProofType = zkTeeProofTypes[_zkTeeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(zkProofType, batchIds);

        // The transition has no conflicting proofs yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(zkProofType));
        assertEq(ts.numConflictingProofs, 0);

        vm.warp(block.timestamp + 2);

        // Push a conflicting ZK + TEE proof
        vm.startPrank(Bob);
        _pushConflictingProof(zkTeeProofType, batchIds);
        vm.stopPrank();

        // The transition now has a conflicting proof
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 1);
        assertEq(ts.blockHashes[1], conflictingBlockHash(1));
        assertTrue(ts.proofTypes[1].equals(zkTeeProofType));
        // and, bond receiver is updated to Bob
        assertEq(ts.bondReceiver, Bob);
    }

    // Conflicts with existing TEE proof
    // ----------------------------------

    function test_inbox_push_conflicting_TEE_proof_for_existing_TEE_proof(
        uint256 _teeIndex1,
        uint256 _teeIndex2
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _teeIndex1 = bound(_teeIndex1, 0, teeProofTypes.length - 1);
        _teeIndex2 = bound(_teeIndex2, 0, teeProofTypes.length - 1);
        LibProofType.ProofType teeProofType1 = teeProofTypes[_teeIndex1];
        LibProofType.ProofType teeProofType2 = teeProofTypes[_teeIndex2];

        vm.assume(_teeIndex1 != _teeIndex2);

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type 1
        _proveBatchesWithProofType(teeProofType1, batchIds);

        // The transition has no conflicting proofs yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(teeProofType1));
        assertEq(ts.numConflictingProofs, 0);

        vm.warp(block.timestamp + 2);

        // Push a conflicting TEE proof type 2
        _pushConflictingProof(teeProofType2, batchIds);

        // The transition now has a conflicting proof
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 1);
        assertEq(ts.blockHashes[1], conflictingBlockHash(1));
        assertTrue(ts.proofTypes[1].equals(teeProofType2));
    }

    function test_inbox_push_conflicting_ZK_proof_for_existing_TEE_proof(
        uint256 _teeIndex,
        uint256 _zkIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(teeProofType, batchIds);

        // The transition has no conflicting proofs yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(teeProofType));
        assertEq(ts.numConflictingProofs, 0);

        vm.warp(block.timestamp + 2);

        // Push a conflicting ZK proof
        _pushConflictingProof(zkProofType, batchIds);

        // The transition now has a conflicting proof
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 1);
        assertEq(ts.blockHashes[1], conflictingBlockHash(1));
        assertTrue(ts.proofTypes[1].equals(zkProofType));
    }

    function test_inbox_batch_is_finalised_when_conflicting_ZK_TEE_proof_is_pushed_for_existing_TEE_proof(
        uint256 _teeIndex,
        uint256 _zkTeeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        _zkTeeIndex = bound(_zkTeeIndex, 0, zkTeeProofTypes.length - 1);
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];
        LibProofType.ProofType zkTeeProofType = zkTeeProofTypes[_zkTeeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(teeProofType, batchIds);

        // The transition has no conflicting proofs yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(teeProofType));
        assertEq(ts.numConflictingProofs, 0);

        vm.warp(block.timestamp + 2);

        // Push a conflicting ZK + TEE proof
        _pushConflictingProof(zkTeeProofType, batchIds);

        // The transition now has a conflicting proof
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 1);
        assertEq(ts.blockHashes[1], conflictingBlockHash(1));
        assertTrue(ts.proofTypes[1].equals(zkTeeProofType));
        // and, batch is finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, finalising proof index is updated to 1
        ITaikoInbox.Batch memory batch = inbox.getBatch(1);
        assertEq(batch.finalisingProofIndex, 1);
    }

    function test_inbox_batch_is_finalised_when_conflicting_TEE_proof_gets_matching_ZK_proof(
        uint256 _teeIndex1,
        uint256 _teeIndex2,
        uint256 _zkIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _teeIndex1 = bound(_teeIndex1, 0, teeProofTypes.length - 1);
        _teeIndex2 = bound(_teeIndex2, 0, teeProofTypes.length - 1);
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        LibProofType.ProofType teeProofType1 = teeProofTypes[_teeIndex1];
        LibProofType.ProofType teeProofType2 = teeProofTypes[_teeIndex2];
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];

        vm.assume(_teeIndex1 != _teeIndex2);

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type 1
        _proveBatchesWithProofType(teeProofType1, batchIds);

        vm.warp(block.timestamp + 2);

        // Push a conflicting TEE proof type 2
        _pushConflictingProof(teeProofType2, batchIds);

        // The transition now has a conflicting proof
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 1);
        assertEq(ts.blockHashes[1], conflictingBlockHash(1));
        assertTrue(ts.proofTypes[1].equals(teeProofType2));
        // but not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Push a matching ZK proof for the conflicting TEE proof type 2
        _pushConflictingProof(zkProofType, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, finalising proof type is updated to ZK + TEE
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[1].equals(zkProofType.combine(teeProofType2)));
        // and, finalising proof index is updated to 1
        ITaikoInbox.Batch memory batch = inbox.getBatch(1);
        assertEq(batch.finalisingProofIndex, 1);
    }

    function test_inbox_sender_becomes_bond_receiver_when_conflicting_ZK_TEE_proof_is_pushed_for_existing_TEE_proof(
        uint256 _teeIndex,
        uint256 _zkTeeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        _zkTeeIndex = bound(_zkTeeIndex, 0, zkTeeProofTypes.length - 1);
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];
        LibProofType.ProofType zkTeeProofType = zkTeeProofTypes[_zkTeeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(teeProofType, batchIds);

        // The transition has no conflicting proofs yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(teeProofType));
        assertEq(ts.numConflictingProofs, 0);

        vm.warp(block.timestamp + 2);

        // Push a conflicting ZK + TEE proof
        vm.startPrank(Bob);
        _pushConflictingProof(zkTeeProofType, batchIds);
        vm.stopPrank();

        // The transition now has a conflicting proof
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 1);
        assertEq(ts.blockHashes[1], conflictingBlockHash(1));
        assertTrue(ts.proofTypes[1].equals(zkTeeProofType));
        // and, bond receiver is updated to Bob
        assertEq(ts.bondReceiver, Bob);
        // and, finalising proof index is updated to 1
        ITaikoInbox.Batch memory batch = inbox.getBatch(1);
        assertEq(batch.finalisingProofIndex, 1);
    }

    // ----------------------------------------------------------
    // Cooldown Period
    // ----------------------------------------------------------

    function test_inbox_batch_is_finalised_when_existing_ZK_proof_has_no_conflicts_within_cooldown_period(
        uint256 _zkIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(zkProofType, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(zkProofType));
        assertEq(ts.numConflictingProofs, 0);

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
        assertTrue(ts.proofTypes[0].equals(zkProofType));
        assertEq(ts.numConflictingProofs, 0);
    }

    function test_inbox_batch_is_finalised_when_existing_TEE_proof_has_no_conflicts_within_cooldown_period(
        uint256 _teeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(teeProofType, batchIds);

        // The transition has no conflicting proofs yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertTrue(ts.proofTypes[0].equals(teeProofType));
        assertEq(ts.numConflictingProofs, 0);

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
        assertTrue(ts.proofTypes[0].equals(teeProofType));
        assertEq(ts.numConflictingProofs, 0);
    }

    function test_inbox_dao_receives_liveness_bond_when_ZK_proof_is_finalised_via_cooldown_period(
        uint256 _zkIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using ZK proof type
        _proveBatchesWithProofType(zkProofType, batchIds);

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
        assertTrue(ts.proofTypes[0].equals(zkProofType));
        assertEq(ts.numConflictingProofs, 0);
        // and, liveness bond is sent to DAO
        assertEq(
            inbox.bondBalanceOf(TaikoInbox(address(inbox)).dao()), pacayaConfig().livenessBondBase
        );
    }

    function test_inbox_dao_receives_liveness_bond_when_TEE_proof_is_finalised_via_cooldown_period(
        uint256 _teeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(teeProofType, batchIds);

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
        assertTrue(ts.proofTypes[0].equals(teeProofType));
        assertEq(ts.numConflictingProofs, 0);
        // and, liveness bond is sent to DAO
        assertEq(
            inbox.bondBalanceOf(TaikoInbox(address(inbox)).dao()), pacayaConfig().livenessBondBase
        );
    }

    function test_inbox_batch_cannot_be_finalised_via_cooldown_period_if_there_are_conflicting_proofs(
        uint256 _zkIndex,
        uint256 _teeIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        LibProofType.ProofType zkProofType = zkProofTypes[_zkIndex];
        LibProofType.ProofType teeProofType = teeProofTypes[_teeIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using TEE proof type
        _proveBatchesWithProofType(teeProofType, batchIds);

        // Push a conflicting ZK proof
        _pushConflictingProof(zkProofType, batchIds);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow + 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should still not be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
    }

    // ----------------------------------------------------------
    // Verifier Upgradeability
    // ----------------------------------------------------------

    function test_inbox_verifiers_of_conflicting_proof_are_marked_for_upgrade(
        uint256 _zkTeeIndex,
        uint256 _teeIndex,
        uint256 _zkIndex
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        _zkTeeIndex = bound(_zkTeeIndex, 0, zkTeeProofTypes.length - 1);
        _teeIndex = bound(_teeIndex, 0, teeProofTypes.length - 1);
        _zkIndex = bound(_zkIndex, 0, zkProofTypes.length - 1);
        LibProofType.ProofType zkTeeProofType = zkTeeProofTypes[_zkTeeIndex];
        LibProofType.ProofType teeConflictingProofType = teeProofTypes[_teeIndex];
        LibProofType.ProofType zkConflictingProofType = zkProofTypes[_zkIndex];

        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Push a conflicting TEE proof with salt 1
        _pushConflictingProof(teeConflictingProofType, batchIds, 1);

        // Push a conflicting ZK proof with salt 2
        _pushConflictingProof(zkConflictingProofType, batchIds, 2);

        // Push finalising proof
        _proveBatchesWithProofType(zkTeeProofType, batchIds);

        // The batch is now finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // but, it contains 2 conflicting proofs
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.numConflictingProofs, 2);
        assertTrue(ts.proofTypes[0].equals(teeConflictingProofType));
        assertTrue(ts.proofTypes[1].equals(zkConflictingProofType));
        // and, the finalising proof
        assertTrue(ts.proofTypes[2].equals(zkTeeProofType));
        // Finalising proof index is updated to 2
        ITaikoInbox.Batch memory batch = inbox.getBatch(1);
        assertEq(batch.finalisingProofIndex, 2);
        // and, conflicting ZK + conflicting TEE verifier is upgradeable
        assertTrue(
            verifier.proofTypeToUpgrade().equals(
                zkConflictingProofType.combine(teeConflictingProofType)
            )
        );
    }
}
