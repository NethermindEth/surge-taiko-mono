// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { LibBlobs } from "../libs/LibBlobs.sol";
import { IInbox } from "./IInbox.sol";
import { ICheckpointStore } from "src/shared/signal/ICheckpointStore.sol";

/// @title IRealTimeInbox
/// @notice Interface for the real-time proving inbox.
/// @dev This inbox combines proposal and proof verification into a single atomic operation.
///      Proposer checks (lookahead, PreconfWhitelist) and bond logic are scrapped for this POC.
/// @custom:security-contact security@nethermind.io
interface IRealTimeInbox {
    /// @notice Simplified configuration for real-time proving inbox.
    struct Config {
        /// @notice The proof verifier contract (SurgeVerifier).
        address proofVerifier;
        /// @notice The signal service contract address.
        address signalService;
        /// @notice The percentage of basefee paid to coinbase.
        uint8 basefeeSharingPctg;
    }

    /// @notice Input data for the propose function.
    struct ProposeInput {
        /// @notice Blob reference for proposal data.
        LibBlobs.BlobReference blobReference;
        /// @notice L1 signal slots to relay to L2.
        /// @dev All signal slots will be included in the first anchor tx of the first block in POC.
        bytes32[] signalSlots;
        /// @notice The max L1 block number to verify linkage.
        /// @dev blockhash(maxAnchorBlockNumber) must be non-zero.
        uint48 maxAnchorBlockNumber;
    }

    /// @notice Transient proposal (not stored on-chain, only hashed).
    struct Proposal {
        /// @notice Hash of the parent proposal (zero for genesis).
        bytes32 parentProposalHash;
        /// @notice The height of highest anchor block.
        uint48 maxAnchorBlockNumber;
        /// @notice The hash of the highest anchor block.
        bytes32 maxAnchorBlockHash;
        /// @notice The percentage of base fee paid to coinbase.
        uint8 basefeeSharingPctg;
        /// @notice Array of derivation sources.
        IInbox.DerivationSource[] sources;
        /// @notice Hash of signal slots to be set on L2.
        bytes32 signalSlotsHash;
    }

    /// @notice Simplified commitment containing only proposal hash and checkpoint.
    struct Commitment {
        /// @notice Hash of the proposal being proven.
        bytes32 proposalHash;
        /// @notice The checkpoint data for the proven state.
        ICheckpointStore.Checkpoint checkpoint;
    }

    // ---------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------

    /// @notice Emitted on successful propose-and-prove.
    /// @param proposalHash The hash of the proposal.
    /// @param parentProposalHash The hash of the parent proposal.
    /// @param maxAnchorBlockNumber The L1 anchor block number.
    /// @param basefeeSharingPctg The basefee sharing percentage.
    /// @param sources Array of derivation sources.
    /// @param signalSlots Array of signal slots to be set on L2.
    /// @param checkpoint The checkpoint data saved.
    event ProposedAndProved(
        bytes32 indexed proposalHash,
        bytes32 parentProposalHash,
        uint48 maxAnchorBlockNumber,
        uint8 basefeeSharingPctg,
        IInbox.DerivationSource[] sources,
        bytes32[] signalSlots,
        ICheckpointStore.Checkpoint checkpoint
    );

    /// @notice Emitted when the inbox is activated.
    /// @param genesisProposalHash The genesis proposal hash.
    event Activated(bytes32 genesisProposalHash);

    // ---------------------------------------------------------------
    // External Transactional Functions
    // ---------------------------------------------------------------

    /// @notice Activates the inbox with a genesis proposal hash.
    /// @dev Must be called by the owner before propose() can be used.
    /// @param _genesisProposalHash The genesis proposal hash to set.
    function activate(bytes32 _genesisProposalHash) external;

    /// @notice Proposes new L2 blocks with real-time proof verification.
    /// @dev Combines proposal submission and proof verification into a single atomic operation.
    /// @param _data The encoded ProposeInput struct.
    /// @param _checkpoint The checkpoint data for the proven state.
    /// @param _proof The validity proof (SurgeVerifier SubProof[] encoded).
    function propose(
        bytes calldata _data,
        ICheckpointStore.Checkpoint calldata _checkpoint,
        bytes calldata _proof
    )
        external;

    // ---------------------------------------------------------------
    // External View Functions
    // ---------------------------------------------------------------

    /// @notice Returns the hash of the last accepted proposal.
    /// @return The last proposal hash (bytes32(0) before activation).
    function getLastProposalHash() external view returns (bytes32);

    /// @notice Returns the configuration parameters.
    /// @return config_ The configuration struct.
    function getConfig() external view returns (Config memory config_);
}
