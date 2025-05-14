// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./SurgeComposeVerifier.sol";

/// @title SurgeVerifier
/// @notice Verifier based on Surge's finality gadget
/// @custom:security-contact security@nethermind.io
contract SurgeVerifier is SurgeComposeVerifier {
    uint256[50] private __gap;

    constructor(
        address _taikoInbox,
        address _sgxRethVerifier,
        address _risc0RethVerifier,
        address _sp1RethVerifier
    )
        SurgeComposeVerifier(
            _taikoInbox,
            address(0),
            address(0),
            address(0),
            _sgxRethVerifier,
            _risc0RethVerifier,
            _sp1RethVerifier
        )
    { }

    function areVerifiersSufficient(address[] memory _verifiers)
        internal
        view
        override
        returns (bool, ITaikoInbox.ProofType)
    {
        if (_verifiers.length == 2) {
            if (_verifiers[0] == sgxRethVerifier) {
                return (
                    _verifiers[1] == risc0RethVerifier || _verifiers[1] == sp1RethVerifier,
                    ITaikoInbox.ProofType.ZK_TEE
                );
            } else if (_verifiers[0] == risc0RethVerifier || _verifiers[0] == sp1RethVerifier) {
                return (_verifiers[1] == sgxRethVerifier, ITaikoInbox.ProofType.ZK);
            }
        } else if (_verifiers.length == 1) {
            if (_verifiers[0] == sgxRethVerifier) {
                return (true, ITaikoInbox.ProofType.TEE);
            } else {
                return (
                    _verifiers[0] == risc0RethVerifier || _verifiers[0] == sp1RethVerifier,
                    ITaikoInbox.ProofType.ZK
                );
            }
        }

        return (false, ITaikoInbox.ProofType.INVALID);
    }
}
