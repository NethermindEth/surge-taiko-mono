// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title LibProofBitmap
/// @dev Bindings from this library can be used to build the proof bitmap on client side
/// @custom:security-contact security@nethermind.io
library LibProofBitmap {
    uint8 internal constant RISC0_RETH = 1; // 0b00000001
    uint8 internal constant SP1_RETH = 1 << 1; // 0b00000010

    type ProofBitmap is uint8;

    /// @dev Combines two proof bitmaps
    function merge(ProofBitmap _x, ProofBitmap _y) internal pure returns (ProofBitmap) {
        return ProofBitmap.wrap(ProofBitmap.unwrap(_x) | ProofBitmap.unwrap(_y));
    }

    /// @dev SWAR method to find number of set bits
    function numProofs(ProofBitmap _x) internal pure returns (uint8) {
        uint8 x = ProofBitmap.unwrap(_x);
        x = x - ((x >> 1) & 0x55);
        x = (x & 0x33) + ((x >> 2) & 0x33);
        x = (x + (x >> 4)) & 0x0F;
        return x;
    }

    /// @dev Checks if two ProofBitmap values are equal
    function equals(ProofBitmap _a, ProofBitmap _b) internal pure returns (bool) {
        return ProofBitmap.unwrap(_a) == ProofBitmap.unwrap(_b);
    }

    function toUint8(ProofBitmap _x) internal pure returns (uint8) {
        return ProofBitmap.unwrap(_x);
    }
}
