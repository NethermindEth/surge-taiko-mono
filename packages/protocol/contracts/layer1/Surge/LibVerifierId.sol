// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title LibVerifierId
/// @dev The bindings from this library shall be used by the client to compute the verifier id.
/// @custom:security-contact security@nethermind.io
library LibVerifierId {
    enum Verifier {
        RISC0_RETH,
        SP1_RETH
    }

    function getVerifierId(Verifier _verifier) public pure returns (uint8) {
        return uint8(_verifier);
    }
}
