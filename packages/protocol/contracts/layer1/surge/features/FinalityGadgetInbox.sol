// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { Inbox } from "../../core/impl/Inbox.sol";
import { LibHashOptimized } from "../../core/libs/LibHashOptimized.sol";
import { SurgeVerifier } from "../SurgeVerifier.sol";
import { LibProofBitmap } from "../libs/LibProofBitmap.sol";

abstract contract FinalityGadgetInbox is Inbox {
    using LibProofBitmap for LibProofBitmap.ProofBitmap;

    /// @dev Emitted when conflicting proofs are detected for a proposal
    /// @param firstProposalId The ID of the first proposal in the conflict
    /// @param conflictingProofBitmap Bitmap representing the conflicting proofs
    event ConflictingProofsDetected(
        uint48 indexed firstProposalId, LibProofBitmap.ProofBitmap conflictingProofBitmap
    );

    function proveConflicts(
        bytes calldata _commitments,
        bytes[] calldata _proofs
    )
        external
    {
        // TODO: move to optimised decoding
        Commitment[] memory commitments = abi.decode(_commitments, (Commitment[]));

        LibProofBitmap.ProofBitmap conflictingProofBitmap;

        // Multiple conflicting commitments should be provided
        require(commitments.length > 0, Surge_OnlyOneCommitmentProvided());

        for (uint256 i; i < commitments.length; ++i) {
            // Conflict checks are restricted to a single proposal for simplicity
            require(commitments[i].transitions.length == 1, Surge_MoreThanOneTransitionProvided());

            // The first commitment is expected to be the finalising commitment, while
            // the rest are expected to be conflicting
            if (i > 0) {
                // Ensure consistency between the provided commitments
                _validateCommitmentConsistency(commitments[i - 1], commitments[i]);

                // Verify proof validity and merge the flag with the bitmap
                LibProofBitmap.ProofBitmap proofBitmap = SurgeVerifier(_proofVerifier)
                    .verifyProof(false, LibHashOptimized.hashCommitment(commitments[i]), _proofs[i]);
                conflictingProofBitmap = conflictingProofBitmap.merge(proofBitmap);
            } else {
                // Set `_requireThreshold` to true to ensure this is a finalising commitment
                SurgeVerifier(_proofVerifier)
                    .verifyProof(true, LibHashOptimized.hashCommitment(commitments[i]), _proofs[i]);
            }
        }

        // Mark the conflicting verifiers upgradeable
        // Note: This fails if the verifier has already been marked upgradeable for a conflict
        // at the given proposal id
        SurgeVerifier(_proofVerifier)
            .markVerifiersUpgradeable(commitments[0].firstProposalId, conflictingProofBitmap, true);
        emit ConflictingProofsDetected(commitments[0].firstProposalId, conflictingProofBitmap);
    }

    /// @dev Validates that two commitments are consistent
    /// @param _previousCommitment The previous commitment to compare
    /// @param _currentCommitment The current commitment to compare
    function _validateCommitmentConsistency(
        Commitment memory _previousCommitment,
        Commitment memory _currentCommitment
    )
        internal
        pure
    {
        // ---------------------------------------------------------
        // Fields expected to be same between commitments
        // ---------------------------------------------------------

        require(
            _currentCommitment.firstProposalId == _previousCommitment.firstProposalId,
            Surge_FirstProposalIdMustNotDiffer()
        );
        require(
            _currentCommitment.firstProposalParentBlockHash
                == _previousCommitment.firstProposalParentBlockHash,
            Surge_FirstProposalParentBlockHashMustNotDiffer()
        );
        require(
            _currentCommitment.lastProposalHash == _previousCommitment.lastProposalHash,
            Surge_LastProposalHashMustNotDiffer()
        );
        require(
            _currentCommitment.endBlockNumber == _previousCommitment.endBlockNumber,
            Surge_EndBlockNumberMustNotDiffer()
        );
        require(
            _currentCommitment.endStateRoot == _previousCommitment.endStateRoot,
            Surge_EndStateRootMustNotDiffer()
        );

        // ---------------------------------------------------------
        // Fields expected to be different between commitments
        // ---------------------------------------------------------

        require(
            _currentCommitment.transitions[0].blockHash
                != _previousCommitment.transitions[0].blockHash,
            Surge_TransitionBlockhashMustDiffer()
        );

        // `actualProver` may or may not be different between commitments
    }

    // ---------------------------------------------------------------
    // Overrides
    // ---------------------------------------------------------------

    /// @dev Override the handler to call the SurgeVerifier that requires a proof threshold for the
    /// verification to pass.
    function _handleProofVerification(
        Commitment memory _commitment,
        bytes calldata _proof
    )
        internal
        view
        override
    {
        SurgeVerifier(_proofVerifier)
            .verifyProof(true, LibHashOptimized.hashCommitment(_commitment), _proof);
    }

    // ---------------------------------------------------------------
    // Custom errors
    // ---------------------------------------------------------------

    error Surge_EndBlockNumberMustNotDiffer();
    error Surge_EndStateRootMustNotDiffer();
    error Surge_FirstCommitmentMustBeFinalising();
    error Surge_FirstProposalIdMustNotDiffer();
    error Surge_FirstProposalParentBlockHashMustNotDiffer();
    error Surge_LastProposalHashMustNotDiffer();
    error Surge_MoreThanOneTransitionProvided();
    error Surge_OnlyOneCommitmentProvided();
    error Surge_TransitionBlockhashMustDiffer();
}
