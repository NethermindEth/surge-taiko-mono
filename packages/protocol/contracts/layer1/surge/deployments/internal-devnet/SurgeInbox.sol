// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { FinalityGadgetInbox } from "../../features/FinalityGadgetInbox.sol";
import { RollbackInbox } from "../../features/RollbackInbox.sol";
import { IInbox } from "src/layer1/core/iface/IInbox.sol";
import { Inbox } from "src/layer1/core/impl/Inbox.sol";

/// @title SurgeInbox
/// @notice Surge inbox implementation for internal-devnet deployment
/// @custom:security-contact security@nethermind.io
contract SurgeInbox is FinalityGadgetInbox, RollbackInbox {
    /// @param _config The inbox configuration
    /// @param _maxFinalizationDelay The maximum grace period after which the chain can be
    /// rollbacked to the last finalized proposal
    constructor(
        IInbox.Config memory _config,
        uint48 _maxFinalizationDelay
    )
        Inbox(_config)
        RollbackInbox(_maxFinalizationDelay)
    { }

    /// @dev Resolves diamond inheritance conflict for _handleOnPropose
    function _handleOnPropose() internal override(Inbox, RollbackInbox) {
        super._handleOnPropose();
    }

    /// @dev Resolves diamond inheritance conflict for _handleOnProve
    function _handleOnProve() internal override(Inbox, RollbackInbox) {
        super._handleOnProve();
    }

    /// @dev Resolves diamond inheritance conflict for _handleProofVerification
    function _handleProofVerification(
        Commitment memory _commitment,
        bytes calldata _proof
    )
        internal
        view
        override(Inbox, FinalityGadgetInbox)
    {
        super._handleProofVerification(_commitment, _proof);
    }
}
