// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { IInbox } from "../../core/iface/IInbox.sol";
import { Inbox } from "../../core/impl/Inbox.sol";

/// @title RollbackInbox
/// @dev This is intended to be for demonstration only.
/// @notice Inbox variant that implements chain rollback and limp mode feature.
/// @custom:security-contact security@nethermind.io
contract RollbackInbox is Inbox {
    struct RollbackMetadata {
        /// @notice When true, you may only propose along with a proof.
        bool inLimpMode;
        /// @notice When true, the next proposal must begin at `rollbackToProposalId`.
        bool rollbackOnNextProposal;
        /// @notice The proposal id to use for the next proposal.
        uint48 proposalIdOnNextProposal;
    }

    /// @notice Emitted when a rollback is performed on the chain.
    /// @param firstProposalId The ID of the first proposal to be rollbacked.
    /// @param lastProposalId The ID of the last proposal to be rollbacked.
    event Rollback(uint48 firstProposalId, uint48 lastProposalId);

    /// @dev The maximum period of non-finalization after which a chain rollback can be
    /// permissionlessly triggered.
    uint48 internal immutable _maxFinalizationWindow;

    /// @dev 1 slot
    RollbackMetadata internal _rollbackMetadata;

    uint256[49] private __gap;

    /// @notice Initializes the RollbackInbox contract
    /// @param _config Default configuration struct containing all constructor parameters
    /// @param __maxFinalizationWindow Maximum finalization window is seconds
    constructor(IInbox.Config memory _config, uint48 __maxFinalizationWindow) Inbox(_config) {
        _maxFinalizationWindow = __maxFinalizationWindow;
    }

    // ---------------------------------------------------------------
    // Feature specific external functions
    // ---------------------------------------------------------------

    /// @notice Triggers a rollback of the chain state to the last finalized proposal if the
    /// finalization window has passed.
    /// @dev Sets the contract into limp mode and flags that the next proposal must begin at a rollbacked proposalId.
    /// @param _coreState The current core state at the chain head, to check against parent proposal.
    /// @param _headProposals Array containing the current chain head proposal(s) for verification.
    /// @param _firstUnfinalizedProposal The first proposal after the last finalized one
    function rollback(
        CoreState calldata _coreState,
        Proposal[] calldata _headProposals,
        Proposal calldata _firstUnfinalizedProposal
    )
        external
    {
        // Use the default head proposal verification
        super._handleChainHeadVerification(_headProposals);

        // Do not rollback twice for the same state
        require(!_rollbackMetadata.rollbackOnNextProposal, AlreadyInRollbackedState());

        // Essential input validations
        require(_hashCoreState(_coreState) == _headProposals[0].coreStateHash, InvalidCoreState());
        require(
            _coreState.lastFinalizedProposalId == _firstUnfinalizedProposal.id - 1,
            InvalidLastFinalizedProposal()
        );
        _checkProposalHash(_firstUnfinalizedProposal);

        // Max finalization window must have been crossed
        require(
            block.timestamp > (_firstUnfinalizedProposal.timestamp + _maxFinalizationWindow),
            MaxFinalizationWindowNotCrossed()
        );

        // Enforce a chain rollback on next proposal
        _rollbackMetadata = RollbackMetadata(true, true, _firstUnfinalizedProposal.id);

        emit Rollback(_firstUnfinalizedProposal.id, _headProposals[0].id);
    }

    /// @notice Sets the contract into limp mode, disabling proposal and proving operations.
    /// @dev Only callable externally, e.g., by a governance or admin address if such access is desired.
    function setLimpMode(bool _inLimpMode) external {
        // Insert appropriate access control in production (e.g., onlyOwner or onlyGovernance)
        _rollbackMetadata.inLimpMode = _inLimpMode;
    }

    /// @notice Submits a proposal and immediately proves its validity in a single transaction.
    /// @dev Combines propose and prove logic for atomic submission during limp mode.
    /// @param _lookahead Data for the preconf lookahead contract
    /// @param _proposeData ABI-encoded ProposeInput struct for the proposal.
    /// @param _proveData ABI-encoded ProveInput struct containing exactly one proposal,
    /// which must be the proposed one.
    /// @param _proof Cryptographic proof for the proposal
    function proposeAndProve(
        bytes calldata _lookahead,
        bytes calldata _proposeData,
        bytes calldata _proveData,
        bytes calldata _proof
    )
        external
    {
        ProposeInput memory _proposeInput = abi.decode(_proposeData, (ProposeInput));
        ProveInput memory _proveInput = abi.decode(_proveData, (ProveInput));

        // Validate that the proposal is accompoanied with its proof
        require(
            _proveInput.proposals.length == 1
                && _proveInput.proposals[0].id == _proposeInput.parentProposals[0].id + 1,
            InvalidDataForLimpMode()
        );

        _propose(_lookahead, _proposeInput);
        _prove(_proveInput, _proof);
    }

    // ---------------------------------------------------------------
    // Surge Handlers Overrides
    // ---------------------------------------------------------------

    /// @dev If we expect a rollback on the next proposal, we simply verify the head against the
    /// data stored by the `rollback` function.
    function _handleChainHeadVerification(Proposal[] memory _parentProposals) internal override {
        RollbackMetadata memory meta = _rollbackMetadata;

        if (meta.rollbackOnNextProposal) {
            // If the chain is in the rollback mode, we assume that the `rollback` function recorded
            // correct data, and simply check that the parent proposal is the proposal that we have
            // rollbacked to.
            _checkProposalHash(_parentProposals[0]);
            require(
                _parentProposals[0].id == meta.proposalIdOnNextProposal - 1, InvalidParentProposal()
            );
            _rollbackMetadata.rollbackOnNextProposal = false;
        } else {
            super._handleChainHeadVerification(_parentProposals);
        }
    }

    /// @dev Block direct proposing in limp mode
    function _handleOnPropose(ProposeInput memory _input) internal view override {
        require(!_rollbackMetadata.inLimpMode, InLimpMode());
    }

    /// @dev Block direct proving in limp mode
    function _handleOnProve(ProveInput memory _input) internal view override {
        require(!_rollbackMetadata.inLimpMode, InLimpMode());
    }

    // ---------------------------------------------------------------
    // Custom errors
    // ---------------------------------------------------------------

    error AlreadyInRollbackedState();
    error CannotProposeWithoutProofInLimpMode();
    error InLimpMode();
    error InvalidCoreState();
    error InvalidDataForLimpMode();
    error InvalidLastFinalizedProposal();
    error InvalidParentProposal();
    error MaxFinalizationWindowNotCrossed();
}

