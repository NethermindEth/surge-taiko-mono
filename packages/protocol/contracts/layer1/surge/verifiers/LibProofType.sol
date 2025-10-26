// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title LibProofType
/// @dev This library offers a custom type to represent proofs.
/// @custom:security-contact security@nethermind.io
library LibProofType {
    // This represents a bitmap of proof types, allowing for up to 16 distinct proof types.
    // Bitmap layout: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, SP1_RETH, RISC0_RETH, TDX_RETH, SGX_RETH]
    type ProofType is uint16;

    uint8 internal constant NUM_PROOF_TYPES = 4;

    uint16 internal constant ZK_MASK = 0x18; // 0b11000
    uint16 internal constant TEE_MASK = 0x07; // 0b00111

    // Invidual proof types
    // --------------------

    /// @dev Empty proof type (0b0000)
    function empty() internal pure returns (ProofType) {
        return ProofType.wrap(0x00);
    }

    /// @dev SGX Reth proof type (0b00001)
    function sgxReth() internal pure returns (ProofType) {
        return ProofType.wrap(0x01);
    }

    /// @dev TDX Nethermind proof type (0b00010)
    function tdxNethermind() internal pure returns (ProofType) {
        return ProofType.wrap(0x02);
    }

    /// @dev Azure TDX Nethermind proof type (0b00100)
    function azureTdxNethermind() internal pure returns (ProofType) {
        return ProofType.wrap(0x04);
    }

    /// @dev RISC-0 Reth proof type (0b01000)
    function risc0Reth() internal pure returns (ProofType) {
        return ProofType.wrap(0x08);
    }

    /// @dev SP1 Reth proof type (0b10000)
    function sp1Reth() internal pure returns (ProofType) {
        return ProofType.wrap(0x10);
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
}
