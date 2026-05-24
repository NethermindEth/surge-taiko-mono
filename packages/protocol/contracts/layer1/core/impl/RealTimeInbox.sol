// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { IForcedInclusionStore } from "../iface/IForcedInclusionStore.sol";
import { IInbox } from "../iface/IInbox.sol";
import { IRealTimeInbox } from "../iface/IRealTimeInbox.sol";
import { LibBlobs } from "../libs/LibBlobs.sol";
import { LibForcedInclusion } from "../libs/LibForcedInclusion.sol";
import { SurgeVerifier } from "src/layer1/surge/SurgeVerifier.sol";
import { EssentialContract } from "src/shared/common/EssentialContract.sol";
import { LibAddress } from "src/shared/libs/LibAddress.sol";
import { ICheckpointStore } from "src/shared/signal/ICheckpointStore.sol";
import { ISignalService } from "src/shared/signal/ISignalService.sol";

/// @title RealTimeInbox
/// @notice Inbox contract that combines proposal and proof verification into a single atomic
/// operation. Each call to `propose()` submits a proposal, verifies a ZK proof, and finalizes
/// the state in one transaction.
/// @dev Proposer checks (lookahead, PreconfWhitelist), bond logic, ring buffer storage, and
///      prover whitelist are all scrapped for this real-time proving POC. Forced inclusions
///      are retained: external users enqueue blobs via `saveForcedInclusion`; the proposer
///      consumes them via `numForcedInclusions` on the propose input.
/// @dev WARNING: This contract is vulnerable to proposal frontrunning. A malicious actor can observe
///      a pending `propose()` transaction in the mempool and submit the same proposal with their own
///      address to steal credit. In production, an `actualProver` field (msg.sender) should be included
///      in the Commitment hash so that the proof is bound to a specific sender and cannot be replayed.
/// @custom:security-contact security@nethermind.io
contract RealTimeInbox is IRealTimeInbox, EssentialContract {
    using LibAddress for address;
    using LibForcedInclusion for LibForcedInclusion.Storage;
    // ---------------------------------------------------------------
    // Immutable Variables
    // ---------------------------------------------------------------

    /// @notice The proof verifier contract.
    SurgeVerifier internal immutable _proofVerifier;

    /// @notice Signal service responsible for checkpoints and signal relay.
    ISignalService internal immutable _signalService;

    /// @notice The percentage of basefee paid to coinbase.
    uint8 internal immutable _basefeeSharingPctg;

    /// @notice The delay in seconds after which a forced inclusion is "due".
    uint16 internal immutable _forcedInclusionDelay;

    /// @notice The base fee in Gwei for the forced-inclusion dynamic-fee curve.
    uint64 internal immutable _forcedInclusionFeeInGwei;

    /// @notice Queue size at which the forced-inclusion fee doubles.
    uint64 internal immutable _forcedInclusionFeeDoubleThreshold;

    // ---------------------------------------------------------------
    // State Variables
    // ---------------------------------------------------------------

    /// @notice Block hash of the last finalized L2 block. Serves as the chain head.
    bytes32 public lastFinalizedBlockHash;

    /// @dev Storage for the forced-inclusion FIFO queue. Uses 2 storage slots.
    LibForcedInclusion.Storage private _forcedInclusionStorage;

    uint256[47] private __gap;

    // ---------------------------------------------------------------
    // Transient Storage for Pending Proposals
    // ---------------------------------------------------------------
    //
    // Layout: one pending proposal per transaction (by convention).
    // tentativePropose -> sets these. finalizePropose -> reads + clears.
    // Solidity's automatic revert semantics ensure that if finalize reverts,
    // the entire tx (including tentative's effects) is rolled back.
    //
    // Storing only a hash of the required return signals is sufficient — the
    // caller of finalizePropose re-provides the signal list, which we verify
    // against the stored hash before checking each signal's existence on L1.

    /// @notice Marker for pending proposal presence (1 = pending, 0 = none).
    uint256 private transient _pendingExists;
    /// @notice Proposal hash used as the pending proposal identifier.
    bytes32 private transient _pendingProposalId;
    /// @notice Hash of the required-return signal slots recorded by tentativePropose.
    bytes32 private transient _pendingRequiredSignalsHash;

    // ---------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------

    /// @dev Initializes immutable configuration.
    /// @param _config Configuration struct.
    constructor(Config memory _config) {
        require(_config.proofVerifier != address(0), "config: proofVerifier");
        require(_config.signalService != address(0), "config: signalService");
        require(_config.basefeeSharingPctg <= 100, "config: basefeeSharingPctg");
        require(_config.forcedInclusionFeeDoubleThreshold > 0, "config: feeDoubleThreshold");

        _proofVerifier = SurgeVerifier(_config.proofVerifier);
        _signalService = ISignalService(_config.signalService);
        _basefeeSharingPctg = _config.basefeeSharingPctg;
        _forcedInclusionDelay = _config.forcedInclusionDelay;
        _forcedInclusionFeeInGwei = _config.forcedInclusionFeeInGwei;
        _forcedInclusionFeeDoubleThreshold = _config.forcedInclusionFeeDoubleThreshold;
    }

    // ---------------------------------------------------------------
    // External Functions
    // ---------------------------------------------------------------

    /// @notice Initializes the owner of the inbox.
    /// @param _owner The owner of this contract.
    function init(address _owner) external initializer {
        __Essential_init(_owner);
    }

    /// @inheritdoc IRealTimeInbox
    function activate(bytes32 _genesisBlockHash) external onlyOwner {
        require(lastFinalizedBlockHash == bytes32(0), AlreadyActivated());
        require(_genesisBlockHash != bytes32(0), InvalidGenesisBlockHash());

        lastFinalizedBlockHash = _genesisBlockHash;
        emit Activated(_genesisBlockHash);
    }

    /// @inheritdoc IRealTimeInbox
    function propose(
        bytes calldata _data,
        ICheckpointStore.Checkpoint calldata _checkpoint,
        bytes calldata _proof
    )
        external
        nonReentrant
    {
        require(lastFinalizedBlockHash != bytes32(0), NotActivated());

        // Capture current chain head before it is updated
        bytes32 prevFinalizedBlockHash = lastFinalizedBlockHash;

        // Build proposal from input and get its hash
        (bytes32 proposalHash, Proposal memory proposal, bytes32[] memory signalSlots) =
            _buildProposal(_data);

        // Verify proof and finalize (updates lastFinalizedBlockHash)
        _verifyAndFinalize(proposalHash, prevFinalizedBlockHash, _checkpoint, _proof);

        // Emit event with raw signal slots for driver derivation
        emit ProposedAndProved(
            proposalHash,
            prevFinalizedBlockHash,
            proposal.maxAnchorBlockNumber,
            proposal.basefeeSharingPctg,
            proposal.sources,
            signalSlots,
            _checkpoint
        );
    }

    /// @inheritdoc IRealTimeInbox
    function tentativePropose(
        bytes calldata _data,
        ICheckpointStore.Checkpoint calldata _checkpoint,
        bytes calldata _proof
    )
        external
        nonReentrant
        returns (bytes32 proposalId_)
    {
        require(lastFinalizedBlockHash != bytes32(0), NotActivated());
        require(_pendingExists == 0, PendingProposalAlreadyExists());

        bytes32 prevFinalizedBlockHash = lastFinalizedBlockHash;

        // Build the proposal. The anchor on L2 consumes the UNION of existingSignals
        // and requiredReturnSignals as fast signals; the proof commits to this union.
        // Only the existingSignals are verified here — requiredReturnSignals are
        // verified in finalizePropose after the L1 callback produces them.
        (
            bytes32 proposalHash,
            Proposal memory proposal,
            bytes32[] memory unionSignalSlots,
            bytes32 requiredSignalsHash
        ) = _buildProposalV2(_data);

        // Verify proof + save checkpoint + update chain head.
        // Saving the checkpoint here is what lets Bridge.processMessage (called later
        // in this multicall for L2→L1 signals) verify storage proofs against the L2
        // state root. If finalizePropose reverts at the end, the entire transaction
        // is rolled back, unwinding every effect of this call.
        _verifyAndFinalize(proposalHash, prevFinalizedBlockHash, _checkpoint, _proof);

        // Emit the canonical event so the driver can derive the L2 block as usual.
        // If the tx reverts (e.g. finalizePropose fails), logs are discarded.
        emit ProposedAndProved(
            proposalHash,
            prevFinalizedBlockHash,
            proposal.maxAnchorBlockNumber,
            proposal.basefeeSharingPctg,
            proposal.sources,
            unionSignalSlots,
            _checkpoint
        );

        // Record pending state in transient storage. Finalize reads + clears this.
        _pendingExists = 1;
        _pendingProposalId = proposalHash;
        _pendingRequiredSignalsHash = requiredSignalsHash;

        emit TentativeProposed(proposalHash, requiredSignalsHash);

        proposalId_ = proposalHash;
    }

    /// @inheritdoc IRealTimeInbox
    function finalizePropose(bytes32[] calldata _requiredReturnSignals) external {
        require(_pendingExists == 1, NoPendingProposal());

        // Verify the provided list matches the hash committed by tentativePropose.
        // This binds the caller to the exact signal set the L2 block was proven with.
        bytes32 providedHash = hashSignalSlots(_requiredReturnSignals);
        require(providedHash == _pendingRequiredSignalsHash, RequiredSignalsMismatch());

        // Verify every required signal now exists on L1. These are produced by the
        // L1 callback of Bridge.processMessage earlier in this multicall.
        for (uint256 i; i < _requiredReturnSignals.length; ++i) {
            require(
                _signalService.isSignalSent(_requiredReturnSignals[i]),
                RequiredSignalNotSent(_requiredReturnSignals[i])
            );
        }

        // Clear pending state (transient — also auto-clears at tx end).
        _pendingExists = 0;
        _pendingProposalId = bytes32(0);
        _pendingRequiredSignalsHash = bytes32(0);
    }

    /// @inheritdoc IForcedInclusionStore
    /// @dev Reverts before activation: forced inclusions cannot be enqueued until the inbox
    /// is activated, since their `blobSlice.timestamp` would not be meaningfully comparable
    /// to anything pre-genesis.
    function saveForcedInclusion(LibBlobs.BlobReference memory _blobReference)
        external
        payable
        nonReentrant
    {
        require(lastFinalizedBlockHash != bytes32(0), NotActivated());

        uint256 refund = _forcedInclusionStorage.saveForcedInclusion(
            _forcedInclusionFeeInGwei, _forcedInclusionFeeDoubleThreshold, _blobReference
        );

        // Refund excess payment to the sender
        if (refund > 0) {
            msg.sender.sendEtherAndVerify(refund);
        }
    }

    // ---------------------------------------------------------------
    // Encoding / Decoding / Hashing Functions
    // ---------------------------------------------------------------

    /// @notice Encodes a ProposeInput struct into bytes.
    /// @param _input The ProposeInput to encode.
    /// @return encoded_ The ABI-encoded bytes.
    function encodeProposeInput(ProposeInput calldata _input)
        public
        pure
        returns (bytes memory encoded_)
    {
        return abi.encode(_input);
    }

    /// @notice Decodes bytes into a ProposeInput struct.
    /// @param _data The ABI-encoded ProposeInput.
    /// @return input_ The decoded ProposeInput.
    function decodeProposeInput(bytes calldata _data)
        public
        pure
        returns (ProposeInput memory input_)
    {
        return abi.decode(_data, (ProposeInput));
    }

    /// @notice Hashes a Proposal struct.
    /// @param _proposal The Proposal to hash.
    /// @return The keccak256 hash.
    function hashProposal(Proposal memory _proposal) public pure returns (bytes32) {
        return keccak256(abi.encode(_proposal));
    }

    /// @notice Hashes a Commitment struct.
    /// @param _commitment The Commitment to hash.
    /// @return The keccak256 hash.
    function hashCommitment(Commitment memory _commitment) public pure returns (bytes32) {
        return keccak256(abi.encode(_commitment));
    }

    /// @notice Hashes an array of signal slots.
    /// @param _signalSlots The signal slots to hash.
    /// @return The keccak256 hash (bytes32(0) if empty).
    function hashSignalSlots(bytes32[] memory _signalSlots) public pure returns (bytes32) {
        if (_signalSlots.length == 0) return bytes32(0);
        return keccak256(abi.encode(_signalSlots));
    }

    // ---------------------------------------------------------------
    // External View Functions
    // ---------------------------------------------------------------

    /// @inheritdoc IRealTimeInbox
    function getLastFinalizedBlockHash() external view returns (bytes32) {
        return lastFinalizedBlockHash;
    }

    /// @inheritdoc IRealTimeInbox
    function getConfig() external view returns (Config memory config_) {
        config_ = Config({
            proofVerifier: address(_proofVerifier),
            signalService: address(_signalService),
            basefeeSharingPctg: _basefeeSharingPctg,
            forcedInclusionDelay: _forcedInclusionDelay,
            forcedInclusionFeeInGwei: _forcedInclusionFeeInGwei,
            forcedInclusionFeeDoubleThreshold: _forcedInclusionFeeDoubleThreshold
        });
    }

    /// @inheritdoc IForcedInclusionStore
    function getCurrentForcedInclusionFee() external view returns (uint64 feeInGwei_) {
        return _forcedInclusionStorage.getCurrentForcedInclusionFee(
            _forcedInclusionFeeInGwei, _forcedInclusionFeeDoubleThreshold
        );
    }

    /// @inheritdoc IForcedInclusionStore
    function getForcedInclusions(
        uint48 _start,
        uint48 _maxCount
    )
        external
        view
        returns (IForcedInclusionStore.ForcedInclusion[] memory inclusions_)
    {
        return _forcedInclusionStorage.getForcedInclusions(_start, _maxCount);
    }

    /// @inheritdoc IForcedInclusionStore
    function getForcedInclusionState() external view returns (uint48 head_, uint48 tail_) {
        return _forcedInclusionStorage.getForcedInclusionState();
    }

    // ---------------------------------------------------------------
    // Internal Functions
    // ---------------------------------------------------------------

    /// @dev Decodes input, validates it, and builds the transient proposal.
    /// @param _data The encoded ProposeInput.
    /// @return proposalHash_ The hash of the proposal.
    /// @return proposal_ The built proposal struct.
    /// @return signalSlots_ The raw signal slots from the input.
    function _buildProposal(bytes calldata _data)
        internal
        returns (bytes32 proposalHash_, Proposal memory proposal_, bytes32[] memory signalSlots_)
    {
        ProposeInput memory input = decodeProposeInput(_data);
        signalSlots_ = input.signalSlots;

        // Validate anchor block - blockhash returns 0 for blocks older than 256
        bytes32 anchorHash = blockhash(input.maxAnchorBlockNumber);
        require(anchorHash != bytes32(0), MaxAnchorBlockTooOld());

        // Verify signal slots and compute hash
        bytes32 signalSlotsHash = _verifySignalSlots(input.signalSlots);

        // Validate blob reference
        LibBlobs.BlobSlice memory blobSlice = LibBlobs.validateBlobReference(input.blobReference);
        // Zero timestamp so it doesn't become part of the proposal hash that must be proven.
        // The driver can derive the blob timestamp from the L1 block that contains the event.
        blobSlice.timestamp = 0;

        // Build derivation sources: forced inclusions (if any) come first, the proposer's
        // own blob last. Also dequeues from the FI queue and forwards fees to msg.sender.
        IInbox.DerivationSource[] memory sources =
            _consumeForcedInclusions(input.numForcedInclusions, blobSlice);

        // Build proposal (standalone — no parent linkage)
        proposal_ = Proposal({
            maxAnchorBlockNumber: input.maxAnchorBlockNumber,
            maxAnchorBlockHash: anchorHash,
            basefeeSharingPctg: _basefeeSharingPctg,
            sources: sources,
            signalSlotsHash: signalSlotsHash
        });

        proposalHash_ = hashProposal(proposal_);
    }

    /// @dev Decodes ProposeInputV2, validates existing signals, and builds the proposal.
    ///      The anchor on L2 consumes the union of existing + required signals as fast
    ///      signals; the proof commits to this union. Only existing signals are
    ///      verified on L1 here; required signals are verified at finalize.
    /// @param _data The encoded ProposeInputV2.
    /// @return proposalHash_ The hash of the proposal.
    /// @return proposal_ The built proposal struct.
    /// @return unionSignalSlots_ Union of existing + required signal slots.
    /// @return requiredSignalsHash_ Hash of the required return signal list.
    function _buildProposalV2(bytes calldata _data)
        internal
        returns (
            bytes32 proposalHash_,
            Proposal memory proposal_,
            bytes32[] memory unionSignalSlots_,
            bytes32 requiredSignalsHash_
        )
    {
        ProposeInputV2 memory input = abi.decode(_data, (ProposeInputV2));

        // Validate anchor block - blockhash returns 0 for blocks older than 256
        bytes32 anchorHash = blockhash(input.maxAnchorBlockNumber);
        require(anchorHash != bytes32(0), MaxAnchorBlockTooOld());

        // Verify existing signals are already on L1.
        for (uint256 i; i < input.existingSignals.length; ++i) {
            require(
                _signalService.isSignalSent(input.existingSignals[i]),
                SignalSlotNotSent(input.existingSignals[i])
            );
        }

        // Build the union of signal slots. The L2 block's anchor consumes both
        // lists as fast signals, and the ZK proof commits to the union hash.
        unionSignalSlots_ = _concatSignals(input.existingSignals, input.requiredReturnSignals);
        requiredSignalsHash_ = hashSignalSlots(input.requiredReturnSignals);
        bytes32 signalSlotsHash = hashSignalSlots(unionSignalSlots_);

        // Validate blob reference
        LibBlobs.BlobSlice memory blobSlice = LibBlobs.validateBlobReference(input.blobReference);
        blobSlice.timestamp = 0;

        // Build derivation sources: forced inclusions (if any) come first, the proposer's
        // own blob last. Also dequeues from the FI queue and forwards fees to msg.sender.
        IInbox.DerivationSource[] memory sources =
            _consumeForcedInclusions(input.numForcedInclusions, blobSlice);

        proposal_ = Proposal({
            maxAnchorBlockNumber: input.maxAnchorBlockNumber,
            maxAnchorBlockHash: anchorHash,
            basefeeSharingPctg: _basefeeSharingPctg,
            sources: sources,
            signalSlotsHash: signalSlotsHash
        });

        proposalHash_ = hashProposal(proposal_);
    }

    /// @dev Dequeues forced inclusions from the queue, builds the combined sources array
    ///      (forced inclusions first, proposer's own blob last), and forwards accumulated fees
    ///      to `msg.sender`. Reverts if the proposer fails to consume an overdue inclusion.
    /// @param _numForcedInclusionsRequested The number of forced inclusions requested by the proposer.
    /// @param _proposerBlobSlice The proposer's own blob slice (already validated by caller).
    /// @return sources_ Sources array of length `toProcess + 1`. Forced inclusions occupy
    ///         indices `[0, toProcess)`; the proposer's blob is at index `toProcess`.
    function _consumeForcedInclusions(
        uint16 _numForcedInclusionsRequested,
        LibBlobs.BlobSlice memory _proposerBlobSlice
    )
        private
        returns (IInbox.DerivationSource[] memory sources_)
    {
        unchecked {
            LibForcedInclusion.Storage storage $ = _forcedInclusionStorage;

            (uint48 head, uint48 tail) = ($.head, $.tail);
            uint256 available = uint256(tail) - uint256(head);

            uint256 toProcess = uint256(_numForcedInclusionsRequested) > available
                ? available
                : uint256(_numForcedInclusionsRequested);

            uint48 headAfter = head + uint48(toProcess);

            // If unconsumed inclusions remain and the next one is past its delay, the proposer
            // is censoring it — revert.
            if (available > toProcess) {
                require(
                    !$.isOldestForcedInclusionDue(headAfter, tail, _forcedInclusionDelay),
                    UnprocessedForcedInclusionIsDue()
                );
            }

            sources_ = new IInbox.DerivationSource[](toProcess + 1);

            uint256 totalFees;
            for (uint256 i; i < toProcess; ++i) {
                IForcedInclusionStore.ForcedInclusion storage inclusion = $.queue[head + i];
                sources_[i] = IInbox.DerivationSource(true, inclusion.blobSlice);
                totalFees += uint256(inclusion.feeInGwei);
            }

            // Proposer's own blob always last (drives the canonical L2 progression).
            sources_[toProcess] = IInbox.DerivationSource(false, _proposerBlobSlice);

            if (toProcess > 0) {
                $.head = headAfter;
            }

            if (totalFees > 0) {
                msg.sender.sendEtherAndVerify(totalFees * 1 gwei);
            }
        }
    }

    /// @dev Concatenates two signal slot arrays. Existing signals come first,
    ///      followed by required return signals. Order is part of the proof's commitment.
    function _concatSignals(
        bytes32[] memory _existing,
        bytes32[] memory _required
    )
        internal
        pure
        returns (bytes32[] memory result_)
    {
        result_ = new bytes32[](_existing.length + _required.length);
        for (uint256 i; i < _existing.length; ++i) {
            result_[i] = _existing[i];
        }
        for (uint256 i; i < _required.length; ++i) {
            result_[_existing.length + i] = _required[i];
        }
    }

    /// @dev Verifies the proof, saves checkpoint, and updates chain head.
    /// @param _proposalHash The proposal hash.
    /// @param _lastFinalizedBlockHash The block hash of the last finalized L2 block.
    /// @param _checkpoint The checkpoint to save.
    /// @param _proof The ZK proof bytes.
    function _verifyAndFinalize(
        bytes32 _proposalHash,
        bytes32 _lastFinalizedBlockHash,
        ICheckpointStore.Checkpoint calldata _checkpoint,
        bytes calldata _proof
    )
        internal
    {
        // Build commitment and hash it
        bytes32 commitmentHash = hashCommitment(
            Commitment({
                proposalHash: _proposalHash,
                lastFinalizedBlockHash: _lastFinalizedBlockHash,
                checkpoint: _checkpoint
            })
        );

        // Verify proof via SurgeVerifier
        _proofVerifier.verifyProof(true, commitmentHash, _proof);

        // Save checkpoint to signal service
        _signalService.saveCheckpoint(_checkpoint);

        // Update chain head to the new finalized block hash
        lastFinalizedBlockHash = _checkpoint.blockHash;
    }

    /// @dev Verifies signal slots exist on L1 and returns their hash.
    /// @param _signalSlots The signal slots to verify.
    /// @return signalSlotsHash_ The keccak256 hash of slots (bytes32(0) if empty).
    function _verifySignalSlots(bytes32[] memory _signalSlots)
        internal
        view
        returns (bytes32 signalSlotsHash_)
    {
        if (_signalSlots.length == 0) return bytes32(0);

        for (uint256 i; i < _signalSlots.length; ++i) {
            require(
                _signalService.isSignalSent(_signalSlots[i]), SignalSlotNotSent(_signalSlots[i])
            );
        }
        signalSlotsHash_ = hashSignalSlots(_signalSlots);
    }

    // ---------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------

    error AlreadyActivated();
    error MaxAnchorBlockTooOld();
    error InvalidGenesisBlockHash();
    error NotActivated();
    error SignalSlotNotSent(bytes32 slot);
    error PendingProposalAlreadyExists();
    error NoPendingProposal();
    error RequiredSignalsMismatch();
    error RequiredSignalNotSent(bytes32 slot);
    error UnprocessedForcedInclusionIsDue();
}
