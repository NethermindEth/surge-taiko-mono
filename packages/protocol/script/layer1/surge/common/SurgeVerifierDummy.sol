// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { Ownable2Step } from "@openzeppelin/contracts/access/Ownable2Step.sol";
import { LibProofBitmap } from "src/layer1/Surge/LibProofBitmap.sol";

/// @title SurgeVerifierDummy
/// @notice Dummy verifier that skips proof verification - for testing only
/// @dev Same as SurgeVerifier but verifyProof doesn't route to internal verifiers
/// @custom:security-contact security@nethermind.io
contract SurgeVerifierDummy is Ownable2Step {
    using LibProofBitmap for LibProofBitmap.ProofBitmap;

    address public immutable inbox;

    struct InternalVerifier {
        /// @dev Address of the proof specific verifier, eg: SP1, RISC0, etc.
        address addr;
        /// @dev When `true` the timelock on the security council can be bypassed
        /// to allow for instantly upgrading this verifier's address
        bool allowInstantUpgrade;
    }

    struct SubProof {
        /// @dev The bit flag of the proof type that can be resolved from `LibProofBitmap.sol`
        LibProofBitmap.ProofBitmap proofBitFlag;
        /// @dev The cryptographic proof
        bytes data;
    }

    /// @notice Mapping from bit flag to an internal verifier contract that implements IProofVerifier
    mapping(LibProofBitmap.ProofBitmap proofBitFlag => InternalVerifier verifier) public verifiers;

    /// @dev Emitted when a verifier is updated
    /// @param proofBitFlag The proof bit flag of the verifier
    /// @param oldVerifier The previous verifier address
    /// @param newVerifier The new verifier address
    event VerifierUpdated(
        LibProofBitmap.ProofBitmap indexed proofBitFlag, address oldVerifier, address newVerifier
    );

    /// @param _owner The initial owner
    constructor(address _inbox, address _owner) {
        inbox = _inbox;
        _transferOwnership(_owner);
    }

    /// @notice Sets or updates the verifier for a given proof bit flag
    /// @param _proofBitFlag The proof bit flag used to route proofs
    /// @param _verifierAddr The verifier contract address (must implement IProofVerifier)
    function setVerifier(
        LibProofBitmap.ProofBitmap _proofBitFlag,
        address _verifierAddr
    )
        external
        onlyOwner
    {
        address oldVerifierAddr = verifiers[_proofBitFlag].addr;
        verifiers[_proofBitFlag] = InternalVerifier(_verifierAddr, false);
        emit VerifierUpdated(_proofBitFlag, oldVerifierAddr, _verifierAddr);
    }

    /// @notice Instantly upgrades the address of the internal verifier for a given proof bit flag,
    ///         bypassing the timelock if `allowInstantUpgrade` is true for that verifier.
    /// @dev Only callable by the contract owner and only if the old verifier allows instant upgrade.
    /// @param _proofBitFlag The proof bit flag of the internal verifier to upgrade.
    /// @param _verifierAddr The new verifier contract address (must implement IProofVerifier).
    function setVerifierInstant(
        LibProofBitmap.ProofBitmap _proofBitFlag,
        address _verifierAddr
    )
        external
        onlyOwner
    {
        InternalVerifier memory oldVerifier = verifiers[_proofBitFlag];
        require(oldVerifier.addr != address(0), Surge_InvalidProofBitFlag());
        require(oldVerifier.allowInstantUpgrade, Surge_InstantUpgradeNotAllowed());
        verifiers[_proofBitFlag] = InternalVerifier(_verifierAddr, false);
        emit VerifierUpdated(_proofBitFlag, oldVerifier.addr, _verifierAddr);
    }

    /// @notice Marks verifiers as upgradeable (allows instant upgrade) or not, according to bits set in the provided bitmap.
    /// @param _proofBitmap The full bitmap indicating which verifiers to update.
    /// @param _allowInstantUpgrade Whether instant upgrade should be allowed for these verifiers.
    function markVerifiersUpgradeable(
        LibProofBitmap.ProofBitmap _proofBitmap,
        bool _allowInstantUpgrade
    )
        external
    {
        require(msg.sender == inbox, Surge_CallerIsNotInbox());

        uint8 flags = _proofBitmap.toUint8();
        for (uint8 i = 0; i < 8; ++i) {
            if ((flags & (1 << i)) != 0) {
                LibProofBitmap.ProofBitmap bit = LibProofBitmap.ProofBitmap.wrap(uint8(1 << i));
                InternalVerifier storage verifier = verifiers[bit];
                require(verifier.addr != address(0), Surge_InvalidProofBitFlag());
                verifier.allowInstantUpgrade = _allowInstantUpgrade;
            }
        }
    }

    /// @notice Returns merged bitmap without verifying proofs
    /// @dev DUMMY: Does not route to internal verifiers, just merges bitmaps
    /// @param _proof The proof data containing an array of sub proofs
    /// @return mergedBitmap_ The merged bitmap of all proof flags
    function verifyProof(
        uint256,
        bytes32,
        bytes calldata _proof
    )
        external
        pure
        returns (LibProofBitmap.ProofBitmap mergedBitmap_)
    {
        SubProof[] memory subProofs = abi.decode(_proof, (SubProof[]));

        for (uint256 i; i < subProofs.length; ++i) {
            mergedBitmap_ = mergedBitmap_.merge(subProofs[i].proofBitFlag);
        }
    }

    // ---------------------------------------------------------------
    // Custom Errors
    // ---------------------------------------------------------------

    error Surge_CallerIsNotInbox();
    error Surge_InstantUpgradeNotAllowed();
    error Surge_InvalidProofBitFlag();
}
