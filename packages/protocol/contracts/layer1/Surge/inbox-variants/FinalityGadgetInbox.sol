// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { IInbox } from "../../core/iface/IInbox.sol";
import { Inbox } from "../../core/impl/Inbox.sol";
import { SurgeVerifier } from "../SurgeVerifier.sol";

/// @title FinalityGadgetInbox
/// @dev This is intended to be for demonstration only.
/// @notice Inbox variant that implements finality gadget logic for surge handlers.
/// @custom:security-contact security@nethermind.io
contract FinalityGadgetInbox is Inbox {
    /// @notice Initializes the FinalityGadgetInbox contract
    /// @param _config Configuration struct containing all constructor parameters
    constructor(IInbox.Config memory _config) Inbox(_config) { }

    // ---------------------------------------------------------------
    // Surge Handlers Overrides
    // ---------------------------------------------------------------

    /// @dev Validates that the verifier id sent in the proof input matches the one extract from
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

        uint256 verifierId =
            SurgeVerifier(_proofVerifier).verifyProof(proposalAge, aggregatedProvingHash, _proof);
        require(verifierId == _input.verifierId, InvalidVerifierId());
    }

    /// @dev Handles a transition conflict by marking the detected conflict and enabling instant upgrade
    /// for the conflicting verifiers.
    function _handleTransitionConflict(
        TransitionRecordHashAndDeadline storage _entry,
        uint8 _conflictingVerifierId
    )
        internal
        override
    {
        // Existing state changes
        conflictingTransitionDetected = true;
        _entry.finalizationDeadline = type(uint40).max;

        // Allow instant upgrade of conflicting verifiers
        SurgeVerifier(_proofVerifier).markVerifierUpgradeable(_entry.verifierId, true);
        SurgeVerifier(_proofVerifier).markVerifierUpgradeable(_conflictingVerifierId, true);
    }

    // ---------------------------------------------------------------
    // Custom errors
    // ---------------------------------------------------------------

    error InvalidVerifierId();
}

