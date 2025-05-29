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

    // --------------------------------------------------------------------------------------------
    // Happy cases
    // --------------------------------------------------------------------------------------------

    // ZK + TEE
    // --------

    function test_inbox_batch_is_finalised_immediately_with_SGX_SP1_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove using SGX + SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX_SP1, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    function test_inbox_batch_is_finalised_immediately_with_SGX_RISC0_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove using SGX + RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX_RISC0, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    // ZK followed by TEE
    // --------------------

    function test_inbox_batch_is_finalised_when_SP1_proof_is_followed_by_matching_SGX_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_SP1
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
    }

    function test_inbox_batch_is_finalised_when_RISC0_proof_is_followed_by_matching_SGX_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_RISC0
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_RISC0));
    }

    // TEE followed by ZK
    // --------------------

    function test_inbox_batch_is_finalised_when_SGX_proof_is_followed_by_matching_SP1_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_SP1
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
    }

    function test_inbox_batch_is_finalised_when_SGX_proof_is_followed_by_matching_RISC0_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_RISC0
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_RISC0));
    }

    // Misc
    // ----

    function test_inbox_sender_of_the_matching_proof_becomes_bond_receiver()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Alice proves the batch using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // Bob proves the batch using matching SGX proof type
        vm.startPrank(Bob);
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);
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

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // Proof type is set to SP1
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));

        // Prove using RISC0 proof type (i.e ZK again)
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // Proof type is still SP1, signaling that proving was skipped
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
    }

    // --------------------------------------------------------------------------------------------
    // Challenge cases
    // --------------------------------------------------------------------------------------------

    // ZK or ZK + TEE Challenging a ZK proof
    // -------------------------------------

    function test_inbox_an_SP1_proof_can_be_challenged_by_a_RISC0_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using RISC0 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.challenged, true);
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SP1));
        // and, a new cooldown period begins
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
    }

    function test_inbox_a_RISC0_proof_can_be_challenged_by_an_SP1_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SP1 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.challenged, true);
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.RISC0));
        // and, a new cooldown period begins
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
    }

    function test_inbox_batch_is_finalised_when_SP1_proof_is_challenged_by_an_SGX_RISC0_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX_RISC0 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_RISC0, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_RISC0));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SP1));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, batch is finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    function test_inbox_batch_is_finalised_when_RISC0_proof_is_challenged_by_an_SGX_SP1_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX_SP1 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_SP1, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        console2.log("Here");
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.RISC0));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, batch is finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    function test_inbox_batch_is_finalised_when_challenging_SP1_proof_gets_matching_SGX_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        vm.warp(block.timestamp + 2);

        // Challenge using SP1 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The batch is challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // but not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove challenged transitions using SGX proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_SP1
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, challengedBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
    }

    function test_inbox_batch_is_finalised_when_challenging_RISC0_proof_gets_matching_SGX_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        vm.warp(block.timestamp + 2);

        // Challenge using RISC0 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The batch is challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // but not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove challenged transitions using SGX proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_RISC0
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, challengedBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_RISC0));
    }

    function test_inbox_sender_becomes_bond_receiver_when_SP1_proof_is_challenged_by_an_SGX_RISC0_proof(
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX_RISC0 proof type
        vm.startPrank(Bob);
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_RISC0, batchIds);
        vm.stopPrank();

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_RISC0));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, bond receiver is updated to Bob
        assertEq(ts.bondReceiver, Bob);
    }

    function test_inbox_sender_becomes_bond_receiver_when_RISC0_proof_is_challenged_by_an_SGX_SP1_proof(
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX_SP1 proof type
        vm.startPrank(Bob);
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_SP1, batchIds);
        vm.stopPrank();

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, bond receiver is updated to Bob
        assertEq(ts.bondReceiver, Bob);
    }

    // ZK cannot be challenged by TEE
    // ------------------------------

    function test_inbox_an_SP1_proof_cannot_be_challenged_by_an_SGX_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Attempt to challenge using SGX proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The transition should remain unchanged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp - 2);
        assertEq(ts.blockHash, correctBlockhash(1));
    }

    function test_inbox_an_RISC0_proof_cannot_be_challenged_by_an_SGX_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Attempt to challenge using SGX proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The transition should remain unchanged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp - 2);
        assertEq(ts.blockHash, correctBlockhash(1));
    }

    // Challenging a TEE proof
    // ------------------------

    // Note: This should become SGX <> TDX whenever that is introduced
    function test_inbox_an_SGX_proof_can_be_challenged_by_another_SGX_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX proof type again
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, true);
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SGX));
        // and, a new cooldown period begins
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
    }

    function test_inbox_an_SGX_proof_can_be_challenged_by_an_SP1_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SP1 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
        assertEq(ts.challenged, true);
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SGX));
        // and, a new cooldown period begins
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
    }

    function test_inbox_an_SGX_proof_can_be_challenged_by_a_RISC0_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using RISC0 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
        assertEq(ts.challenged, true);
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SGX));
        // and, a new cooldown period begins
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
    }

    // Note: This should become SGX <> TDX + ZK whenever TDX is introduced
    function test_inbox_batch_is_finalised_when_SGX_proof_is_challenged_by_an_SGX_SP1_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX + SP1 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_SP1, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, batch is finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    function test_inbox_batch_is_finalised_when_SGX_proof_is_challenged_by_an_SGX_RISC0_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX + RISC0 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_RISC0, batchIds);

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_RISC0));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, batch is finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
    }

    function test_inbox_batch_is_finalised_when_challenging_SGX_proof_gets_matching_SP1_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Note: This will be changed to TDX later
        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // but not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove challenged transitions using SP1 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_SP1
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, challengedBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
    }

    function test_inbox_batch_is_finalised_when_challenging_SGX_proof_gets_matching_RISC0_proof()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Note: This will be changed to TDX later
        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // but not finalised yet
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);

        // Prove challenged transitions using RISC0 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The batch is now finalised
        stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_RISC0
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, challengedBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_RISC0));
    }

    function test_inbox_sender_becomes_bond_receiver_when_SGX_proof_is_challenged_by_an_SGX_SP1_proof(
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX + SP1 proof type
        vm.startPrank(Bob);
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_SP1, batchIds);
        vm.stopPrank();

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, bond receiver is updated to Bob
        assertEq(ts.bondReceiver, Bob);
    }

    function test_inbox_sender_becomes_bond_receiver_when_SGX_proof_is_challenged_by_an_SGX_RISC0_proof(
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, false);
        assertEq(ts.createdAt, block.timestamp);

        vm.warp(block.timestamp + 2);

        // Challenge using SGX + RISC0 proof type
        vm.startPrank(Bob);
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_RISC0, batchIds);
        vm.stopPrank();

        // The transition is now challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_RISC0));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SGX));
        assertEq(ts.challenged, true);
        assertEq(ts.createdAt, block.timestamp);
        assertEq(ts.blockHash, challengedBlockhash(1));
        // and, bond receiver is updated to Bob
        assertEq(ts.bondReceiver, Bob);
    }

    // ----------------------------------------------------------
    // Cooldown Period
    // ----------------------------------------------------------

    function test_inbox_batch_is_finalised_when_SGX_proof_is_not_challenged_within_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
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
        // and, proof type remains SGX as it was not challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
    }

    function test_inbox_batch_is_finalised_when_SP1_proof_is_not_challenged_within_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
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
        // and, proof type remains SP1 as it was not challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
    }

    function test_inbox_batch_is_finalised_when_RISC0_proof_is_not_challenged_within_cooldown_period(
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // The batch is not challenged yet
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
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
        // and, proof type remains RISC0 as it was not challenged
        ts = inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
    }

    function test_inbox_dao_receives_liveness_bond_when_SGX_proof_is_finalised_via_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow + 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should now be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type remains SGX as it was not challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX));
        // and, liveness bond is sent to DAO
        assertEq(
            inbox.bondBalanceOf(TaikoInbox(address(inbox)).dao()), pacayaConfig().livenessBondBase
        );
    }

    function test_inbox_dao_receives_liveness_bond_when_SP1_proof_is_finalised_via_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow + 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should now be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type remains SP1 as it was not challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SP1));
        // and, liveness bond is sent to DAO
        assertEq(
            inbox.bondBalanceOf(TaikoInbox(address(inbox)).dao()), pacayaConfig().livenessBondBase
        );
    }

    function test_inbox_dao_receives_liveness_bond_when_RISC0_proof_is_finalised_via_cooldown_period(
    )
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow + 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should now be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type remains RISC0 as it was not challenged
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(ts.blockHash, correctBlockhash(1));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.RISC0));
        // and, liveness bond is sent to DAO
        assertEq(
            inbox.bondBalanceOf(TaikoInbox(address(inbox)).dao()), pacayaConfig().livenessBondBase
        );
    }

    function test_inbox_challenged_SGX_proof_cannot_be_finalised_via_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        // Note: This will be changed to TDX later
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // Challenge using SGX proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX, batchIds);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow + 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should still not be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
    }

    function test_inbox_challenged_SP1_proof_cannot_be_finalised_via_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // Challenge using SP1 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SP1, batchIds);

        // Warp time to after the cooldown period ends
        vm.warp(block.timestamp + pacayaConfig().cooldownWindow + 1);

        // Attempt to finalise
        inbox.verifyBatches(1);

        // The batch should still not be finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 0);
    }

    function test_inbox_challenged_RISC0_proof_cannot_be_finalised_via_cooldown_period()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // Challenge using RISC0 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.RISC0, batchIds);

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

    function test_inbox_verifier_is_upgradeable_when_challenged_SGX_proof_is_finalised()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SGX proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SGX, batchIds);

        // Challenge using SGX + SP1 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_SP1, batchIds);

        // The batch is now finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_SP1
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SGX));
        // and, SGX verifier is upgradeable
        assertEq(uint8(verifier.proofTypeToUpgrade()), uint8(LibProofType.ProofType.SGX));
    }

    function test_inbox_verifier_is_upgradeable_when_challenged_SP1_proof_is_finalised()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using SP1 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.SP1, batchIds);

        // Challenge using SP1 + RISC0 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_RISC0, batchIds);

        // The batch is now finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_RISC0
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_RISC0));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.SP1));
        // and, SP1 verifier is upgradeable
        assertEq(uint8(verifier.proofTypeToUpgrade()), uint8(LibProofType.ProofType.SP1));
    }

    function test_inbox_verifier_is_upgradeable_when_challenged_RISC0_proof_is_finalised()
        external
        transactBy(Alice)
        WhenMultipleBatchesAreProposedWithDefaultParameters(1)
    {
        uint64[] memory batchIds = new uint64[](1);
        batchIds[0] = 1;

        // Prove using RISC0 proof type
        _proveBatchesWithProofType(LibProofType.ProofType.RISC0, batchIds);

        // Challenge using SGX + SP1 proof type
        _challengeTransitionWithProofType(LibProofType.ProofType.SGX_SP1, batchIds);

        // The batch is now finalised
        ITaikoInbox.Stats2 memory stats2 = inbox.getStats2();
        assertEq(stats2.lastVerifiedBatchId, 1);
        // and, proof type is updated to SGX_SP1
        ITaikoInbox.TransitionState memory ts =
            inbox.getTransitionByParentHash(1, correctBlockhash(0));
        assertEq(uint8(ts.proofType), uint8(LibProofType.ProofType.SGX_SP1));
        assertEq(uint8(ts.challengedProofType), uint8(LibProofType.ProofType.RISC0));
        // and, RISC0 verifier is upgradeable
        assertEq(uint8(verifier.proofTypeToUpgrade()), uint8(LibProofType.ProofType.RISC0));
    }
}
