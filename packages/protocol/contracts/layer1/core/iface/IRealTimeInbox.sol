// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { LibBlobs } from "../libs/LibBlobs.sol";
import { IForcedInclusionStore } from "./IForcedInclusionStore.sol";
import { IInbox } from "./IInbox.sol";
import { ICheckpointStore } from "src/shared/signal/ICheckpointStore.sol";

/// @title IRealTimeInbox
/// @notice Interface for the real-time proving inbox.
/// @dev This inbox combines proposal and proof verification into a single atomic operation.
///      Proposer checks (lookahead, PreconfWhitelist) and bond logic are scrapped for this POC.
/// @dev Forced inclusions are supported: any user may pay a fee to enqueue a blob via
///      `saveForcedInclusion`. The proposer consumes them by setting `numForcedInclusions`
///      on `ProposeInput`/`ProposeInputV2`. If the oldest unconsumed inclusion is older than
///      `forcedInclusionDelay`, proposing reverts unless that inclusion is consumed.
/// @custom:security-contact security@nethermind.io
interface IRealTimeInbox is IForcedInclusionStore {
    /// @notice Simplified configuration for real-time proving inbox.
    struct Config {
        /// @notice The proof verifier contract (SurgeVerifier).
        address proofVerifier;
        /// @notice The signal service contract address.
        address signalService;
        /// @notice The percentage of basefee paid to coinbase.
        uint8 basefeeSharingPctg;
        /// @notice The delay in seconds after which a forced inclusion is "due" — proposing
        ///         must consume it or revert.
        uint16 forcedInclusionDelay;
        /// @notice The base fee in Gwei used for the forced-inclusion dynamic-fee curve.
        uint64 forcedInclusionFeeInGwei;
        /// @notice The pending-queue size at which the forced-inclusion fee doubles.
        ///         See `IForcedInclusionStore.getCurrentForcedInclusionFee` for formula details.
        uint64 forcedInclusionFeeDoubleThreshold;
    }

    /// @notice Input data for the propose function.
    struct ProposeInput {
        /// @notice Blob reference for proposal data.
        LibBlobs.BlobReference blobReference;
        /// @notice The number of forced inclusions to consume from the queue. Must be at least
        ///         the count of "due" inclusions or proposing will revert.
        uint16 numForcedInclusions;
        /// @notice L1 signal slots to relay to L2.
        /// @dev All signal slots will be included in the first anchor tx of the first block in POC.
        bytes32[] signalSlots;
        /// @notice The max L1 block number to verify linkage.
        /// @dev blockhash(maxAnchorBlockNumber) must be non-zero.
        uint48 maxAnchorBlockNumber;
    }

    /// @notice Input for `tentativePropose` (L2→L1→L2 flow).
    /// @dev Splits signal slots into two categories:
    ///      - existingSignals: already-sent L1→L2 signals, verified immediately
    ///      - requiredReturnSignals: L1 signals that must exist by `finalizePropose`.
    ///        These are injected as fast signals on L2 (the ZK proof commits to them)
    ///        but are produced later in the same L1 multicall by the L1 callback of a
    ///        L2→L1 bridge processMessage.
    /// @dev The anchor on L2 consumes the UNION of both lists as fast signals. The ZK
    ///      proof commits to the union hash. The inbox splits the verification timing.
    struct ProposeInputV2 {
        /// @notice Blob reference for proposal data.
        LibBlobs.BlobReference blobReference;
        /// @notice The number of forced inclusions to consume from the queue. Must be at least
        ///         the count of "due" inclusions or proposing will revert.
        uint16 numForcedInclusions;
        /// @notice L1 signals already on L1, verified at `tentativePropose` time.
        bytes32[] existingSignals;
        /// @notice L1 signals that must exist on L1 by `finalizePropose`. Produced by
        ///         the L1 callback during the multicall (between tentative and finalize).
        bytes32[] requiredReturnSignals;
        /// @notice The max L1 block number to verify linkage.
        uint48 maxAnchorBlockNumber;
    }

    /// @notice Transient proposal (not stored on-chain, only hashed).
    struct Proposal {
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

