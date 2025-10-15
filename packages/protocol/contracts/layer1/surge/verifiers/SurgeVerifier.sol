// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/shared/common/EssentialContract.sol";
import "src/shared/libs/LibStrings.sol";
import "src/layer1/verifiers/IVerifier.sol";
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

    address public immutable taikoInbox;
    /// proofs come from reth client
    InternalVerifier public sgxRethVerifier;
    InternalVerifier public tdxNethermindVerifier;
    InternalVerifier public azureTdxNethermindVerifier;
    InternalVerifier public risc0RethVerifier;
    InternalVerifier public sp1RethVerifier;

    uint256[46] private __gap;

    constructor(address _taikoInbox) EssentialContract(address(0)) {
        taikoInbox = _taikoInbox;
    }

    /// @notice Initializes the contract.
    /// @param _owner The owner of this contract. msg.sender will be used if this value is zero.
    function init(
        address _owner,
        address _sgxRethVerifier,
        address _tdxNethermindVerifier,
        address _azureTdxNethermindVerifier,
        address _risc0RethVerifier,
        address _sp1RethVerifier
    )
        external
        initializer
    {
        __Essential_init(_owner);

        sgxRethVerifier.addr = _sgxRethVerifier;
        tdxNethermindVerifier.addr = _tdxNethermindVerifier;
        azureTdxNethermindVerifier.addr = _azureTdxNethermindVerifier;
        risc0RethVerifier.addr = _risc0RethVerifier;
        sp1RethVerifier.addr = _sp1RethVerifier;
    }

    /// @inheritdoc ISurgeVerifier
    function verifyProof(
        IVerifier.Context[] calldata _ctxs,
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
            IVerifier(verifier).verifyProof(_ctxs, subProofs[i].proof);

            composedProofType = composedProofType.combine(subProofs[i].proofType);
        }

        return composedProofType;
    }

    function markUpgradeable(LibProofType.ProofType _proofType) external onlyFrom(taikoInbox) {
        uint16 pt = LibProofType.ProofType.unwrap(_proofType);

        if ((pt & 0x01) != 0) {
            // SGX Reth (0b00001)
            sgxRethVerifier.upgradeable = true;
        }
        if ((pt & 0x02) != 0) {
            // TDX Reth (0b00010)
            tdxNethermindVerifier.upgradeable = true;
        }
        if ((pt & 0x04) != 0) {
            // Azure TDX Nethermind (0b00100)
            azureTdxNethermindVerifier.upgradeable = true;
        }
        if ((pt & 0x08) != 0) {
            // RISC0 Reth (0b01000)
            risc0RethVerifier.upgradeable = true;
        }
        if ((pt & 0x10) != 0) {
            // SP1 Reth (0b10000)
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
        InternalVerifier storage _verifier;
        if (_proofType.equals(LibProofType.sgxReth())) {
            _verifier = sgxRethVerifier;
        } else if (_proofType.equals(LibProofType.tdxNethermind())) {
            _verifier = tdxNethermindVerifier;
        } else if (_proofType.equals(LibProofType.azureTdxNethermind())) {
            _verifier = azureTdxNethermindVerifier;
        } else if (_proofType.equals(LibProofType.sp1Reth())) {
            _verifier = sp1RethVerifier;
        } else if (_proofType.equals(LibProofType.risc0Reth())) {
            _verifier = risc0RethVerifier;
        } else {
            revert INVALID_PROOF_TYPE();
        }

        require(_verifier.upgradeable, VERIFIER_NOT_MARKED_UPGRADEABLE());
        _verifier.addr = _newVerifier;
        _verifier.upgradeable = false;
    }

    function _getVerifierFromProofType(LibProofType.ProofType _proofType)
        internal
        view
        returns (address)
    {
        if (_proofType.equals(LibProofType.sgxReth())) {
            return sgxRethVerifier.addr;
        } else if (_proofType.equals(LibProofType.tdxNethermind())) {
            return tdxNethermindVerifier.addr;
        } else if (_proofType.equals(LibProofType.azureTdxNethermind())) {
            return azureTdxNethermindVerifier.addr;
        } else if (_proofType.equals(LibProofType.sp1Reth())) {
            return sp1RethVerifier.addr;
        } else if (_proofType.equals(LibProofType.risc0Reth())) {
            return risc0RethVerifier.addr;
        } else {
            revert INVALID_PROOF_TYPE();
        }
    }
}
