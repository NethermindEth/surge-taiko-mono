// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { ECDSA } from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import { IProofVerifier } from "src/layer1/verifiers/IProofVerifier.sol";
import { LibPublicInput } from "src/layer1/verifiers/LibPublicInput.sol";

/// @title ProofVerifierDummy
/// @dev Signature-based verifier that sits behind SurgeVerifier as a sub-verifier.
///      Mirrors the public-input hashing of real verifiers (ZiskVerifier, SP1Verifier)
///      so the offchain pipeline (Raiko → Catalyst → SubProof encoding) works unchanged.
///
///      Public input: hash("VERIFY_PROOF", chainId, address(this), commitmentHash, address(0))
///      _proof:       65-byte ECDSA signature over that public input from the trusted signer.
/// @custom:security-contact security@nethermind.io
contract ProofVerifierDummy is IProofVerifier {
    using ECDSA for bytes32;

    /// @notice The trusted signer address
    address public immutable signer;

    /// @notice The L2 chain id (same as real verifiers)
    uint64 public immutable taikoChainId;

    /// @param _signer The trusted signer address
    /// @param _taikoChainId The L2 chain id used in public input hashing
    constructor(address _signer, uint64 _taikoChainId) {
        if (_signer == address(0)) revert InvalidSigner();
        signer = _signer;
        taikoChainId = _taikoChainId;
    }

    /// @inheritdoc IProofVerifier
    function verifyProof(
        uint256, /*_proposalAge*/
        bytes32 _commitmentHash,
        bytes calldata _proof
    )
        external
        view
    {
        // Hash public input identically to ZiskVerifier / SP1Verifier so
        // the same Raiko signing pipeline produces a valid signature.
        bytes32 publicInput = LibPublicInput.hashPublicInputs(
            _commitmentHash,
            address(this), // verifier contract address
            address(0), // proofSigner (not used for non-SGX)
            taikoChainId
        );

        address recoveredSigner = publicInput.recover(_proof);
        if (recoveredSigner != signer) revert InvalidSignature();
    }

    // ---------------------------------------------------------------
    // Custom Errors
    // ---------------------------------------------------------------

    error InvalidSigner();
    error InvalidSignature();
}