    /// @notice Commitment binding proposal, prior finalized state, and new checkpoint.
    struct Commitment {
        /// @notice Hash of the proposal being proven.
        bytes32 proposalHash;
        /// @notice Block hash of the last finalized L2 block (the proof's starting state).
        bytes32 lastFinalizedBlockHash;
        /// @notice The checkpoint data for the proven state.
        ICheckpointStore.Checkpoint checkpoint;
    }

    // ---------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------

    /// @notice Emitted on successful propose-and-prove.
    /// @param proposalHash The hash of the proposal.
    /// @param lastFinalizedBlockHash The block hash of the last finalized L2 block before this proposal.
    /// @param maxAnchorBlockNumber The L1 anchor block number.
    /// @param basefeeSharingPctg The basefee sharing percentage.
    /// @param sources Array of derivation sources.
    /// @param signalSlots Array of signal slots to be set on L2.
    /// @param checkpoint The checkpoint data saved.
    event ProposedAndProved(
        bytes32 indexed proposalHash,
        bytes32 lastFinalizedBlockHash,
        uint48 maxAnchorBlockNumber,
        uint8 basefeeSharingPctg,
        IInbox.DerivationSource[] sources,
        bytes32[] signalSlots,
        ICheckpointStore.Checkpoint checkpoint
    );

    /// @notice Emitted when the inbox is activated.
    /// @param genesisBlockHash The genesis block hash.
    event Activated(bytes32 genesisBlockHash);

    /// @notice Emitted when a tentative proposal is registered. The proposal is not
    ///         finalized until `finalizePropose` completes in the same transaction.
    /// @param proposalId Identifier for the pending proposal (equal to its proposal hash).
    /// @param requiredReturnSignalsHash Hash of the required return signal list.
    event TentativeProposed(bytes32 indexed proposalId, bytes32 requiredReturnSignalsHash);

    // ---------------------------------------------------------------
    // External Transactional Functions
    // ---------------------------------------------------------------

    /// @notice Activates the inbox with a genesis block hash.
    /// @dev Must be called by the owner before propose() can be used.
    /// @param _genesisBlockHash The genesis block hash to set as the initial finalized state.
    function activate(bytes32 _genesisBlockHash) external;

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

    /// @notice Tentatively proposes L2 blocks whose validity depends on L1 return signals
    ///         that will be produced later in the same L1 multicall.
    /// @dev Must be paired with `finalizePropose` in the same transaction. Between the two,
    ///      the pending checkpoint is saved so `Bridge.processMessage` can verify
    ///      L2→L1 signals produced by the proposed L2 block. The L1 callback triggered by
    ///      `processMessage` is expected to produce the `requiredReturnSignals`.
    /// @dev Emits `ProposedAndProved` and updates the chain head inside this call.
    ///      If `finalizePropose` reverts or is not called with matching signals, the
    ///      entire transaction reverts (by caller convention — the builder always pairs them).
    /// @param _data Encoded `ProposeInputV2` struct.
    /// @param _checkpoint Checkpoint data for the proven state.
    /// @param _proof The validity proof (SurgeVerifier SubProof[] encoded).
    /// @return proposalId_ Identifier for the pending proposal (= proposal hash).
    function tentativePropose(
        bytes calldata _data,
        ICheckpointStore.Checkpoint calldata _checkpoint,
        bytes calldata _proof
    )
        external
        returns (bytes32 proposalId_);

    /// @notice Finalizes a tentative proposal by verifying that all required return signals
    ///         have been produced on L1 during the multicall.
    /// @dev Must be called in the same transaction as `tentativePropose`. Reverts if any
    ///      required signal is missing or if the signal list does not hash-match what
    ///      was committed in `tentativePropose`.
    /// @param _requiredReturnSignals The required return signal slots (same order and
    ///        contents as in `ProposeInputV2.requiredReturnSignals`).
    function finalizePropose(bytes32[] calldata _requiredReturnSignals) external;

    // ---------------------------------------------------------------
    // External View Functions
    // ---------------------------------------------------------------

    /// @notice Returns the block hash of the last finalized L2 block.
    /// @return The last finalized block hash (bytes32(0) before activation).
    function getLastFinalizedBlockHash() external view returns (bytes32);

    /// @notice Returns the configuration parameters.
    /// @return config_ The configuration struct.
    function getConfig() external view returns (Config memory config_);
}
