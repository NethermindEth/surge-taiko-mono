// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/layer1/verifiers/compose/AnyTwoVerifier.sol";

/// @title SurgeVerifier
/// @notice Verifier for Surge protocol that requires any two of SGX, RISC0, or SP1 verifiers
/// @custom:security-contact security@nethermind.io
contract SurgeVerifier is AnyTwoVerifier {
    uint256[50] private __gap;

    constructor(
        address _taikoInbox,
        address _sgxRethVerifier,
        address _risc0RethVerifier,
        address _sp1RethVerifier
    )
        AnyTwoVerifier(_taikoInbox, _sgxRethVerifier, _risc0RethVerifier, _sp1RethVerifier)
    { }
}
