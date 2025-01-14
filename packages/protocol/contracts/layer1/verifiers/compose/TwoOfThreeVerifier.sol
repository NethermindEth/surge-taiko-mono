// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/shared/common/LibStrings.sol";
import "./ComposeVerifier.sol";

/// @title TwoOfThreeVerifier
/// @dev Surge: Allows using any two of the three verifier from SGX / RISC0 / SP1
/// @custom:security-contact security@nethermind.io
contract TwoOfThreeVerifier is ComposeVerifier {
    uint256[50] private __gap;

    /// @inheritdoc ComposeVerifier
    function getSubVerifiersAndThreshold()
        public
        view
        override
        returns (address[] memory verifiers_, uint256 numSubProofs_)
    {
        verifiers_ = new address[](3);
        verifiers_[0] = resolve(LibStrings.B_TIER_ZKVM_RISC0, true);
        verifiers_[1] = resolve(LibStrings.B_TIER_ZKVM_SP1, true);
        verifiers_[2] = resolve(LibStrings.B_TIER_SGX, true);
        numSubProofs_ = 2;
    }
}
