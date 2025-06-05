// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/shared/common/EssentialContract.sol";
import "src/shared/libs/LibStrings.sol";
import "./ISurgeVerifier.sol";
import "./LibProofType.sol";

/// @title SurgeVerifier
/// @notice This contract is a verifier that composes multiple sub-verifiers to validate
/// proofs.
/// It ensures that a set of sub-proofs are verified by their respective verifiers before
/// considering the overall proof as valid.
/// @custom:security-contact security@nethermind.io
contract SurgeVerifier is EssentialContract, ISurgeVerifier {
    using LibProofType for LibProofType.ProofType;

    struct SubProof {
        // This is a single proof type i.e SGX_RETH / SP1_RETH / RISC0_RETH / TDX_RETH
        LibProofType.ProofType proofType;
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

    uint256[44] private __gap;

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

    error INVALID_PROOF_TYPE();
    error UPGRADE_NOT_SUPPORTED_BY_PROOF_TYPE();
    error VERIFIER_NOT_MARKED_UPGRADEABLE();

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
        LibProofType.ProofType composedProofType = LibProofType.empty();
        uint256 size = subProofs.length;

        for (uint256 i; i < size; ++i) {
            address verifier = _getVerifierFromProofType(subProofs[i].proofType);
            ISurgeVerifier(verifier).verifyProof(_ctxs, subProofs[i].proof);

            composedProofType = composedProofType.combine(subProofs[i].proofType);
        }

        return composedProofType;
    }

    function markUpgradeable(LibProofType.ProofType _proofType) external onlyFrom(taikoInbox) {
        uint16 pt = LibProofType.ProofType.unwrap(_proofType);

        if ((pt & 0x01) != 0) {
            // SGX Reth (0b0001)
            sgxRethVerifier.upgradeable = true;
        }
        if ((pt & 0x02) != 0) {
            // TDX Reth (0b0010)
            tdxGethVerifier.upgradeable = true;
        }
        if ((pt & 0x04) != 0) {
            // RISC0 Reth (0b0100)
            risc0RethVerifier.upgradeable = true;
        }
        if ((pt & 0x08) != 0) {
            // SP1 Reth (0b1000)
            sp1RethVerifier.upgradeable = true;
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
        if (_proofType.equals(LibProofType.sgxReth())) {
            _verifier = sgxRethVerifier;
        } else if (_proofType.equals(LibProofType.sp1Reth())) {
            _verifier = sp1RethVerifier;
        } else if (_proofType.equals(LibProofType.risc0Reth())) {
            _verifier = risc0RethVerifier;
        } else {
            revert UPGRADE_NOT_SUPPORTED_BY_PROOF_TYPE();
        }

        require(_verifier.upgradeable, VERIFIER_NOT_MARKED_UPGRADEABLE());
        _verifier.addr = _newVerifier;
    }

    function _getVerifierFromProofType(LibProofType.ProofType _proofType)
        internal
        view
        returns (address)
    {
        if (_proofType.equals(LibProofType.sgxReth())) {
            return sgxRethVerifier.addr;
        } else if (_proofType.equals(LibProofType.sp1Reth())) {
            return sp1RethVerifier.addr;
        } else if (_proofType.equals(LibProofType.risc0Reth())) {
            return risc0RethVerifier.addr;
        } else {
            revert INVALID_PROOF_TYPE();
        }
    }
}
