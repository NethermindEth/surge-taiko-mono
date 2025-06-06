// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/layer1/based/ITaikoInbox.sol";
import "src/layer1/surge/verifiers/LibProofType.sol";

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
        returns (LibProofType.ProofType);

    /// @notice Marks the verifier for a proof type as upgradeable.
    /// @dev Should be called by the inbox contract.
    /// @param _proofType The proof type to mark as upgradeable.
    function markUpgradeable(LibProofType.ProofType _proofType) external;

    /// @notice Upgrades the verifier for a proof type.
    /// @dev Called by the owner of the parent compose verifier
    /// @param _proofType The proof type to upgrade.
    /// @param _newVerifier The address of the new verifier.
    function upgradeVerifier(LibProofType.ProofType _proofType, address _newVerifier) external;
}
