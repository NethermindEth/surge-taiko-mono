// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { IProofVerifier } from "../verifiers/IProofVerifier.sol";
import { Ownable2Step } from "@openzeppelin/contracts/access/Ownable2Step.sol";

/// @title SurgeVerifier
/// @notice Routes proof verification to internal verifiers
/// @custom:security-contact security@nethermind.io
contract SurgeVerifier is Ownable2Step {
    address public immutable inbox;

    struct InternalVerifier {
        /// @dev Address of the proof specific verifier, eg: SP1, RISC0, etc.
        address addr;
        /// @dev When `true` the timelock on the security council can be bypassed
        /// to allow for instantly upgrading this verifier's address
        bool allowInstantUpgrade;
    }

    struct Proof {
        /// @dev Unique id of the verifier that can be resolved from `LibVerifierId.sol`
        uint8 verifierId;
        /// @dev The cryptographic proof
        bytes data;
    }

    /// @notice Mapping from verifier id to an internal verifier contract that implements IProofVerifier
    mapping(uint8 verifierId => InternalVerifier verifier) public verifiers;

    /// @dev Emitted when a verifier is updated
    /// @param verifierId The id of the verifier
    /// @param oldVerifier The previous verifier address
    /// @param newVerifier The new verifier address
    event VerifierUpdated(uint8 indexed verifierId, address oldVerifier, address newVerifier);

    /// @param _owner The initial owner
    constructor(address _inbox, address _owner) {
        inbox = _inbox;
        _transferOwnership(_owner);
    }

    /// @notice Sets or updates the verifier for a given verifier id
    /// @param _verifierId The id used to route proofs
    /// @param _verifierAddr The verifier contract address (must implement IProofVerifier)
    function setVerifier(uint8 _verifierId, address _verifierAddr) external onlyOwner {
        address oldVerifierAddr = verifiers[_verifierId].addr;
        verifiers[_verifierId] = InternalVerifier(_verifierAddr, false);
        emit VerifierUpdated(_verifierId, oldVerifierAddr, _verifierAddr);
    }

    /// @notice Instantly upgrades the address of the internal verifier for a given verifier ID,
    ///         bypassing the timelock if `allowInstantUpgrade` is true for that verifier.
    /// @dev Only callable by the contract owner and only if the old verifier allows instant upgrade.
    /// @param _verifierId The ID of the internal verifier to upgrade.
    /// @param _verifierAddr The new verifier contract address (must implement IProofVerifier).
    function setVerifierInstant(
        uint8 _verifierId,
        address _verifierAddr
    )
        external
        onlyOwner
    {
        InternalVerifier memory oldVerifier = verifiers[_verifierId];
        require(oldVerifier.addr != address(0), InvalidVerifierId());
        require(oldVerifier.allowInstantUpgrade, InstantUpgradeNotAllowed());
        verifiers[_verifierId] = InternalVerifier(_verifierAddr, false);
        emit VerifierUpdated(_verifierId, oldVerifier.addr, _verifierAddr);
    }

    /// @notice Marks a verifier as upgradeable (allows instant upgrade) or not.
    /// @param _verifierId The id of the verifier to update.
    /// @param _allowInstantUpgrade Whether instant upgrade should be allowed for this verifier.
    function markVerifierUpgradeable(
        uint8 _verifierId,
        bool _allowInstantUpgrade
    )
        external
    {
        require(msg.sender == inbox, CallerIsNotInbox());
        InternalVerifier storage verifier = verifiers[_verifierId];
        require(verifier.addr != address(0), InvalidVerifierId());
        verifier.allowInstantUpgrade = _allowInstantUpgrade;
    }

    /// @notice Verifies a validity proof for a state transition
    /// @dev This function must revert if the proof is invalid
    /// @dev This is a presumed extension of `IProofVerifier` and returns the verifier id.
    /// @param _proposalAge The age in seconds of the proposal being proven. Only set for
    ///        single-proposal proofs (calculated as block.timestamp - proposal.timestamp).
    ///        For multi-proposal batches, this is always 0, meaning "not applicable".
    ///        Verifiers should interpret _proposalAge == 0 as "not applicable" rather than
    ///        "instant proof". This parameter enables age-based verification logic, such as
    ///        detecting and handling prover-killer proposals differently.
    /// @param _transitionsHash The hash of the transitions to verify
    /// @param _proof The proof data for the transitions
    /// @return verifierId_ The verifier id extracted from the proof
    function verifyProof(
        uint256 _proposalAge,
        bytes32 _transitionsHash,
        bytes calldata _proof
    )
        external
        view
        returns (uint8 verifierId_)
    {
        Proof memory proof = abi.decode(_proof, (Proof));

        verifierId_ = proof.verifierId;
        address verifierAddr = verifiers[verifierId_].addr;
        if (verifierAddr == address(0)) revert InvalidVerifierId();

        IProofVerifier(verifierAddr).verifyProof(_proposalAge, _transitionsHash, proof.data);
    }

    // ---------------------------------------------------------------
    // Custom Errors
    // ---------------------------------------------------------------
    error CallerIsNotInbox();
    error InstantUpgradeNotAllowed();
    error InvalidVerifierId();
    error VerifierNotSet();
}

