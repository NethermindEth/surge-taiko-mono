// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/shared/common/EssentialContract.sol";
import "src/shared/libs/LibStrings.sol";
import "./ISurgeVerifier.sol";
import "./LibProofType.sol";

/// @title SurgeComposeVerifier
/// @notice This contract is an abstract verifier that composes multiple sub-verifiers to validate
/// proofs.
/// It ensures that a set of sub-proofs are verified by their respective verifiers before
/// considering the overall proof as valid.
/// @custom:security-contact security@nethermind.io
abstract contract SurgeComposeVerifier is EssentialContract, ISurgeVerifier {
    uint256[50] private __gap;

    struct SubProof {
        address verifier;
        bytes proof;
    }

    struct Verifier {
        bool upgradeable;
        address addr;
    }

    address public immutable taikoInbox;
    /// The sgx/tdx-GethVerifier is the core verifier required in every proof.
    /// All other proofs share its status root, despite different public inputs
    /// due to different verification types.
    /// proofs come from geth client
    Verifier public sgxGethVerifier;
    Verifier public tdxGethVerifier;
    /// op for test purpose
    Verifier public opVerifier;
    /// proofs come from reth client
    Verifier public sgxRethVerifier;
    Verifier public risc0RethVerifier;
    Verifier public sp1RethVerifier;

    constructor(
        address _taikoInbox,
        address _sgxGethVerifier,
        address _tdxGethVerifier,
        address _opVerifier,
        address _sgxRethVerifier,
        address _risc0RethVerifier,
        address _sp1RethVerifier
    )
        EssentialContract(address(0))
    {
        taikoInbox = _taikoInbox;
        sgxGethVerifier.addr = _sgxGethVerifier;
        tdxGethVerifier.addr = _tdxGethVerifier;
        opVerifier.addr = _opVerifier;
        sgxRethVerifier.addr = _sgxRethVerifier;
        risc0RethVerifier.addr = _risc0RethVerifier;
        sp1RethVerifier.addr = _sp1RethVerifier;
    }

    error SCV_INVALID_SUB_VERIFIER();
    error SCV_INVALID_SUB_VERIFIER_ORDER();
    error SCV_INVALID_PROOF_TYPE();
    error SCV_UPGRADE_NOT_SUPPORTED_BY_PROOF_TYPE();
    error SCV_VERIFIER_NOT_MARKED_UPGRADEABLE();

    /// @notice Initializes the contract.
    /// @param _owner The owner of this contract. msg.sender will be used if this value is zero.
    function init(address _owner) external initializer {
        __Essential_init(_owner);
    }

    /// @inheritdoc ISurgeVerifier
    function verifyProof(
        Context[] calldata _ctxs,
        bytes calldata _proof
    )
        external
        onlyFrom(taikoInbox)
        returns (LibProofType.ProofType)
    {
        SubProof[] memory subProofs = abi.decode(_proof, (SubProof[]));
        uint256 size = subProofs.length;
        address[] memory verifiers = new address[](size);

        address verifier;

        for (uint256 i; i < size; ++i) {
            require(subProofs[i].verifier != address(0), SCV_INVALID_SUB_VERIFIER());
            require(subProofs[i].verifier > verifier, SCV_INVALID_SUB_VERIFIER_ORDER());

            verifier = subProofs[i].verifier;
            ISurgeVerifier(verifier).verifyProof(_ctxs, subProofs[i].proof);

            verifiers[i] = verifier;
        }

        LibProofType.ProofType proofType = getProofTypeFromVerifiers(verifiers);
        require(proofType != LibProofType.ProofType.INVALID, SCV_INVALID_PROOF_TYPE());

        return proofType;
    }

    function markUpgradeable(LibProofType.ProofType _proofType) external onlyFrom(taikoInbox) {
        if (_proofType == LibProofType.ProofType.SGX) {
            sgxRethVerifier.upgradeable = true;
        } else if (_proofType == LibProofType.ProofType.SP1) {
            sp1RethVerifier.upgradeable = true;
        } else if (_proofType == LibProofType.ProofType.RISC0) {
            risc0RethVerifier.upgradeable = true;
        } else {
            revert SCV_UPGRADE_NOT_SUPPORTED_BY_PROOF_TYPE();
        }
    }

    function upgradeVerifier(
        LibProofType.ProofType _proofType,
        address _newVerifier
    )
        external
        onlyOwner
    {
        Verifier storage _verifier;
        if (_proofType == LibProofType.ProofType.SGX) {
            _verifier = sgxRethVerifier;
        } else if (_proofType == LibProofType.ProofType.SP1) {
            _verifier = sp1RethVerifier;
        } else if (_proofType == LibProofType.ProofType.RISC0) {
            _verifier = risc0RethVerifier;
        } else {
            revert SCV_UPGRADE_NOT_SUPPORTED_BY_PROOF_TYPE();
        }

        require(_verifier.upgradeable, SCV_VERIFIER_NOT_MARKED_UPGRADEABLE());
        _verifier.addr = _newVerifier;
    }

    function getProofTypeFromVerifiers(address[] memory _verifiers)
        internal
        view
        virtual
        returns (LibProofType.ProofType);
}
