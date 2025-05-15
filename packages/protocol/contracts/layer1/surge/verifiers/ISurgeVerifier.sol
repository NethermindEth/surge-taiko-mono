// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/layer1/based/ITaikoInbox.sol";

/// @title ISurgeVerifier
/// @notice Defines the function that handles proof verification.
/// @custom:security-contact security@nethermind.io
interface ISurgeVerifier {
    struct Context {
        uint64 batchId;
        bytes32 metaHash;
        ITaikoInbox.Transition transition;
    }

    /// @notice Verifies multiple proofs. This function must throw if the proof cannot be verified.
    /// @param _ctxs The array of contexts for the proof verifications.
    /// @param _proof The batch proof to verify.
    /// @return proofType The type of proof.
    function verifyProof(
        Context[] calldata _ctxs,
        bytes calldata _proof
    )
        external
        returns (ITaikoInbox.ProofType);
}
