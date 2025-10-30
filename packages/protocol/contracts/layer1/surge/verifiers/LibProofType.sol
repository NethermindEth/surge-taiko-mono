// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title LibProofType
/// @dev This library offers a custom type to represent proofs.
/// @custom:security-contact security@nethermind.io
library LibProofType {
    // This represents a bitmap of proof types, allowing for up to 16 distinct proof types.
    type ProofType is uint16;

    uint8 internal constant NUM_PROOF_TYPES = 7;

    uint16 internal constant ZK_MASK = 0x30; // 0b110000
    uint16 internal constant TEE_MASK = 0x0F; // 0b001111

    // Invidual proof types
    // --------------------

    /// @dev Empty proof type (0b00000)
    function empty() internal pure returns (ProofType) {
        return ProofType.wrap(0x00);
    }

    /// @dev SGX Reth proof type (0b000001)
    function sgxReth() internal pure returns (ProofType) {
        return ProofType.wrap(0x01);
    }

    /// @dev SGX Geth proof type (0b000010)
    function sgxGeth() internal pure returns (ProofType) {
        return ProofType.wrap(0x02);
    }

    /// @dev TDX Nethermind proof type (0b000100)
    function tdxNethermind() internal pure returns (ProofType) {
        return ProofType.wrap(0x04);
    }

    /// @dev Azure TDX Nethermind proof type (0b001000)
    function azureTdxNethermind() internal pure returns (ProofType) {
        return ProofType.wrap(0x08);
    }

    /// @dev RISC-0 Reth proof type (0b010000)
    function risc0Reth() internal pure returns (ProofType) {
        return ProofType.wrap(0x10);
    }

    /// @dev SP1 Reth proof type (0b100000)
    function sp1Reth() internal pure returns (ProofType) {
        return ProofType.wrap(0x20);
    }

    // ZK / TEE type detectors
    // -----------------------

    function isZkProof(ProofType _proofType) internal pure returns (bool) {
        uint16 pt = ProofType.unwrap(_proofType);
        return (pt & ZK_MASK) != 0 && (pt & TEE_MASK) == 0;
    }

    function isTeeProof(ProofType _proofType) internal pure returns (bool) {
        uint16 pt = ProofType.unwrap(_proofType);
        return (pt & ZK_MASK) == 0 && (pt & TEE_MASK) != 0;
    }

    function isZkTeeProof(ProofType _proofType) internal pure returns (bool) {
        uint16 pt = ProofType.unwrap(_proofType);
        return (pt & ZK_MASK) != 0 && (pt & TEE_MASK) != 0;
    }

    function isFinalizingProof(ProofType _proofType) internal pure returns (bool) {
        return isZkTeeProof(_proofType);
    }

    // Misc helpers
    // ------------

    function equals(ProofType _proofType1, ProofType _proofType2) internal pure returns (bool) {
        uint16 pt1 = ProofType.unwrap(_proofType1);
        uint16 pt2 = ProofType.unwrap(_proofType2);
        return pt1 == pt2;
    }

    function combine(
        ProofType _proofType1,
        ProofType _proofType2
    )
        internal
        pure
        returns (ProofType)
    {
        uint16 pt1 = ProofType.unwrap(_proofType1);
        uint16 pt2 = ProofType.unwrap(_proofType2);
        return ProofType.wrap(pt1 | pt2);
    }

    function includes(
        ProofType _proofType1,
        ProofType _proofType2
    )
        internal
        pure
        returns (bool)
    {
        uint16 pt1 = ProofType.unwrap(_proofType1);
        uint16 pt2 = ProofType.unwrap(_proofType2);
        return (pt1 & pt2) != 0;
    }
}
