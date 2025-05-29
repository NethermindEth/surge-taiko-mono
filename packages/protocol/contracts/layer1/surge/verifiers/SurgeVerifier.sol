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

    function getProofTypeFromVerifiers(address[] memory _verifiers)
        internal
        view
        override
        returns (LibProofType.ProofType)
    {
        if (_verifiers.length == 2) {
            if (_verifiers[0] == sgxRethVerifier.addr) {
                if (_verifiers[1] == sp1RethVerifier.addr) {
                    return LibProofType.ProofType.SGX_SP1;
                } else if (_verifiers[1] == risc0RethVerifier.addr) {
                    return LibProofType.ProofType.SGX_RISC0;
                }
            } else if (_verifiers[0] == risc0RethVerifier.addr) {
                if (_verifiers[1] == sgxRethVerifier.addr) {
                    return LibProofType.ProofType.SGX_RISC0;
                }
            } else if (_verifiers[0] == sp1RethVerifier.addr) {
                if (_verifiers[1] == sgxRethVerifier.addr) {
                    return LibProofType.ProofType.SGX_SP1;
                }
            }
        } else if (_verifiers.length == 1) {
            if (_verifiers[0] == sgxRethVerifier.addr) {
                return LibProofType.ProofType.SGX;
            } else if (_verifiers[0] == risc0RethVerifier.addr) {
                return LibProofType.ProofType.RISC0;
            } else if (_verifiers[0] == sp1RethVerifier.addr) {
                return LibProofType.ProofType.SP1;
            }
        }

        return LibProofType.ProofType.INVALID;
    }
}
