// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { Inbox } from "../../core/impl/Inbox.sol";

abstract contract RollbackInbox is Inbox {
    /// @notice Emitted when a rollback operation is performed over a range of proposals.
    /// @param firstProposalId The ID of the first proposal reverted in the rollback.
    /// @param lastProposalId The ID of the last proposal reverted in the rollback.
    event Rollbacked(uint256 firstProposalId, uint256 lastProposalId);

    /// @notice Emitted when limp mode is enabled or disabled.
    /// @param enabled True if limp mode was enabled, false if disabled.
    event LimpModeSet(bool enabled);

    /// @dev Maximum grace period after which the chain can be rollbacked to the last finalized proposal.
    uint48 public immutable maxFinalizationDelayBeforeRollback;

    /// @dev When set to `true`, proposals must be accompanied with the associated proof
    bool public inLimpMode;

    uint256[49] private __gap;

    constructor(uint48 _maxFinalizationDelayBeforeRollback) {
        maxFinalizationDelayBeforeRollback = _maxFinalizationDelayBeforeRollback;
    }

    /// @notice Rolls back unfinalized proposals if the finalization window has been exceeded.
    /// @dev This allows recovery when the chain has stalled without finalization for too long,
    /// for instance when a prover killer block has been published.
    function rollback() external {
        // Check if the last finalization exceeds the maxFinalizationDelayBeforeRollback
        require(
            block.timestamp > _coreState.lastFinalizedTimestamp + maxFinalizationDelayBeforeRollback,
            Surge_RollbackNotAllowed()
        );

        uint48 lastFinalizedProposalId = _coreState.lastFinalizedProposalId;
        uint48 nextProposalId = _coreState.nextProposalId;

        // Only rollback if there are unfinalized proposals
        require(nextProposalId > lastFinalizedProposalId + 1, Surge_NoProposalsToRollback());

        // Rollback to the last finalized proposal and enable limp mode
        // When in limp mode, proposals must be accompanied with the associated proof
        _coreState.nextProposalId = lastFinalizedProposalId + 1;
        inLimpMode = true;

        emit Rollbacked(lastFinalizedProposalId + 1, nextProposalId - 1);
    }

    /// @notice Allows the owner to enable or disable limp mode.
    /// @param _val True to enable limp mode, false to disable.
    function setLimpMode(bool _val) external onlyOwner {
        inLimpMode = _val;
        emit LimpModeSet(_val);
    }

    /// @notice Proposes a new batch and immediately proves it, ensuring the head is always finalized.
    /// @param _lookahead Additional data used for lookahead operations.
    /// @param _proposeData The encoded proposal input data.
    /// @param _proveData The encoded ProveInput struct.
    /// @param _proof Validity proof for the batch of proposals.
    function proposeAndProve(
        bytes calldata _lookahead,
        bytes calldata _proposeData,
        bytes calldata _proveData,
        bytes calldata _proof
    )
        external
    {
        _propose(_lookahead, _proposeData);
        _prove(_proveData, _proof);

        // Verify that the head of the chain is finalized
        require(
            _coreState.lastFinalizedProposalId == _coreState.nextProposalId - 1,
            Surge_HeadMustBeFinalizedInLimpMode()
        );
    }

    // ---------------------------------------------------------------
    // Overrides
    // ---------------------------------------------------------------

    /// @dev Disable calling the `propose(..)` function directly when in limp mode.
    /// Proposals must be accompanied by the associated proof.
    function _handleOnPropose() internal view override {
        require(!inLimpMode, Surge_CannotProposeDirectlyInLimpMode());
    }

    /// @dev Disable calling the `prove(..)` function directly when in limp mode.
    /// Proposals and proofs must be jointly submitted.
    function _handleOnProve() internal view override {
        require(!inLimpMode, Surge_CannotProveDirectlyInLimpMode());
    }

    // ---------------------------------------------------------------
    // Custom errors
    // ---------------------------------------------------------------

    error Surge_CannotProposeDirectlyInLimpMode();
    error Surge_CannotProveDirectlyInLimpMode();
    error Surge_HeadMustBeFinalizedInLimpMode();
    error Surge_NoProposalsToRollback();
    error Surge_RollbackNotAllowed();
}
