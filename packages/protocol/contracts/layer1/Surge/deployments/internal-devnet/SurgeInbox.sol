// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { FinalityGadgetInbox } from "../../features/FinalityGadgetInbox.sol";
import { IInbox } from "src/layer1/core/iface/IInbox.sol";
import { Inbox } from "src/layer1/core/impl/Inbox.sol";

/// @title SurgeInbox
/// @notice Surge inbox implementation for internal-devnet deployment
/// @dev Inherits from FinalityGadgetInbox providing optimistic finality gadget functionality
/// @custom:security-contact security@nethermind.io
contract SurgeInbox is FinalityGadgetInbox {
    /// @param _config The inbox configuration
    /// @param _optimisticFallbackDelay The delay before which a single proof with no conflicts
    /// can be used for finalising a transition
    /// @param _finalisingProofCount The minimum number of distinct proofs required for a
    /// transition to be finalising
    constructor(
        IInbox.Config memory _config,
        uint256 _optimisticFallbackDelay,
        uint8 _finalisingProofCount
    )
        Inbox(_config)
        FinalityGadgetInbox(_optimisticFallbackDelay, _finalisingProofCount)
    { }
}

