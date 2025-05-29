// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

library LibProofType {
    enum ProofType {
        // Only SGX proof is provided
        SGX,
        // Only SP1 proof is provided
        SP1,
        // Only RISC0 proof is provided
        RISC0,
        // Combined proof from SGX and SP1
        SGX_SP1,
        // Combined proof from SGX and RISC0
        SGX_RISC0,
        // Represents an invalid or unrecognized proof type
        INVALID
    }

    error LPT_COMBINATION_NOT_SUPPORTED();

    function isZkProof(ProofType _proofType) internal pure returns (bool) {
        return _proofType == ProofType.RISC0 || _proofType == ProofType.SP1;
    }

    function isTeeProof(ProofType _proofType) internal pure returns (bool) {
        return _proofType == ProofType.SGX;
    }

    function isZkTeeProof(ProofType _proofType) internal pure returns (bool) {
        return _proofType == ProofType.SGX_SP1 || _proofType == ProofType.SGX_RISC0;
    }

    function combine(
        ProofType _proofType1,
        ProofType _proofType2
    )
        internal
        pure
        returns (ProofType)
    {
        if (
            (_proofType1 == ProofType.SGX && _proofType2 == ProofType.SP1)
                || (_proofType1 == ProofType.SP1 && _proofType2 == ProofType.SGX)
        ) {
            return ProofType.SGX_SP1;
        } else if (
            (_proofType1 == ProofType.SGX && _proofType2 == ProofType.RISC0)
                || (_proofType1 == ProofType.RISC0 && _proofType2 == ProofType.SGX)
        ) {
            return ProofType.SGX_RISC0;
        } else {
            revert LPT_COMBINATION_NOT_SUPPORTED();
        }
    }
}
