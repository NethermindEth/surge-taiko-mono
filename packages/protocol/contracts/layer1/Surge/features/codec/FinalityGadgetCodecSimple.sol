// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { LibProofBitmap } from "../../LibProofBitmap.sol";
import { FinalityGadgetInbox } from "../FinalityGadgetInbox.sol";
import { CodecSimple } from "src/layer1/core/impl/CodecSimple.sol";

/// @title FinalityGadgetCodecSimple
/// @notice Codec contract extending CodecSimple with Surge FinalityGadget-specific encoding/decoding
/// @dev Provides encoding and decoding functions for extra data used in FinalityGadgetInbox
/// @custom:security-contact security@nethermind.io
contract FinalityGadgetCodecSimple is CodecSimple {
    // ---------------------------------------------------------------
    // Extra Data Encoding/Decoding Functions (for prove path)
    // ---------------------------------------------------------------

    /// @notice Encodes the extra data for FinalityGadgetInbox prove operations
    /// @param _proofBitmap The proof bitmap indicating which proof variants were used
    /// @param _surgeTransitionRecords The array of Surge transition record arrays
    /// @return encoded_ The encoded extra data bytes
    function encodeProveExtra(
        LibProofBitmap.ProofBitmap _proofBitmap,
        FinalityGadgetInbox.SurgeTransitionRecord[][] calldata _surgeTransitionRecords
    )
        external
        pure
        returns (bytes memory encoded_)
    {
        return abi.encode(_proofBitmap, _surgeTransitionRecords);
    }

    /// @notice Decodes the extra data for FinalityGadgetInbox prove operations
    /// @param _data The encoded extra data bytes
    /// @return proofBitmap_ The proof bitmap indicating which proof variants were used
    /// @return surgeTransitionRecords_ The array of Surge transition record arrays
    function decodeProveExtra(bytes calldata _data)
        external
        pure
        returns (
            LibProofBitmap.ProofBitmap proofBitmap_,
            FinalityGadgetInbox.SurgeTransitionRecord[][] memory surgeTransitionRecords_
        )
    {
        return abi.decode(
            _data, (LibProofBitmap.ProofBitmap, FinalityGadgetInbox.SurgeTransitionRecord[][])
        );
    }

    // ---------------------------------------------------------------
    // Extra Data Encoding/Decoding Functions (for propose path)
    // ---------------------------------------------------------------

    /// @notice Encodes the extra data for FinalityGadgetInbox propose operations
    /// @dev Used during finalization when only transition records are needed
    /// @param _surgeTransitionRecords The array of Surge transition record arrays
    /// @return encoded_ The encoded extra data bytes
    function encodeProposeExtra(FinalityGadgetInbox
                .SurgeTransitionRecord[][] calldata _surgeTransitionRecords)
        external
        pure
        returns (bytes memory encoded_)
    {
        return abi.encode(_surgeTransitionRecords);
    }

    /// @notice Decodes the extra data for FinalityGadgetInbox propose operations
    /// @dev Used during finalization when only transition records are needed
    /// @param _data The encoded extra data bytes
    /// @return surgeTransitionRecords_ The array of Surge transition record arrays
    function decodeProposeExtra(bytes calldata _data)
        external
        pure
        returns (FinalityGadgetInbox.SurgeTransitionRecord[][] memory surgeTransitionRecords_)
    {
        return abi.decode(_data, (FinalityGadgetInbox.SurgeTransitionRecord[][]));
    }

    // ---------------------------------------------------------------
    // Hashing Functions
    // ---------------------------------------------------------------

    /// @notice Hashes a single SurgeTransitionRecord
    /// @param _record The SurgeTransitionRecord to hash
    /// @return The keccak256 hash of the encoded record
    function hashSurgeTransitionRecord(FinalityGadgetInbox.SurgeTransitionRecord calldata _record)
        external
        pure
        returns (bytes32)
    {
        return keccak256(abi.encode(_record));
    }

    /// @notice Hashes an array of SurgeTransitionRecords
    /// @param _records The SurgeTransitionRecord array to hash
    /// @return The keccak256 hash of the encoded records
    function hashSurgeTransitionRecords(FinalityGadgetInbox
                .SurgeTransitionRecord[] calldata _records)
        external
        pure
        returns (bytes32)
    {
        return keccak256(abi.encode(_records));
    }
}

