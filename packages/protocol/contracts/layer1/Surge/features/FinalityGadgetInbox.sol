// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { Inbox } from "../../core/impl/Inbox.sol";
import { LibBondInstruction } from "../../core/libs/LibBondInstruction.sol";
import { LibProofBitmap } from "../LibProofBitmap.sol";
import { SurgeVerifier } from "../SurgeVerifier.sol";
import { LibBonds } from "src/shared/libs/LibBonds.sol";

/// @title FinalityGadgetInbox
/// @dev An extension of `InboxOptimized1` that accomodates the optimisitc fallback finality gadget.
/// - Transition records are now wrapped in a `SurgeTransitionRecord` that also contains the proof
/// variants.
/// - To accomodate conflicts, the tuple (proposalId, parentTransitionHash) maps to the hash of
/// an array of Surge transition records.
/// - Except for _buildAndSaveTransitionRecords, no other internal function related to storing or
/// managing a transition record is used.
/// @custom:security-contact security@nethermind.io
abstract contract FinalityGadgetInbox is Inbox {
    using LibProofBitmap for LibProofBitmap.ProofBitmap;

    /// @dev Delay before which a single proof with no conflicts can be used for finalising
    /// a transition
    uint256 internal immutable _optimisticFallbackDelay;

    /// @dev Minimum number of distinct proofs required for a transition to be finalising
    uint8 internal immutable _finalisingProofCount;

    /// @param conflictingProofBitmap A proof bitmap of all proofs conflicting with the
    /// finalising proof
    event ConflictingProofsDetected(LibProofBitmap.ProofBitmap conflictingProofBitmap);

    /// @dev A wrapper around `TransitionRecordHashAndDeadline` with an added `proofBitmap`
    /// that stores the different proof variants used for the transition
    struct SurgeTransitionRecord {
        LibProofBitmap.ProofBitmap proofBitmap;
        TransitionRecordHashAndDeadline hashAndDeadline;
    }

    /// @dev Synonymous with `InboxOptimized1.ReusableTransitionRecord`.
    struct ReusableSurgeTransitionRecord {
        uint48 proposalId;
        bytes26 partialParentTransitionHash;
        /// @dev Keccak(SurgeTransitionRecord[])
        /// @dev Ideally, we may only have a single transition record
        bytes32 surgeTransitionRecordsHash;
    }

    /// @dev This overrides the mapping `_transitionRecordAndHashDeadline` from `Inbox`.
    mapping(bytes32 compositKey => bytes32 surgeTransitionRecordsHash) internal
        _surgeTransitionRecordsHashes;

    /// @dev This is comparable to the mapping `_resuableTransitionRecords` from `InboxOptimized1`.
    mapping(uint256 bufferSlot => ReusableSurgeTransitionRecord reusableRecord) internal
        _surgeTransitionRecordsHashesBuffer;

    uint256[48] private __gap;

    // ---------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------

    constructor(uint256 __optimisticFallbackDelay, uint8 __finalisingProofCount) {
        _optimisticFallbackDelay = __optimisticFallbackDelay;
        _finalisingProofCount = __finalisingProofCount;
    }

    // ---------------------------------------------------------------
    // Surge specific external functions
    // ---------------------------------------------------------------

    /// @notice Retrieves the hash of Surge transition records for a specific proposal and parent
    /// transition.
    /// @param _proposalId The ID of the proposal containing the transition
    /// @param _parentTransitionHash The hash of the parent transition in the proof chain
    /// @return The hash of the Surge transition records array
    function getSurgeTransitionRecordsHash(
        uint48 _proposalId,
        bytes32 _parentTransitionHash
    )
        external
        view
        returns (bytes32)
    {
        return _getSurgeTransitionRecordsHash(_proposalId, _parentTransitionHash);
    }

    // ---------------------------------------------------------------
    // Internal Overrides
    // ---------------------------------------------------------------

    /// @dev An override of `Inbox._finalize(..)` that uses the pre-image of Surge transition
    /// records stored in `_input.extra` to find the finalising records.
    function _finalize(ProposeInput memory _input)
        internal
        override
        returns (CoreState memory coreState_, LibBonds.BondInstruction[] memory bondInstructions_)
    {
        unchecked {
            CoreState memory coreState = _input.coreState;
            uint48 proposalId = coreState.lastFinalizedProposalId + 1;
            uint256 lastFinalizedRecordIdx;
            uint256 finalizedCount;
            uint256 transitionCount = _input.transitionRecords.length;
            uint256 currentTimestamp = block.timestamp;
            uint256 totalBondInstructionCount;

            // Surge: Pull the pre-images of Surge transition records from the extra data
            // in the input.
            SurgeTransitionRecord[][] memory surgeTransitionRecords =
                _decodeProposeExtra(_input.extra);

            for (uint256 i; i < _maxFinalizationCount; ++i) {
                // Check if there are more proposals to finalize
                if (proposalId >= coreState.nextProposalId) break;

                // Try to finalize the current proposal
                // Surge: Use a new handler (_processFinalizingTransitionRecord) to find the finalising
                // record details.
                (bytes26 recordHash, uint48 finalizationDeadline) = _processFinalizingTransitionRecord(
                    proposalId, coreState.lastFinalizedTransitionHash, surgeTransitionRecords[i]
                );

                if (i >= transitionCount) {
                    if (recordHash == 0) break;

                    if (currentTimestamp >= finalizationDeadline) {
                        revert TransitionRecordNotProvided();
                    }

                    break;
                }

                if (recordHash == 0) break;

                TransitionRecord memory transitionRecord = _input.transitionRecords[i];

                require(
                    _hashTransitionRecord(transitionRecord) == recordHash,
                    TransitionRecordHashMismatchWithStorage()
                );

                coreState.lastFinalizedProposalId = proposalId;
                coreState.lastFinalizedTransitionHash = transitionRecord.transitionHash;

                uint256 bondInstructionLen = transitionRecord.bondInstructions.length;
                for (uint256 j; j < bondInstructionLen; ++j) {
                    coreState.bondInstructionsHash = LibBonds.aggregateBondInstruction(
                        coreState.bondInstructionsHash, transitionRecord.bondInstructions[j]
                    );
                }

                totalBondInstructionCount += bondInstructionLen;

                require(transitionRecord.span > 0, InvalidSpan());

                uint48 nextProposalId = proposalId + transitionRecord.span;
                require(nextProposalId <= coreState.nextProposalId, SpanOutOfBounds());

                proposalId = nextProposalId;

                // Update state for successful finalization
                lastFinalizedRecordIdx = i;
                ++finalizedCount;
            }

            // Update checkpoint if any proposals were finalized and minimum delay has passed
            if (finalizedCount > 0) {
                _syncCheckpointIfNeeded(
                    _input.checkpoint,
                    _input.transitionRecords[lastFinalizedRecordIdx].checkpointHash,
                    coreState
                );
            }

            if (totalBondInstructionCount > 0) {
                bondInstructions_ = new LibBonds.BondInstruction[](totalBondInstructionCount);
                uint256 bondInstructionIndex;

                for (uint256 i; i < finalizedCount; ++i) {
                    LibBonds.BondInstruction[] memory instructions =
                    _input.transitionRecords[i].bondInstructions;
                    uint256 instructionsLen = instructions.length;

                    for (uint256 j; j < instructionsLen; ++j) {
                        bondInstructions_[bondInstructionIndex++] = instructions[j];
                    }
                }
            }

            return (coreState, bondInstructions_);
        }
    }

    /// @dev A modification of `InboxOptimized1._buildAndSaveTransitionRecords`
    function _buildAndSaveTransitionRecords(ProveInput memory _input) internal override {
        if (_input.proposals.length == 1) {
            _buildAndSaveSurgeTransitionRecord(_input);
        } else {
            _buildAndSaveAggregatedSurgeTransitionRecord(_input);
        }
    }

    /// @dev Validates that the proof bitmap sent in the proof input matches the one extract from
    /// the proof bytes.
    function _handleProofVerification(
        ProveInput memory _input,
        bytes calldata _proof
    )
        internal
        view
        override
    {
        uint256 proposalAge;
        if (_input.proposals.length == 1) {
            unchecked {
                proposalAge = block.timestamp - _input.proposals[0].timestamp;
            }
        }

        bytes32 aggregatedProvingHash =
            _hashTransitionsWithMetadata(_input.transitions, _input.metadata);

        (LibProofBitmap.ProofBitmap proofBitmapInput,) = _decodeProveExtra(_input.extra);

        LibProofBitmap.ProofBitmap proofBitmapComputed =
            SurgeVerifier(_proofVerifier).verifyProof(proposalAge, aggregatedProvingHash, _proof);
        require(proofBitmapComputed.equals(proofBitmapInput), Surge_InvalidProofBitmapInput());
    }

    /// @dev Any external function utilising this internal helper should fail since we are not longer
    /// using Taiko's record store
    function _getTransitionRecordHashAndDeadline(
        uint48,
        bytes32
    )
        internal
        pure
        override
        returns (bytes26, uint48)
    {
        revert Surge_NotSupported();
    }

    // ---------------------------------------------------------------
    // Surge specific internal functions
    // ---------------------------------------------------------------

    /// @dev Loops through the Surge transition records to find the finalising record, and marks
    /// the verifiers used in the conflicting records as upgradeable.
    /// @dev In the ideal case, there will only be one record with sufficient proofs.
    function _processFinalizingTransitionRecord(
        uint48 _proposalId,
        bytes32 _parentTransitionHash,
        SurgeTransitionRecord[] memory _surgeTransitionRecords
    )
        internal
        returns (bytes26 recordHash_, uint48 finalizationDeadline_)
    {
        _validateSurgeTransitionRecords(
            _surgeTransitionRecords,
            _getSurgeTransitionRecordsHash(_proposalId, _parentTransitionHash)
        );

        LibProofBitmap.ProofBitmap conflictingProofBitmap;

        for (uint256 i; i < _surgeTransitionRecords.length; ++i) {
            if (_surgeTransitionRecords[i].proofBitmap.numProofs() >= _finalisingProofCount) {
                // The transition has sufficient number of proofs for finalisation

                recordHash_ = _surgeTransitionRecords[i].hashAndDeadline.recordHash;
                finalizationDeadline_ =
                _surgeTransitionRecords[i].hashAndDeadline.finalizationDeadline;
            } else if (_surgeTransitionRecords.length == 1) {
                // Optimisitic path when no conflicting proofs

                // If the optimistic delay is met
                if (
                    block.timestamp
                        > _surgeTransitionRecords[i].hashAndDeadline.finalizationDeadline
                            + _optimisticFallbackDelay
                ) {
                    recordHash_ = _surgeTransitionRecords[i].hashAndDeadline.recordHash;
                    finalizationDeadline_ =
                    _surgeTransitionRecords[i].hashAndDeadline.finalizationDeadline;
                }
            } else {
                // These are conflicting proofs

                conflictingProofBitmap =
                    conflictingProofBitmap.merge(_surgeTransitionRecords[i].proofBitmap);
            }
        }

        // Allow instant upgrades of all conflicting verifiers
        if (recordHash_ != bytes26(0) && conflictingProofBitmap.toUint8() > 0) {
            SurgeVerifier(_proofVerifier).markVerifiersUpgradeable(conflictingProofBitmap, true);
            emit ConflictingProofsDetected(conflictingProofBitmap);
        }
    }

    /// @dev Loads the Surge transition record from the storage
    function _getSurgeTransitionRecordsHash(
        uint48 _proposalId,
        bytes32 _parentTransitionHash
    )
        internal
        view
        returns (bytes32)
    {
        uint256 bufferSlot = _proposalId % _ringBufferSize;
        ReusableSurgeTransitionRecord storage surgeRecord =
            _surgeTransitionRecordsHashesBuffer[bufferSlot];

        if (
            surgeRecord.proposalId == _proposalId
                && surgeRecord.partialParentTransitionHash == bytes26(_parentTransitionHash)
        ) {
            return surgeRecord.surgeTransitionRecordsHash;
        } else {
            bytes32 compositeKey = _composeTransitionKey(_proposalId, _parentTransitionHash);
            return _surgeTransitionRecordsHashes[compositeKey];
        }
    }

    /// @dev Synonymous with `InboxOptimized1._processSingleTransitionAtIndex`.
    function _buildAndSaveSurgeTransitionRecord(ProveInput memory _input) internal {
        _validateTransition(_input.proposals[0], _input.transitions[0]);

        TransitionRecord memory transitionRecord = _buildTransitionRecord(
            _input.proposals[0], _input.transitions[0], _input.metadata[0]
        );

        (
            LibProofBitmap.ProofBitmap proofBitmap,
            SurgeTransitionRecord[][] memory surgeTransitionRecords
        ) = _decodeProveExtra(_input.extra);

        _setSurgeTransitionRecord(
            _input.proposals[0].id,
            _input.transitions[0],
            _input.metadata[0],
            transitionRecord,
            proofBitmap,
            surgeTransitionRecords[0]
        );
    }

    /// @dev Synonymous with `Inbox._buildAndSaveAggregatedTransitionRecords`.
    function _buildAndSaveAggregatedSurgeTransitionRecord(ProveInput memory _input) internal {
        unchecked {
            for (uint256 i; i < _input.proposals.length; ++i) {
                _validateTransition(_input.proposals[i], _input.transitions[i]);
            }

            TransitionRecord memory currentRecord = _buildTransitionRecord(
                _input.proposals[0], _input.transitions[0], _input.metadata[0]
            );

            uint48 currentGroupStartId = _input.proposals[0].id;
            uint256 firstIndex;

            (
                LibProofBitmap.ProofBitmap proofBitmap,
                SurgeTransitionRecord[][] memory surgeTransitionRecords
            ) = _decodeProveExtra(_input.extra);
            uint256 surgeTransitionRecordsIndex;

            for (uint256 i = 1; i < _input.proposals.length; ++i) {
                if (_input.proposals[i].id == currentGroupStartId + currentRecord.span) {
                    TransitionRecord memory nextRecord = _buildTransitionRecord(
                        _input.proposals[i], _input.transitions[i], _input.metadata[i]
                    );
                    if (nextRecord.bondInstructions.length == 0) {
                        // Keep current instructions unchanged
                    } else if (currentRecord.bondInstructions.length == 0) {
                        currentRecord.bondInstructions = nextRecord.bondInstructions;
                    } else {
                        currentRecord.bondInstructions = LibBondInstruction.mergeBondInstructions(
                            currentRecord.bondInstructions, nextRecord.bondInstructions
                        );
                    }
                    currentRecord.transitionHash = nextRecord.transitionHash;
                    currentRecord.checkpointHash = nextRecord.checkpointHash;
                    currentRecord.span++;
                } else {
                    // Save current group and start new one
                    _setSurgeTransitionRecord(
                        currentGroupStartId,
                        _input.transitions[firstIndex],
                        _input.metadata[firstIndex],
                        currentRecord,
                        proofBitmap,
                        surgeTransitionRecords[surgeTransitionRecordsIndex]
                    );
                    surgeTransitionRecordsIndex++;

                    // Reset for new group
                    currentGroupStartId = _input.proposals[i].id;
                    firstIndex = i;
                    currentRecord = _buildTransitionRecord(
                        _input.proposals[i], _input.transitions[i], _input.metadata[i]
                    );
                }
            }

            // Save the final aggregated record
            _setSurgeTransitionRecord(
                currentGroupStartId,
                _input.transitions[firstIndex],
                _input.metadata[firstIndex],
                currentRecord,
                proofBitmap,
                surgeTransitionRecords[surgeTransitionRecordsIndex]
            );
        }
    }

    /// @dev Synonymous with `Inbox._setTransitionRecordHashAndDeadline`.
    function _setSurgeTransitionRecord(
        uint48 _proposalId,
        Transition memory _transition,
        TransitionMetadata memory _metadata,
        TransitionRecord memory _transitionRecord,
        LibProofBitmap.ProofBitmap _proofBitmap,
        SurgeTransitionRecord[] memory _surgeTransitionRecords
    )
        internal
    {
        (, TransitionRecordHashAndDeadline memory hashAndDeadline) =
            _computeTransitionRecordHashAndDeadline(_transitionRecord);

        _storeSurgeTransitionRecord(
            _proposalId,
            _transition.parentTransitionHash,
            hashAndDeadline,
            _proofBitmap,
            _surgeTransitionRecords
        );

        // TODO: Add the proof bitmap if the indexer needs it
        ProvedEventPayload memory payload = ProvedEventPayload({
            proposalId: _proposalId,
            transition: _transition,
            transitionRecord: _transitionRecord,
            metadata: _metadata
        });
        emit Proved(_encodeProvedEventData(payload));
    }

    /// @dev Synonymous with `InboxOptimized1._storeTransitionRecord`.
    function _storeSurgeTransitionRecord(
        uint48 _proposalId,
        bytes32 _parentTransitionHash,
        TransitionRecordHashAndDeadline memory _hashAndDeadline,
        LibProofBitmap.ProofBitmap _proofBitmap,
        SurgeTransitionRecord[] memory _surgeTransitionRecords
    )
        internal
    {
        uint256 bufferSlot = _proposalId % _ringBufferSize;
        ReusableSurgeTransitionRecord storage currentSurgeRecord =
            _surgeTransitionRecordsHashesBuffer[bufferSlot];

        bytes26 partialParentHash = bytes26(_parentTransitionHash);

        SurgeTransitionRecord memory surgeTransitionRecord =
            SurgeTransitionRecord({ proofBitmap: _proofBitmap, hashAndDeadline: _hashAndDeadline });

        if (currentSurgeRecord.proposalId != _proposalId) {
            // No existing transition for the proposal

            currentSurgeRecord.proposalId = _proposalId;
            currentSurgeRecord.partialParentTransitionHash = partialParentHash;
            currentSurgeRecord.surgeTransitionRecordsHash =
                _hashSurgeTransitionRecord(surgeTransitionRecord);
        } else if (currentSurgeRecord.partialParentTransitionHash == partialParentHash) {
            // Potentially conflicting or duplicate transition

            _validateSurgeTransitionRecords(
                _surgeTransitionRecords, currentSurgeRecord.surgeTransitionRecordsHash
            );
            currentSurgeRecord.surgeTransitionRecordsHash = _getUpdatedSurgeTransitionRecordsHash(
                _surgeTransitionRecords, _hashAndDeadline, _proofBitmap
            );
        } else {
            // Parallel transitions go into the composite key indexed mapping

            bytes32 compositeKey = _composeTransitionKey(_proposalId, _parentTransitionHash);
            _surgeTransitionRecordsHashes[compositeKey] = _getUpdatedSurgeTransitionRecordsHash(
                _surgeTransitionRecords, _hashAndDeadline, _proofBitmap
            );
        }
    }

    /// @dev If the proof belongs to a matching transition, the records proof bitmap is
    /// updated to accomodate the potentially new proof variant. Otherwise, the conflicting
    /// transition and proof is appended as a new transition record.
    function _getUpdatedSurgeTransitionRecordsHash(
        SurgeTransitionRecord[] memory _surgeTransitionRecords,
        TransitionRecordHashAndDeadline memory _hashAndDeadline,
        LibProofBitmap.ProofBitmap _proofBitmap
    )
        internal
        pure
        returns (bytes32)
    {
        // TODO: This can be optimized to return the hash right at the start if
        // _surgeTransitionRecords is empty

        for (uint256 i; i < _surgeTransitionRecords.length; ++i) {
            // If a matching transition record is found, merge the proof bitmaps
            if (
                _surgeTransitionRecords[i].hashAndDeadline.recordHash == _hashAndDeadline.recordHash
            ) {
                _surgeTransitionRecords[i].proofBitmap =
                    _surgeTransitionRecords[i].proofBitmap.merge(_proofBitmap);
                return _hashSurgeTransitionRecords(_surgeTransitionRecords);
            }
        }

        SurgeTransitionRecord memory newSurgeTransitionRecord =
            SurgeTransitionRecord({ proofBitmap: _proofBitmap, hashAndDeadline: _hashAndDeadline });

        // We have a conflicting proof that needs to be appended to the list of records
        // TODO: Move to optimised hashing
        SurgeTransitionRecord[] memory updatedSurgeTransitionRecords =
            new SurgeTransitionRecord[](_surgeTransitionRecords.length + 1);
        for (uint256 i = 0; i < _surgeTransitionRecords.length; ++i) {
            updatedSurgeTransitionRecords[i] = _surgeTransitionRecords[i];
        }
        updatedSurgeTransitionRecords[_surgeTransitionRecords.length] = newSurgeTransitionRecord;
        return _hashSurgeTransitionRecords(updatedSurgeTransitionRecords);
    }

    function _validateSurgeTransitionRecords(
        SurgeTransitionRecord[] memory _record,
        bytes32 _recordsHash
    )
        internal
        pure
    {
        require(
            _hashSurgeTransitionRecords(_record) == _recordsHash, Surge_InvalidTransitionRecords()
        );
    }

    // ---------------------------------------------------------------
    // Transition record hashing
    // ---------------------------------------------------------------

    /// @notice Computes the keccak256 hash of a single SurgeTransitionRecord.
    /// @param _record The SurgeTransitionRecord struct to hash.
    /// @return The keccak256 hash of the encoded record.
    function _hashSurgeTransitionRecord(SurgeTransitionRecord memory _record)
        internal
        pure
        returns (bytes32)
    {
        return keccak256(abi.encode(_record));
    }

    /// @notice Computes the keccak256 hash of an array of SurgeTransitionRecords.
    /// @param _records The array of SurgeTransitionRecord structs to hash.
    /// @return The keccak256 hash of the encoded records.
    function _hashSurgeTransitionRecords(SurgeTransitionRecord[] memory _records)
        internal
        pure
        returns (bytes32)
    {
        return keccak256(abi.encode(_records));
    }

    // ---------------------------------------------------------------
    // Extra Data Decoding Functions
    // ---------------------------------------------------------------

    /// @dev Decodes the extra data for FinalityGadgetInbox prove operations
    /// @param _data The encoded extra data bytes
    /// @return proofBitmap_ The proof bitmap indicating which proof variants were used
    /// @return The array of Surge transition record arrays
    function _decodeProveExtra(bytes memory _data)
        internal
        pure
        virtual
        returns (LibProofBitmap.ProofBitmap, SurgeTransitionRecord[][] memory)
    {
        return abi.decode(_data, (LibProofBitmap.ProofBitmap, SurgeTransitionRecord[][]));
    }

    /// @dev Decodes the extra data for FinalityGadgetInbox propose operations
    /// @param _data The encoded extra data bytes
    /// @return The array of Surge transition record arrays
    function _decodeProposeExtra(bytes memory _data)
        internal
        pure
        virtual
        returns (SurgeTransitionRecord[][] memory)
    {
        return abi.decode(_data, (SurgeTransitionRecord[][]));
    }

    // ---------------------------------------------------------------
    // Custom errors
    // ---------------------------------------------------------------

    error Surge_InvalidProofBitmapInput();
    error Surge_InvalidTransitionRecords();
    error Surge_NotSupported();
}
