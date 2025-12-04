// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { IProofVerifier } from "src/layer1/verifiers/IProofVerifier.sol";

/// @title AnyTwoVerifierDummy
/// @notice Dummy verifier that skips proof verification - for testing only
/// @dev Same interface as AnyTwoVerifier but verifyProof doesn't route to internal verifiers
/// @custom:security-contact security@nethermind.io
contract AnyTwoVerifierDummy is IProofVerifier {
    struct SubProof {
        uint8 verifierId;
        bytes proof;
    }

    /// @notice Enum for verifier identification using stable IDs
    uint8 public constant NONE = 0;
    uint8 public constant SGX_GETH = 1;
    uint8 public constant TDX_GETH = 2;
    uint8 public constant OP = 3;
    uint8 public constant SGX_RETH = 4;
    uint8 public constant RISC0_RETH = 5;
    uint8 public constant SP1_RETH = 6;

    /// @notice Immutable verifier addresses (kept for interface compatibility)
    address public immutable sgxGethVerifier;
    address public immutable tdxGethVerifier;
    address public immutable opVerifier;
    address public immutable sgxRethVerifier;
    address public immutable risc0RethVerifier;
    address public immutable sp1RethVerifier;

    constructor(address _risc0RethVerifier, address _sp1RethVerifier) {
        // SGX not needed, set to address(0)
        sgxGethVerifier = address(0);
        tdxGethVerifier = address(0);
        opVerifier = address(0);
        sgxRethVerifier = address(0);
        risc0RethVerifier = _risc0RethVerifier;
        sp1RethVerifier = _sp1RethVerifier;
    }

    /// @inheritdoc IProofVerifier
    /// @dev DUMMY: Does not route to internal verifiers, just validates structure
    function verifyProof(uint256, bytes32, bytes calldata _proof) external view override {
        SubProof[] memory subProofs = abi.decode(_proof, (SubProof[]));
        uint256 size = subProofs.length;
        address[] memory verifiers = new address[](size);

        uint8 lastVerifierId;

        for (uint256 i; i < size; ++i) {
            uint8 verifierId = subProofs[i].verifierId;

            require(verifierId != NONE, CV_INVALID_SUB_VERIFIER());
            require(verifierId > lastVerifierId, CV_INVALID_SUB_VERIFIER_ORDER());

            // DUMMY: Skip actual verification
            verifiers[i] = _getVerifierAddress(verifierId);
            lastVerifierId = verifierId;
        }

        require(_areVerifiersSufficient(verifiers), CV_VERIFIERS_INSUFFICIENT());
    }

    /// @notice Returns the verifier address for a given verifier ID
    /// @param _verifierId The verifier ID to query
    /// @return The address of the verifier (or address(0) if invalid)
    function getVerifierAddress(uint8 _verifierId) public view returns (address) {
        return _getVerifierAddress(_verifierId);
    }

    function _getVerifierAddress(uint8 _verifierId) internal view returns (address) {
        if (_verifierId == SGX_GETH) return sgxGethVerifier;
        if (_verifierId == TDX_GETH) return tdxGethVerifier;
        if (_verifierId == OP) return opVerifier;
        if (_verifierId == SGX_RETH) return sgxRethVerifier;
        if (_verifierId == RISC0_RETH) return risc0RethVerifier;
        if (_verifierId == SP1_RETH) return sp1RethVerifier;
        return address(0);
    }

    /// @dev Checks if the provided verifiers are sufficient
    /// @notice Valid combinations: [RISC0_RETH, SP1_RETH]
    function _areVerifiersSufficient(address[] memory _verifiers) internal view returns (bool) {
        if (_verifiers.length != 2) return false;

        // Only valid combination (without SGX): [RISC0_RETH, SP1_RETH]
        return _verifiers[0] == risc0RethVerifier && _verifiers[1] == sp1RethVerifier;
    }

    error CV_INVALID_SUB_VERIFIER();
    error CV_INVALID_SUB_VERIFIER_ORDER();
    error CV_VERIFIERS_INSUFFICIENT();
}

