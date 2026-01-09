// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package shasta

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// RollbackInboxMetaData contains all meta data concerning the RollbackInbox contract.
var RollbackInboxMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"activate\",\"inputs\":[{\"name\":\"_lastPacayaBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"activationTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decodeProposeInput\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"input_\",\"type\":\"tuple\",\"internalType\":\"structIInbox.ProposeInput\",\"components\":[{\"name\":\"deadline\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blobReference\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobReference\",\"components\":[{\"name\":\"blobStartIndex\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"numBlobs\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]},{\"name\":\"numForcedInclusions\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"decodeProveInput\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"input_\",\"type\":\"tuple\",\"internalType\":\"structIInbox.ProveInput\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Commitment\",\"components\":[{\"name\":\"firstProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"firstProposalParentBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lastProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actualProver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"endBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"transitions\",\"type\":\"tuple[]\",\"internalType\":\"structIInbox.Transition[]\",\"components\":[{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"designatedProver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"forceCheckpointSync\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeProposeInput\",\"inputs\":[{\"name\":\"_input\",\"type\":\"tuple\",\"internalType\":\"structIInbox.ProposeInput\",\"components\":[{\"name\":\"deadline\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blobReference\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobReference\",\"components\":[{\"name\":\"blobStartIndex\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"numBlobs\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]},{\"name\":\"numForcedInclusions\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[{\"name\":\"encoded_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeProveInput\",\"inputs\":[{\"name\":\"_input\",\"type\":\"tuple\",\"internalType\":\"structIInbox.ProveInput\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Commitment\",\"components\":[{\"name\":\"firstProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"firstProposalParentBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lastProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actualProver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"endBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"transitions\",\"type\":\"tuple[]\",\"internalType\":\"structIInbox.Transition[]\",\"components\":[{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"designatedProver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"forceCheckpointSync\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[{\"name\":\"encoded_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"config_\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Config\",\"components\":[{\"name\":\"proofVerifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposerChecker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proverWhitelist\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signalService\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provingWindow\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"maxProofSubmissionDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"ringBufferSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"minForcedInclusionCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"forcedInclusionDelay\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"forcedInclusionFeeInGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"forcedInclusionFeeDoubleThreshold\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"minCheckpointDelay\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"permissionlessInclusionMultiplier\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCoreState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIInbox.CoreState\",\"components\":[{\"name\":\"nextProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastProposalBlockId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastFinalizedProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastFinalizedTimestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastCheckpointTimestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastFinalizedBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentForcedInclusionFee\",\"inputs\":[],\"outputs\":[{\"name\":\"feeInGwei_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getForcedInclusionState\",\"inputs\":[],\"outputs\":[{\"name\":\"head_\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"tail_\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getForcedInclusions\",\"inputs\":[{\"name\":\"_start\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"_maxCount\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[{\"name\":\"inclusions_\",\"type\":\"tuple[]\",\"internalType\":\"structIForcedInclusionStore.ForcedInclusion[]\",\"components\":[{\"name\":\"feeInGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProposalHash\",\"inputs\":[{\"name\":\"_proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashBondInstruction\",\"inputs\":[{\"name\":\"_bondInstruction\",\"type\":\"tuple\",\"internalType\":\"structLibBonds.BondInstruction\",\"components\":[{\"name\":\"proposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"bondType\",\"type\":\"uint8\",\"internalType\":\"enumLibBonds.BondType\"},{\"name\":\"payer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"hashCommitment\",\"inputs\":[{\"name\":\"_commitment\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Commitment\",\"components\":[{\"name\":\"firstProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"firstProposalParentBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lastProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actualProver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"endBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"transitions\",\"type\":\"tuple[]\",\"internalType\":\"structIInbox.Transition[]\",\"components\":[{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"designatedProver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"hashProposal\",\"inputs\":[{\"name\":\"_proposal\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Proposal\",\"components\":[{\"name\":\"id\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endOfSubmissionWindowTimestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"parentProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"originBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"originBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"sources\",\"type\":\"tuple[]\",\"internalType\":\"structIInbox.DerivationSource[]\",\"components\":[{\"name\":\"isForcedInclusion\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"impl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inLimpMode\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inNonReentrant\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"init\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxFinalizationDelayBeforeRollback\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"propose\",\"inputs\":[{\"name\":\"_lookahead\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proposeAndProve\",\"inputs\":[{\"name\":\"_lookahead\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_proposeData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_proveData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"prove\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rollback\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"saveForcedInclusion\",\"inputs\":[{\"name\":\"_blobReference\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobReference\",\"components\":[{\"name\":\"blobStartIndex\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"numBlobs\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setLimpMode\",\"inputs\":[{\"name\":\"_val\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AdminChanged\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BeaconUpgraded\",\"inputs\":[{\"name\":\"beacon\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BondInstructionCreated\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint48\",\"indexed\":true,\"internalType\":\"uint48\"},{\"name\":\"bondInstruction\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structLibBonds.BondInstruction\",\"components\":[{\"name\":\"proposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"bondType\",\"type\":\"uint8\",\"internalType\":\"enumLibBonds.BondType\"},{\"name\":\"payer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ForcedInclusionSaved\",\"inputs\":[{\"name\":\"forcedInclusion\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIForcedInclusionStore.ForcedInclusion\",\"components\":[{\"name\":\"feeInGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboxActivated\",\"inputs\":[{\"name\":\"lastPacayaBlockHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LimpModeSet\",\"inputs\":[{\"name\":\"enabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Proposed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint48\",\"indexed\":true,\"internalType\":\"uint48\"},{\"name\":\"proposer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"parentProposalHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"endOfSubmissionWindowTimestamp\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"sources\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structIInbox.DerivationSource[]\",\"components\":[{\"name\":\"isForcedInclusion\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Proved\",\"inputs\":[{\"name\":\"firstProposalId\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"firstNewProposalId\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"lastProposalId\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"actualProver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"checkpointSynced\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Rollbacked\",\"inputs\":[{\"name\":\"firstProposalId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lastProposalId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ACCESS_DENIED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ActivationRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BlobNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotProposeInCurrentBlock\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckpointDelayHasPassed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DeadlineExceeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ETH_TRANSFER_FAILED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyBatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FUNC_NOT_IMPLEMENTED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FirstProposalIdTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"INVALID_PAUSE_STATUS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectProposalCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LastProposalAlreadyFinalized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LastProposalHashMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LastProposalIdTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthExceedsUint16\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoBlobs\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughCapacity\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParentBlockHashMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProverNotWhitelisted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"REENTRANT_CALL\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Surge_CannotProposeDirectlyInLimpMode\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Surge_CannotProveDirectlyInLimpMode\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Surge_HeadMustBeFinalizedInLimpMode\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Surge_NoProposalsToRollback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Surge_RollbackNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnprocessedForcedInclusionIsDue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_ADDRESS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_VALUE\",\"inputs\":[]}]",
}

// RollbackInboxABI is the input ABI used to generate the binding from.
// Deprecated: Use RollbackInboxMetaData.ABI instead.
var RollbackInboxABI = RollbackInboxMetaData.ABI

// RollbackInbox is an auto generated Go binding around an Ethereum contract.
type RollbackInbox struct {
	RollbackInboxCaller     // Read-only binding to the contract
	RollbackInboxTransactor // Write-only binding to the contract
	RollbackInboxFilterer   // Log filterer for contract events
}

// RollbackInboxCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollbackInboxCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollbackInboxTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollbackInboxTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollbackInboxFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollbackInboxFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollbackInboxSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollbackInboxSession struct {
	Contract     *RollbackInbox    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollbackInboxCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollbackInboxCallerSession struct {
	Contract *RollbackInboxCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// RollbackInboxTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollbackInboxTransactorSession struct {
	Contract     *RollbackInboxTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RollbackInboxRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollbackInboxRaw struct {
	Contract *RollbackInbox // Generic contract binding to access the raw methods on
}

// RollbackInboxCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollbackInboxCallerRaw struct {
	Contract *RollbackInboxCaller // Generic read-only contract binding to access the raw methods on
}

// RollbackInboxTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollbackInboxTransactorRaw struct {
	Contract *RollbackInboxTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollbackInbox creates a new instance of RollbackInbox, bound to a specific deployed contract.
func NewRollbackInbox(address common.Address, backend bind.ContractBackend) (*RollbackInbox, error) {
	contract, err := bindRollbackInbox(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RollbackInbox{RollbackInboxCaller: RollbackInboxCaller{contract: contract}, RollbackInboxTransactor: RollbackInboxTransactor{contract: contract}, RollbackInboxFilterer: RollbackInboxFilterer{contract: contract}}, nil
}

// NewRollbackInboxCaller creates a new read-only instance of RollbackInbox, bound to a specific deployed contract.
func NewRollbackInboxCaller(address common.Address, caller bind.ContractCaller) (*RollbackInboxCaller, error) {
	contract, err := bindRollbackInbox(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxCaller{contract: contract}, nil
}

// NewRollbackInboxTransactor creates a new write-only instance of RollbackInbox, bound to a specific deployed contract.
func NewRollbackInboxTransactor(address common.Address, transactor bind.ContractTransactor) (*RollbackInboxTransactor, error) {
	contract, err := bindRollbackInbox(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxTransactor{contract: contract}, nil
}

// NewRollbackInboxFilterer creates a new log filterer instance of RollbackInbox, bound to a specific deployed contract.
func NewRollbackInboxFilterer(address common.Address, filterer bind.ContractFilterer) (*RollbackInboxFilterer, error) {
	contract, err := bindRollbackInbox(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxFilterer{contract: contract}, nil
}

// bindRollbackInbox binds a generic wrapper to an already deployed contract.
func bindRollbackInbox(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RollbackInboxMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollbackInbox *RollbackInboxRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollbackInbox.Contract.RollbackInboxCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollbackInbox *RollbackInboxRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollbackInbox.Contract.RollbackInboxTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollbackInbox *RollbackInboxRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollbackInbox.Contract.RollbackInboxTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollbackInbox *RollbackInboxCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollbackInbox.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollbackInbox *RollbackInboxTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollbackInbox.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollbackInbox *RollbackInboxTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollbackInbox.Contract.contract.Transact(opts, method, params...)
}

// ActivationTimestamp is a free data retrieval call binding the contract method 0x0423c7de.
//
// Solidity: function activationTimestamp() view returns(uint48)
func (_RollbackInbox *RollbackInboxCaller) ActivationTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "activationTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActivationTimestamp is a free data retrieval call binding the contract method 0x0423c7de.
//
// Solidity: function activationTimestamp() view returns(uint48)
func (_RollbackInbox *RollbackInboxSession) ActivationTimestamp() (*big.Int, error) {
	return _RollbackInbox.Contract.ActivationTimestamp(&_RollbackInbox.CallOpts)
}

// ActivationTimestamp is a free data retrieval call binding the contract method 0x0423c7de.
//
// Solidity: function activationTimestamp() view returns(uint48)
func (_RollbackInbox *RollbackInboxCallerSession) ActivationTimestamp() (*big.Int, error) {
	return _RollbackInbox.Contract.ActivationTimestamp(&_RollbackInbox.CallOpts)
}

// DecodeProposeInput is a free data retrieval call binding the contract method 0xafb63ad4.
//
// Solidity: function decodeProposeInput(bytes _data) pure returns((uint48,(uint16,uint16,uint24),uint8) input_)
func (_RollbackInbox *RollbackInboxCaller) DecodeProposeInput(opts *bind.CallOpts, _data []byte) (IInboxProposeInput, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "decodeProposeInput", _data)

	if err != nil {
		return *new(IInboxProposeInput), err
	}

	out0 := *abi.ConvertType(out[0], new(IInboxProposeInput)).(*IInboxProposeInput)

	return out0, err

}

// DecodeProposeInput is a free data retrieval call binding the contract method 0xafb63ad4.
//
// Solidity: function decodeProposeInput(bytes _data) pure returns((uint48,(uint16,uint16,uint24),uint8) input_)
func (_RollbackInbox *RollbackInboxSession) DecodeProposeInput(_data []byte) (IInboxProposeInput, error) {
	return _RollbackInbox.Contract.DecodeProposeInput(&_RollbackInbox.CallOpts, _data)
}

// DecodeProposeInput is a free data retrieval call binding the contract method 0xafb63ad4.
//
// Solidity: function decodeProposeInput(bytes _data) pure returns((uint48,(uint16,uint16,uint24),uint8) input_)
func (_RollbackInbox *RollbackInboxCallerSession) DecodeProposeInput(_data []byte) (IInboxProposeInput, error) {
	return _RollbackInbox.Contract.DecodeProposeInput(&_RollbackInbox.CallOpts, _data)
}

// DecodeProveInput is a free data retrieval call binding the contract method 0xedbacd44.
//
// Solidity: function decodeProveInput(bytes _data) pure returns(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,address,uint48,bytes32)[]),bool) input_)
func (_RollbackInbox *RollbackInboxCaller) DecodeProveInput(opts *bind.CallOpts, _data []byte) (IInboxProveInput, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "decodeProveInput", _data)

	if err != nil {
		return *new(IInboxProveInput), err
	}

	out0 := *abi.ConvertType(out[0], new(IInboxProveInput)).(*IInboxProveInput)

	return out0, err

}

// DecodeProveInput is a free data retrieval call binding the contract method 0xedbacd44.
//
// Solidity: function decodeProveInput(bytes _data) pure returns(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,address,uint48,bytes32)[]),bool) input_)
func (_RollbackInbox *RollbackInboxSession) DecodeProveInput(_data []byte) (IInboxProveInput, error) {
	return _RollbackInbox.Contract.DecodeProveInput(&_RollbackInbox.CallOpts, _data)
}

// DecodeProveInput is a free data retrieval call binding the contract method 0xedbacd44.
//
// Solidity: function decodeProveInput(bytes _data) pure returns(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,address,uint48,bytes32)[]),bool) input_)
func (_RollbackInbox *RollbackInboxCallerSession) DecodeProveInput(_data []byte) (IInboxProveInput, error) {
	return _RollbackInbox.Contract.DecodeProveInput(&_RollbackInbox.CallOpts, _data)
}

// EncodeProposeInput is a free data retrieval call binding the contract method 0x2f1969b0.
//
// Solidity: function encodeProposeInput((uint48,(uint16,uint16,uint24),uint8) _input) pure returns(bytes encoded_)
func (_RollbackInbox *RollbackInboxCaller) EncodeProposeInput(opts *bind.CallOpts, _input IInboxProposeInput) ([]byte, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "encodeProposeInput", _input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeProposeInput is a free data retrieval call binding the contract method 0x2f1969b0.
//
// Solidity: function encodeProposeInput((uint48,(uint16,uint16,uint24),uint8) _input) pure returns(bytes encoded_)
func (_RollbackInbox *RollbackInboxSession) EncodeProposeInput(_input IInboxProposeInput) ([]byte, error) {
	return _RollbackInbox.Contract.EncodeProposeInput(&_RollbackInbox.CallOpts, _input)
}

// EncodeProposeInput is a free data retrieval call binding the contract method 0x2f1969b0.
//
// Solidity: function encodeProposeInput((uint48,(uint16,uint16,uint24),uint8) _input) pure returns(bytes encoded_)
func (_RollbackInbox *RollbackInboxCallerSession) EncodeProposeInput(_input IInboxProposeInput) ([]byte, error) {
	return _RollbackInbox.Contract.EncodeProposeInput(&_RollbackInbox.CallOpts, _input)
}

// EncodeProveInput is a free data retrieval call binding the contract method 0xc3d3e2f4.
//
// Solidity: function encodeProveInput(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,address,uint48,bytes32)[]),bool) _input) pure returns(bytes encoded_)
func (_RollbackInbox *RollbackInboxCaller) EncodeProveInput(opts *bind.CallOpts, _input IInboxProveInput) ([]byte, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "encodeProveInput", _input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeProveInput is a free data retrieval call binding the contract method 0xc3d3e2f4.
//
// Solidity: function encodeProveInput(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,address,uint48,bytes32)[]),bool) _input) pure returns(bytes encoded_)
func (_RollbackInbox *RollbackInboxSession) EncodeProveInput(_input IInboxProveInput) ([]byte, error) {
	return _RollbackInbox.Contract.EncodeProveInput(&_RollbackInbox.CallOpts, _input)
}

// EncodeProveInput is a free data retrieval call binding the contract method 0xc3d3e2f4.
//
// Solidity: function encodeProveInput(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,address,uint48,bytes32)[]),bool) _input) pure returns(bytes encoded_)
func (_RollbackInbox *RollbackInboxCallerSession) EncodeProveInput(_input IInboxProveInput) ([]byte, error) {
	return _RollbackInbox.Contract.EncodeProveInput(&_RollbackInbox.CallOpts, _input)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,address,address,uint48,uint48,uint256,uint8,uint256,uint16,uint64,uint64,uint16,uint8) config_)
func (_RollbackInbox *RollbackInboxCaller) GetConfig(opts *bind.CallOpts) (IInboxConfig, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(IInboxConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IInboxConfig)).(*IInboxConfig)

	return out0, err

}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,address,address,uint48,uint48,uint256,uint8,uint256,uint16,uint64,uint64,uint16,uint8) config_)
func (_RollbackInbox *RollbackInboxSession) GetConfig() (IInboxConfig, error) {
	return _RollbackInbox.Contract.GetConfig(&_RollbackInbox.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,address,address,uint48,uint48,uint256,uint8,uint256,uint16,uint64,uint64,uint16,uint8) config_)
func (_RollbackInbox *RollbackInboxCallerSession) GetConfig() (IInboxConfig, error) {
	return _RollbackInbox.Contract.GetConfig(&_RollbackInbox.CallOpts)
}

// GetCoreState is a free data retrieval call binding the contract method 0x6aa6a01a.
//
// Solidity: function getCoreState() view returns((uint48,uint48,uint48,uint48,uint48,bytes32))
func (_RollbackInbox *RollbackInboxCaller) GetCoreState(opts *bind.CallOpts) (IInboxCoreState, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "getCoreState")

	if err != nil {
		return *new(IInboxCoreState), err
	}

	out0 := *abi.ConvertType(out[0], new(IInboxCoreState)).(*IInboxCoreState)

	return out0, err

}

// GetCoreState is a free data retrieval call binding the contract method 0x6aa6a01a.
//
// Solidity: function getCoreState() view returns((uint48,uint48,uint48,uint48,uint48,bytes32))
func (_RollbackInbox *RollbackInboxSession) GetCoreState() (IInboxCoreState, error) {
	return _RollbackInbox.Contract.GetCoreState(&_RollbackInbox.CallOpts)
}

// GetCoreState is a free data retrieval call binding the contract method 0x6aa6a01a.
//
// Solidity: function getCoreState() view returns((uint48,uint48,uint48,uint48,uint48,bytes32))
func (_RollbackInbox *RollbackInboxCallerSession) GetCoreState() (IInboxCoreState, error) {
	return _RollbackInbox.Contract.GetCoreState(&_RollbackInbox.CallOpts)
}

// GetCurrentForcedInclusionFee is a free data retrieval call binding the contract method 0xe3053335.
//
// Solidity: function getCurrentForcedInclusionFee() view returns(uint64 feeInGwei_)
func (_RollbackInbox *RollbackInboxCaller) GetCurrentForcedInclusionFee(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "getCurrentForcedInclusionFee")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetCurrentForcedInclusionFee is a free data retrieval call binding the contract method 0xe3053335.
//
// Solidity: function getCurrentForcedInclusionFee() view returns(uint64 feeInGwei_)
func (_RollbackInbox *RollbackInboxSession) GetCurrentForcedInclusionFee() (uint64, error) {
	return _RollbackInbox.Contract.GetCurrentForcedInclusionFee(&_RollbackInbox.CallOpts)
}

// GetCurrentForcedInclusionFee is a free data retrieval call binding the contract method 0xe3053335.
//
// Solidity: function getCurrentForcedInclusionFee() view returns(uint64 feeInGwei_)
func (_RollbackInbox *RollbackInboxCallerSession) GetCurrentForcedInclusionFee() (uint64, error) {
	return _RollbackInbox.Contract.GetCurrentForcedInclusionFee(&_RollbackInbox.CallOpts)
}

// GetForcedInclusionState is a free data retrieval call binding the contract method 0x5ccc1718.
//
// Solidity: function getForcedInclusionState() view returns(uint48 head_, uint48 tail_)
func (_RollbackInbox *RollbackInboxCaller) GetForcedInclusionState(opts *bind.CallOpts) (struct {
	Head *big.Int
	Tail *big.Int
}, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "getForcedInclusionState")

	outstruct := new(struct {
		Head *big.Int
		Tail *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Head = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Tail = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetForcedInclusionState is a free data retrieval call binding the contract method 0x5ccc1718.
//
// Solidity: function getForcedInclusionState() view returns(uint48 head_, uint48 tail_)
func (_RollbackInbox *RollbackInboxSession) GetForcedInclusionState() (struct {
	Head *big.Int
	Tail *big.Int
}, error) {
	return _RollbackInbox.Contract.GetForcedInclusionState(&_RollbackInbox.CallOpts)
}

// GetForcedInclusionState is a free data retrieval call binding the contract method 0x5ccc1718.
//
// Solidity: function getForcedInclusionState() view returns(uint48 head_, uint48 tail_)
func (_RollbackInbox *RollbackInboxCallerSession) GetForcedInclusionState() (struct {
	Head *big.Int
	Tail *big.Int
}, error) {
	return _RollbackInbox.Contract.GetForcedInclusionState(&_RollbackInbox.CallOpts)
}

// GetForcedInclusions is a free data retrieval call binding the contract method 0x40df9866.
//
// Solidity: function getForcedInclusions(uint48 _start, uint48 _maxCount) view returns((uint64,(bytes32[],uint24,uint48))[] inclusions_)
func (_RollbackInbox *RollbackInboxCaller) GetForcedInclusions(opts *bind.CallOpts, _start *big.Int, _maxCount *big.Int) ([]IForcedInclusionStoreForcedInclusion, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "getForcedInclusions", _start, _maxCount)

	if err != nil {
		return *new([]IForcedInclusionStoreForcedInclusion), err
	}

	out0 := *abi.ConvertType(out[0], new([]IForcedInclusionStoreForcedInclusion)).(*[]IForcedInclusionStoreForcedInclusion)

	return out0, err

}

// GetForcedInclusions is a free data retrieval call binding the contract method 0x40df9866.
//
// Solidity: function getForcedInclusions(uint48 _start, uint48 _maxCount) view returns((uint64,(bytes32[],uint24,uint48))[] inclusions_)
func (_RollbackInbox *RollbackInboxSession) GetForcedInclusions(_start *big.Int, _maxCount *big.Int) ([]IForcedInclusionStoreForcedInclusion, error) {
	return _RollbackInbox.Contract.GetForcedInclusions(&_RollbackInbox.CallOpts, _start, _maxCount)
}

// GetForcedInclusions is a free data retrieval call binding the contract method 0x40df9866.
//
// Solidity: function getForcedInclusions(uint48 _start, uint48 _maxCount) view returns((uint64,(bytes32[],uint24,uint48))[] inclusions_)
func (_RollbackInbox *RollbackInboxCallerSession) GetForcedInclusions(_start *big.Int, _maxCount *big.Int) ([]IForcedInclusionStoreForcedInclusion, error) {
	return _RollbackInbox.Contract.GetForcedInclusions(&_RollbackInbox.CallOpts, _start, _maxCount)
}

// GetProposalHash is a free data retrieval call binding the contract method 0xa834725a.
//
// Solidity: function getProposalHash(uint256 _proposalId) view returns(bytes32)
func (_RollbackInbox *RollbackInboxCaller) GetProposalHash(opts *bind.CallOpts, _proposalId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "getProposalHash", _proposalId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetProposalHash is a free data retrieval call binding the contract method 0xa834725a.
//
// Solidity: function getProposalHash(uint256 _proposalId) view returns(bytes32)
func (_RollbackInbox *RollbackInboxSession) GetProposalHash(_proposalId *big.Int) ([32]byte, error) {
	return _RollbackInbox.Contract.GetProposalHash(&_RollbackInbox.CallOpts, _proposalId)
}

// GetProposalHash is a free data retrieval call binding the contract method 0xa834725a.
//
// Solidity: function getProposalHash(uint256 _proposalId) view returns(bytes32)
func (_RollbackInbox *RollbackInboxCallerSession) GetProposalHash(_proposalId *big.Int) ([32]byte, error) {
	return _RollbackInbox.Contract.GetProposalHash(&_RollbackInbox.CallOpts, _proposalId)
}

// HashBondInstruction is a free data retrieval call binding the contract method 0x5a213615.
//
// Solidity: function hashBondInstruction((uint48,uint8,address,address) _bondInstruction) pure returns(bytes32)
func (_RollbackInbox *RollbackInboxCaller) HashBondInstruction(opts *bind.CallOpts, _bondInstruction LibBondsBondInstruction) ([32]byte, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "hashBondInstruction", _bondInstruction)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashBondInstruction is a free data retrieval call binding the contract method 0x5a213615.
//
// Solidity: function hashBondInstruction((uint48,uint8,address,address) _bondInstruction) pure returns(bytes32)
func (_RollbackInbox *RollbackInboxSession) HashBondInstruction(_bondInstruction LibBondsBondInstruction) ([32]byte, error) {
	return _RollbackInbox.Contract.HashBondInstruction(&_RollbackInbox.CallOpts, _bondInstruction)
}

// HashBondInstruction is a free data retrieval call binding the contract method 0x5a213615.
//
// Solidity: function hashBondInstruction((uint48,uint8,address,address) _bondInstruction) pure returns(bytes32)
func (_RollbackInbox *RollbackInboxCallerSession) HashBondInstruction(_bondInstruction LibBondsBondInstruction) ([32]byte, error) {
	return _RollbackInbox.Contract.HashBondInstruction(&_RollbackInbox.CallOpts, _bondInstruction)
}

// HashCommitment is a free data retrieval call binding the contract method 0xcbc148c3.
//
// Solidity: function hashCommitment((uint48,bytes32,bytes32,address,uint48,bytes32,(address,address,uint48,bytes32)[]) _commitment) pure returns(bytes32)
func (_RollbackInbox *RollbackInboxCaller) HashCommitment(opts *bind.CallOpts, _commitment IInboxCommitment) ([32]byte, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "hashCommitment", _commitment)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashCommitment is a free data retrieval call binding the contract method 0xcbc148c3.
//
// Solidity: function hashCommitment((uint48,bytes32,bytes32,address,uint48,bytes32,(address,address,uint48,bytes32)[]) _commitment) pure returns(bytes32)
func (_RollbackInbox *RollbackInboxSession) HashCommitment(_commitment IInboxCommitment) ([32]byte, error) {
	return _RollbackInbox.Contract.HashCommitment(&_RollbackInbox.CallOpts, _commitment)
}

// HashCommitment is a free data retrieval call binding the contract method 0xcbc148c3.
//
// Solidity: function hashCommitment((uint48,bytes32,bytes32,address,uint48,bytes32,(address,address,uint48,bytes32)[]) _commitment) pure returns(bytes32)
func (_RollbackInbox *RollbackInboxCallerSession) HashCommitment(_commitment IInboxCommitment) ([32]byte, error) {
	return _RollbackInbox.Contract.HashCommitment(&_RollbackInbox.CallOpts, _commitment)
}

// HashProposal is a free data retrieval call binding the contract method 0xb28e824e.
//
// Solidity: function hashProposal((uint48,uint48,uint48,address,bytes32,uint48,bytes32,uint8,(bool,(bytes32[],uint24,uint48))[]) _proposal) pure returns(bytes32)
func (_RollbackInbox *RollbackInboxCaller) HashProposal(opts *bind.CallOpts, _proposal IInboxProposal) ([32]byte, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "hashProposal", _proposal)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashProposal is a free data retrieval call binding the contract method 0xb28e824e.
//
// Solidity: function hashProposal((uint48,uint48,uint48,address,bytes32,uint48,bytes32,uint8,(bool,(bytes32[],uint24,uint48))[]) _proposal) pure returns(bytes32)
func (_RollbackInbox *RollbackInboxSession) HashProposal(_proposal IInboxProposal) ([32]byte, error) {
	return _RollbackInbox.Contract.HashProposal(&_RollbackInbox.CallOpts, _proposal)
}

// HashProposal is a free data retrieval call binding the contract method 0xb28e824e.
//
// Solidity: function hashProposal((uint48,uint48,uint48,address,bytes32,uint48,bytes32,uint8,(bool,(bytes32[],uint24,uint48))[]) _proposal) pure returns(bytes32)
func (_RollbackInbox *RollbackInboxCallerSession) HashProposal(_proposal IInboxProposal) ([32]byte, error) {
	return _RollbackInbox.Contract.HashProposal(&_RollbackInbox.CallOpts, _proposal)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_RollbackInbox *RollbackInboxCaller) Impl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "impl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_RollbackInbox *RollbackInboxSession) Impl() (common.Address, error) {
	return _RollbackInbox.Contract.Impl(&_RollbackInbox.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_RollbackInbox *RollbackInboxCallerSession) Impl() (common.Address, error) {
	return _RollbackInbox.Contract.Impl(&_RollbackInbox.CallOpts)
}

// InLimpMode is a free data retrieval call binding the contract method 0xf1d338ae.
//
// Solidity: function inLimpMode() view returns(bool)
func (_RollbackInbox *RollbackInboxCaller) InLimpMode(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "inLimpMode")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// InLimpMode is a free data retrieval call binding the contract method 0xf1d338ae.
//
// Solidity: function inLimpMode() view returns(bool)
func (_RollbackInbox *RollbackInboxSession) InLimpMode() (bool, error) {
	return _RollbackInbox.Contract.InLimpMode(&_RollbackInbox.CallOpts)
}

// InLimpMode is a free data retrieval call binding the contract method 0xf1d338ae.
//
// Solidity: function inLimpMode() view returns(bool)
func (_RollbackInbox *RollbackInboxCallerSession) InLimpMode() (bool, error) {
	return _RollbackInbox.Contract.InLimpMode(&_RollbackInbox.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_RollbackInbox *RollbackInboxCaller) InNonReentrant(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "inNonReentrant")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_RollbackInbox *RollbackInboxSession) InNonReentrant() (bool, error) {
	return _RollbackInbox.Contract.InNonReentrant(&_RollbackInbox.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_RollbackInbox *RollbackInboxCallerSession) InNonReentrant() (bool, error) {
	return _RollbackInbox.Contract.InNonReentrant(&_RollbackInbox.CallOpts)
}

// MaxFinalizationDelayBeforeRollback is a free data retrieval call binding the contract method 0x1fd7286d.
//
// Solidity: function maxFinalizationDelayBeforeRollback() view returns(uint48)
func (_RollbackInbox *RollbackInboxCaller) MaxFinalizationDelayBeforeRollback(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "maxFinalizationDelayBeforeRollback")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxFinalizationDelayBeforeRollback is a free data retrieval call binding the contract method 0x1fd7286d.
//
// Solidity: function maxFinalizationDelayBeforeRollback() view returns(uint48)
func (_RollbackInbox *RollbackInboxSession) MaxFinalizationDelayBeforeRollback() (*big.Int, error) {
	return _RollbackInbox.Contract.MaxFinalizationDelayBeforeRollback(&_RollbackInbox.CallOpts)
}

// MaxFinalizationDelayBeforeRollback is a free data retrieval call binding the contract method 0x1fd7286d.
//
// Solidity: function maxFinalizationDelayBeforeRollback() view returns(uint48)
func (_RollbackInbox *RollbackInboxCallerSession) MaxFinalizationDelayBeforeRollback() (*big.Int, error) {
	return _RollbackInbox.Contract.MaxFinalizationDelayBeforeRollback(&_RollbackInbox.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RollbackInbox *RollbackInboxCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RollbackInbox *RollbackInboxSession) Owner() (common.Address, error) {
	return _RollbackInbox.Contract.Owner(&_RollbackInbox.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RollbackInbox *RollbackInboxCallerSession) Owner() (common.Address, error) {
	return _RollbackInbox.Contract.Owner(&_RollbackInbox.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RollbackInbox *RollbackInboxCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RollbackInbox *RollbackInboxSession) Paused() (bool, error) {
	return _RollbackInbox.Contract.Paused(&_RollbackInbox.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RollbackInbox *RollbackInboxCallerSession) Paused() (bool, error) {
	return _RollbackInbox.Contract.Paused(&_RollbackInbox.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_RollbackInbox *RollbackInboxCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_RollbackInbox *RollbackInboxSession) PendingOwner() (common.Address, error) {
	return _RollbackInbox.Contract.PendingOwner(&_RollbackInbox.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_RollbackInbox *RollbackInboxCallerSession) PendingOwner() (common.Address, error) {
	return _RollbackInbox.Contract.PendingOwner(&_RollbackInbox.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RollbackInbox *RollbackInboxCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RollbackInbox *RollbackInboxSession) ProxiableUUID() ([32]byte, error) {
	return _RollbackInbox.Contract.ProxiableUUID(&_RollbackInbox.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RollbackInbox *RollbackInboxCallerSession) ProxiableUUID() ([32]byte, error) {
	return _RollbackInbox.Contract.ProxiableUUID(&_RollbackInbox.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_RollbackInbox *RollbackInboxCaller) Resolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollbackInbox.contract.Call(opts, &out, "resolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_RollbackInbox *RollbackInboxSession) Resolver() (common.Address, error) {
	return _RollbackInbox.Contract.Resolver(&_RollbackInbox.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_RollbackInbox *RollbackInboxCallerSession) Resolver() (common.Address, error) {
	return _RollbackInbox.Contract.Resolver(&_RollbackInbox.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_RollbackInbox *RollbackInboxTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_RollbackInbox *RollbackInboxSession) AcceptOwnership() (*types.Transaction, error) {
	return _RollbackInbox.Contract.AcceptOwnership(&_RollbackInbox.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_RollbackInbox *RollbackInboxTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RollbackInbox.Contract.AcceptOwnership(&_RollbackInbox.TransactOpts)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _lastPacayaBlockHash) returns()
func (_RollbackInbox *RollbackInboxTransactor) Activate(opts *bind.TransactOpts, _lastPacayaBlockHash [32]byte) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "activate", _lastPacayaBlockHash)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _lastPacayaBlockHash) returns()
func (_RollbackInbox *RollbackInboxSession) Activate(_lastPacayaBlockHash [32]byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.Activate(&_RollbackInbox.TransactOpts, _lastPacayaBlockHash)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _lastPacayaBlockHash) returns()
func (_RollbackInbox *RollbackInboxTransactorSession) Activate(_lastPacayaBlockHash [32]byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.Activate(&_RollbackInbox.TransactOpts, _lastPacayaBlockHash)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_RollbackInbox *RollbackInboxTransactor) Init(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "init", _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_RollbackInbox *RollbackInboxSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _RollbackInbox.Contract.Init(&_RollbackInbox.TransactOpts, _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_RollbackInbox *RollbackInboxTransactorSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _RollbackInbox.Contract.Init(&_RollbackInbox.TransactOpts, _owner)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RollbackInbox *RollbackInboxTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RollbackInbox *RollbackInboxSession) Pause() (*types.Transaction, error) {
	return _RollbackInbox.Contract.Pause(&_RollbackInbox.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RollbackInbox *RollbackInboxTransactorSession) Pause() (*types.Transaction, error) {
	return _RollbackInbox.Contract.Pause(&_RollbackInbox.TransactOpts)
}

// Propose is a paid mutator transaction binding the contract method 0x9791e644.
//
// Solidity: function propose(bytes _lookahead, bytes _data) returns()
func (_RollbackInbox *RollbackInboxTransactor) Propose(opts *bind.TransactOpts, _lookahead []byte, _data []byte) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "propose", _lookahead, _data)
}

// Propose is a paid mutator transaction binding the contract method 0x9791e644.
//
// Solidity: function propose(bytes _lookahead, bytes _data) returns()
func (_RollbackInbox *RollbackInboxSession) Propose(_lookahead []byte, _data []byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.Propose(&_RollbackInbox.TransactOpts, _lookahead, _data)
}

// Propose is a paid mutator transaction binding the contract method 0x9791e644.
//
// Solidity: function propose(bytes _lookahead, bytes _data) returns()
func (_RollbackInbox *RollbackInboxTransactorSession) Propose(_lookahead []byte, _data []byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.Propose(&_RollbackInbox.TransactOpts, _lookahead, _data)
}

// ProposeAndProve is a paid mutator transaction binding the contract method 0xd6545ef0.
//
// Solidity: function proposeAndProve(bytes _lookahead, bytes _proposeData, bytes _proveData, bytes _proof) returns()
func (_RollbackInbox *RollbackInboxTransactor) ProposeAndProve(opts *bind.TransactOpts, _lookahead []byte, _proposeData []byte, _proveData []byte, _proof []byte) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "proposeAndProve", _lookahead, _proposeData, _proveData, _proof)
}

// ProposeAndProve is a paid mutator transaction binding the contract method 0xd6545ef0.
//
// Solidity: function proposeAndProve(bytes _lookahead, bytes _proposeData, bytes _proveData, bytes _proof) returns()
func (_RollbackInbox *RollbackInboxSession) ProposeAndProve(_lookahead []byte, _proposeData []byte, _proveData []byte, _proof []byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.ProposeAndProve(&_RollbackInbox.TransactOpts, _lookahead, _proposeData, _proveData, _proof)
}

// ProposeAndProve is a paid mutator transaction binding the contract method 0xd6545ef0.
//
// Solidity: function proposeAndProve(bytes _lookahead, bytes _proposeData, bytes _proveData, bytes _proof) returns()
func (_RollbackInbox *RollbackInboxTransactorSession) ProposeAndProve(_lookahead []byte, _proposeData []byte, _proveData []byte, _proof []byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.ProposeAndProve(&_RollbackInbox.TransactOpts, _lookahead, _proposeData, _proveData, _proof)
}

// Prove is a paid mutator transaction binding the contract method 0xea191743.
//
// Solidity: function prove(bytes _data, bytes _proof) returns()
func (_RollbackInbox *RollbackInboxTransactor) Prove(opts *bind.TransactOpts, _data []byte, _proof []byte) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "prove", _data, _proof)
}

// Prove is a paid mutator transaction binding the contract method 0xea191743.
//
// Solidity: function prove(bytes _data, bytes _proof) returns()
func (_RollbackInbox *RollbackInboxSession) Prove(_data []byte, _proof []byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.Prove(&_RollbackInbox.TransactOpts, _data, _proof)
}

// Prove is a paid mutator transaction binding the contract method 0xea191743.
//
// Solidity: function prove(bytes _data, bytes _proof) returns()
func (_RollbackInbox *RollbackInboxTransactorSession) Prove(_data []byte, _proof []byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.Prove(&_RollbackInbox.TransactOpts, _data, _proof)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RollbackInbox *RollbackInboxTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RollbackInbox *RollbackInboxSession) RenounceOwnership() (*types.Transaction, error) {
	return _RollbackInbox.Contract.RenounceOwnership(&_RollbackInbox.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RollbackInbox *RollbackInboxTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _RollbackInbox.Contract.RenounceOwnership(&_RollbackInbox.TransactOpts)
}

// Rollback is a paid mutator transaction binding the contract method 0x9afd9d78.
//
// Solidity: function rollback() returns()
func (_RollbackInbox *RollbackInboxTransactor) Rollback(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "rollback")
}

// Rollback is a paid mutator transaction binding the contract method 0x9afd9d78.
//
// Solidity: function rollback() returns()
func (_RollbackInbox *RollbackInboxSession) Rollback() (*types.Transaction, error) {
	return _RollbackInbox.Contract.Rollback(&_RollbackInbox.TransactOpts)
}

// Rollback is a paid mutator transaction binding the contract method 0x9afd9d78.
//
// Solidity: function rollback() returns()
func (_RollbackInbox *RollbackInboxTransactorSession) Rollback() (*types.Transaction, error) {
	return _RollbackInbox.Contract.Rollback(&_RollbackInbox.TransactOpts)
}

// SaveForcedInclusion is a paid mutator transaction binding the contract method 0xdf596d9e.
//
// Solidity: function saveForcedInclusion((uint16,uint16,uint24) _blobReference) payable returns()
func (_RollbackInbox *RollbackInboxTransactor) SaveForcedInclusion(opts *bind.TransactOpts, _blobReference LibBlobsBlobReference) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "saveForcedInclusion", _blobReference)
}

// SaveForcedInclusion is a paid mutator transaction binding the contract method 0xdf596d9e.
//
// Solidity: function saveForcedInclusion((uint16,uint16,uint24) _blobReference) payable returns()
func (_RollbackInbox *RollbackInboxSession) SaveForcedInclusion(_blobReference LibBlobsBlobReference) (*types.Transaction, error) {
	return _RollbackInbox.Contract.SaveForcedInclusion(&_RollbackInbox.TransactOpts, _blobReference)
}

// SaveForcedInclusion is a paid mutator transaction binding the contract method 0xdf596d9e.
//
// Solidity: function saveForcedInclusion((uint16,uint16,uint24) _blobReference) payable returns()
func (_RollbackInbox *RollbackInboxTransactorSession) SaveForcedInclusion(_blobReference LibBlobsBlobReference) (*types.Transaction, error) {
	return _RollbackInbox.Contract.SaveForcedInclusion(&_RollbackInbox.TransactOpts, _blobReference)
}

// SetLimpMode is a paid mutator transaction binding the contract method 0x3fb23ebf.
//
// Solidity: function setLimpMode(bool _val) returns()
func (_RollbackInbox *RollbackInboxTransactor) SetLimpMode(opts *bind.TransactOpts, _val bool) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "setLimpMode", _val)
}

// SetLimpMode is a paid mutator transaction binding the contract method 0x3fb23ebf.
//
// Solidity: function setLimpMode(bool _val) returns()
func (_RollbackInbox *RollbackInboxSession) SetLimpMode(_val bool) (*types.Transaction, error) {
	return _RollbackInbox.Contract.SetLimpMode(&_RollbackInbox.TransactOpts, _val)
}

// SetLimpMode is a paid mutator transaction binding the contract method 0x3fb23ebf.
//
// Solidity: function setLimpMode(bool _val) returns()
func (_RollbackInbox *RollbackInboxTransactorSession) SetLimpMode(_val bool) (*types.Transaction, error) {
	return _RollbackInbox.Contract.SetLimpMode(&_RollbackInbox.TransactOpts, _val)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RollbackInbox *RollbackInboxTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RollbackInbox *RollbackInboxSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RollbackInbox.Contract.TransferOwnership(&_RollbackInbox.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RollbackInbox *RollbackInboxTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RollbackInbox.Contract.TransferOwnership(&_RollbackInbox.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RollbackInbox *RollbackInboxTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RollbackInbox *RollbackInboxSession) Unpause() (*types.Transaction, error) {
	return _RollbackInbox.Contract.Unpause(&_RollbackInbox.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RollbackInbox *RollbackInboxTransactorSession) Unpause() (*types.Transaction, error) {
	return _RollbackInbox.Contract.Unpause(&_RollbackInbox.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_RollbackInbox *RollbackInboxTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_RollbackInbox *RollbackInboxSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _RollbackInbox.Contract.UpgradeTo(&_RollbackInbox.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_RollbackInbox *RollbackInboxTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _RollbackInbox.Contract.UpgradeTo(&_RollbackInbox.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RollbackInbox *RollbackInboxTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RollbackInbox.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RollbackInbox *RollbackInboxSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.UpgradeToAndCall(&_RollbackInbox.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RollbackInbox *RollbackInboxTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RollbackInbox.Contract.UpgradeToAndCall(&_RollbackInbox.TransactOpts, newImplementation, data)
}

// RollbackInboxAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the RollbackInbox contract.
type RollbackInboxAdminChangedIterator struct {
	Event *RollbackInboxAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxAdminChanged represents a AdminChanged event raised by the RollbackInbox contract.
type RollbackInboxAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_RollbackInbox *RollbackInboxFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*RollbackInboxAdminChangedIterator, error) {

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &RollbackInboxAdminChangedIterator{contract: _RollbackInbox.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_RollbackInbox *RollbackInboxFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *RollbackInboxAdminChanged) (event.Subscription, error) {

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxAdminChanged)
				if err := _RollbackInbox.contract.UnpackLog(event, "AdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_RollbackInbox *RollbackInboxFilterer) ParseAdminChanged(log types.Log) (*RollbackInboxAdminChanged, error) {
	event := new(RollbackInboxAdminChanged)
	if err := _RollbackInbox.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the RollbackInbox contract.
type RollbackInboxBeaconUpgradedIterator struct {
	Event *RollbackInboxBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxBeaconUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxBeaconUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxBeaconUpgraded represents a BeaconUpgraded event raised by the RollbackInbox contract.
type RollbackInboxBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_RollbackInbox *RollbackInboxFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*RollbackInboxBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxBeaconUpgradedIterator{contract: _RollbackInbox.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_RollbackInbox *RollbackInboxFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *RollbackInboxBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxBeaconUpgraded)
				if err := _RollbackInbox.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_RollbackInbox *RollbackInboxFilterer) ParseBeaconUpgraded(log types.Log) (*RollbackInboxBeaconUpgraded, error) {
	event := new(RollbackInboxBeaconUpgraded)
	if err := _RollbackInbox.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxBondInstructionCreatedIterator is returned from FilterBondInstructionCreated and is used to iterate over the raw logs and unpacked data for BondInstructionCreated events raised by the RollbackInbox contract.
type RollbackInboxBondInstructionCreatedIterator struct {
	Event *RollbackInboxBondInstructionCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxBondInstructionCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxBondInstructionCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxBondInstructionCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxBondInstructionCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxBondInstructionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxBondInstructionCreated represents a BondInstructionCreated event raised by the RollbackInbox contract.
type RollbackInboxBondInstructionCreated struct {
	ProposalId      *big.Int
	BondInstruction LibBondsBondInstruction
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBondInstructionCreated is a free log retrieval operation binding the contract event 0x8796b99bbf275a983988181e85b42502d9347327f4dde674320bec3879bcc5e1.
//
// Solidity: event BondInstructionCreated(uint48 indexed proposalId, (uint48,uint8,address,address) bondInstruction)
func (_RollbackInbox *RollbackInboxFilterer) FilterBondInstructionCreated(opts *bind.FilterOpts, proposalId []*big.Int) (*RollbackInboxBondInstructionCreatedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "BondInstructionCreated", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxBondInstructionCreatedIterator{contract: _RollbackInbox.contract, event: "BondInstructionCreated", logs: logs, sub: sub}, nil
}

// WatchBondInstructionCreated is a free log subscription operation binding the contract event 0x8796b99bbf275a983988181e85b42502d9347327f4dde674320bec3879bcc5e1.
//
// Solidity: event BondInstructionCreated(uint48 indexed proposalId, (uint48,uint8,address,address) bondInstruction)
func (_RollbackInbox *RollbackInboxFilterer) WatchBondInstructionCreated(opts *bind.WatchOpts, sink chan<- *RollbackInboxBondInstructionCreated, proposalId []*big.Int) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "BondInstructionCreated", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxBondInstructionCreated)
				if err := _RollbackInbox.contract.UnpackLog(event, "BondInstructionCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBondInstructionCreated is a log parse operation binding the contract event 0x8796b99bbf275a983988181e85b42502d9347327f4dde674320bec3879bcc5e1.
//
// Solidity: event BondInstructionCreated(uint48 indexed proposalId, (uint48,uint8,address,address) bondInstruction)
func (_RollbackInbox *RollbackInboxFilterer) ParseBondInstructionCreated(log types.Log) (*RollbackInboxBondInstructionCreated, error) {
	event := new(RollbackInboxBondInstructionCreated)
	if err := _RollbackInbox.contract.UnpackLog(event, "BondInstructionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxForcedInclusionSavedIterator is returned from FilterForcedInclusionSaved and is used to iterate over the raw logs and unpacked data for ForcedInclusionSaved events raised by the RollbackInbox contract.
type RollbackInboxForcedInclusionSavedIterator struct {
	Event *RollbackInboxForcedInclusionSaved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxForcedInclusionSavedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxForcedInclusionSaved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxForcedInclusionSaved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxForcedInclusionSavedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxForcedInclusionSavedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxForcedInclusionSaved represents a ForcedInclusionSaved event raised by the RollbackInbox contract.
type RollbackInboxForcedInclusionSaved struct {
	ForcedInclusion IForcedInclusionStoreForcedInclusion
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterForcedInclusionSaved is a free log retrieval operation binding the contract event 0x18c4fc1e6ac628dbb537b0375bf0efabf1ff2528af1ec22faa74d2da95c29471.
//
// Solidity: event ForcedInclusionSaved((uint64,(bytes32[],uint24,uint48)) forcedInclusion)
func (_RollbackInbox *RollbackInboxFilterer) FilterForcedInclusionSaved(opts *bind.FilterOpts) (*RollbackInboxForcedInclusionSavedIterator, error) {

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "ForcedInclusionSaved")
	if err != nil {
		return nil, err
	}
	return &RollbackInboxForcedInclusionSavedIterator{contract: _RollbackInbox.contract, event: "ForcedInclusionSaved", logs: logs, sub: sub}, nil
}

// WatchForcedInclusionSaved is a free log subscription operation binding the contract event 0x18c4fc1e6ac628dbb537b0375bf0efabf1ff2528af1ec22faa74d2da95c29471.
//
// Solidity: event ForcedInclusionSaved((uint64,(bytes32[],uint24,uint48)) forcedInclusion)
func (_RollbackInbox *RollbackInboxFilterer) WatchForcedInclusionSaved(opts *bind.WatchOpts, sink chan<- *RollbackInboxForcedInclusionSaved) (event.Subscription, error) {

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "ForcedInclusionSaved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxForcedInclusionSaved)
				if err := _RollbackInbox.contract.UnpackLog(event, "ForcedInclusionSaved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseForcedInclusionSaved is a log parse operation binding the contract event 0x18c4fc1e6ac628dbb537b0375bf0efabf1ff2528af1ec22faa74d2da95c29471.
//
// Solidity: event ForcedInclusionSaved((uint64,(bytes32[],uint24,uint48)) forcedInclusion)
func (_RollbackInbox *RollbackInboxFilterer) ParseForcedInclusionSaved(log types.Log) (*RollbackInboxForcedInclusionSaved, error) {
	event := new(RollbackInboxForcedInclusionSaved)
	if err := _RollbackInbox.contract.UnpackLog(event, "ForcedInclusionSaved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxInboxActivatedIterator is returned from FilterInboxActivated and is used to iterate over the raw logs and unpacked data for InboxActivated events raised by the RollbackInbox contract.
type RollbackInboxInboxActivatedIterator struct {
	Event *RollbackInboxInboxActivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxInboxActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxInboxActivated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxInboxActivated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxInboxActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxInboxActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxInboxActivated represents a InboxActivated event raised by the RollbackInbox contract.
type RollbackInboxInboxActivated struct {
	LastPacayaBlockHash [32]byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterInboxActivated is a free log retrieval operation binding the contract event 0xe4356761c97932c05c3ee0859fb1a5e4f91f7a1d7a3752c7d5a72d5cc6ecb2d2.
//
// Solidity: event InboxActivated(bytes32 lastPacayaBlockHash)
func (_RollbackInbox *RollbackInboxFilterer) FilterInboxActivated(opts *bind.FilterOpts) (*RollbackInboxInboxActivatedIterator, error) {

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "InboxActivated")
	if err != nil {
		return nil, err
	}
	return &RollbackInboxInboxActivatedIterator{contract: _RollbackInbox.contract, event: "InboxActivated", logs: logs, sub: sub}, nil
}

// WatchInboxActivated is a free log subscription operation binding the contract event 0xe4356761c97932c05c3ee0859fb1a5e4f91f7a1d7a3752c7d5a72d5cc6ecb2d2.
//
// Solidity: event InboxActivated(bytes32 lastPacayaBlockHash)
func (_RollbackInbox *RollbackInboxFilterer) WatchInboxActivated(opts *bind.WatchOpts, sink chan<- *RollbackInboxInboxActivated) (event.Subscription, error) {

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "InboxActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxInboxActivated)
				if err := _RollbackInbox.contract.UnpackLog(event, "InboxActivated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInboxActivated is a log parse operation binding the contract event 0xe4356761c97932c05c3ee0859fb1a5e4f91f7a1d7a3752c7d5a72d5cc6ecb2d2.
//
// Solidity: event InboxActivated(bytes32 lastPacayaBlockHash)
func (_RollbackInbox *RollbackInboxFilterer) ParseInboxActivated(log types.Log) (*RollbackInboxInboxActivated, error) {
	event := new(RollbackInboxInboxActivated)
	if err := _RollbackInbox.contract.UnpackLog(event, "InboxActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the RollbackInbox contract.
type RollbackInboxInitializedIterator struct {
	Event *RollbackInboxInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxInitialized represents a Initialized event raised by the RollbackInbox contract.
type RollbackInboxInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_RollbackInbox *RollbackInboxFilterer) FilterInitialized(opts *bind.FilterOpts) (*RollbackInboxInitializedIterator, error) {

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RollbackInboxInitializedIterator{contract: _RollbackInbox.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_RollbackInbox *RollbackInboxFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RollbackInboxInitialized) (event.Subscription, error) {

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxInitialized)
				if err := _RollbackInbox.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_RollbackInbox *RollbackInboxFilterer) ParseInitialized(log types.Log) (*RollbackInboxInitialized, error) {
	event := new(RollbackInboxInitialized)
	if err := _RollbackInbox.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxLimpModeSetIterator is returned from FilterLimpModeSet and is used to iterate over the raw logs and unpacked data for LimpModeSet events raised by the RollbackInbox contract.
type RollbackInboxLimpModeSetIterator struct {
	Event *RollbackInboxLimpModeSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxLimpModeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxLimpModeSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxLimpModeSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxLimpModeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxLimpModeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxLimpModeSet represents a LimpModeSet event raised by the RollbackInbox contract.
type RollbackInboxLimpModeSet struct {
	Enabled bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLimpModeSet is a free log retrieval operation binding the contract event 0xcd24087b96e08a050a16febd67bd7a26feeb56ab1102259cde1fd9d3bebd6681.
//
// Solidity: event LimpModeSet(bool enabled)
func (_RollbackInbox *RollbackInboxFilterer) FilterLimpModeSet(opts *bind.FilterOpts) (*RollbackInboxLimpModeSetIterator, error) {

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "LimpModeSet")
	if err != nil {
		return nil, err
	}
	return &RollbackInboxLimpModeSetIterator{contract: _RollbackInbox.contract, event: "LimpModeSet", logs: logs, sub: sub}, nil
}

// WatchLimpModeSet is a free log subscription operation binding the contract event 0xcd24087b96e08a050a16febd67bd7a26feeb56ab1102259cde1fd9d3bebd6681.
//
// Solidity: event LimpModeSet(bool enabled)
func (_RollbackInbox *RollbackInboxFilterer) WatchLimpModeSet(opts *bind.WatchOpts, sink chan<- *RollbackInboxLimpModeSet) (event.Subscription, error) {

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "LimpModeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxLimpModeSet)
				if err := _RollbackInbox.contract.UnpackLog(event, "LimpModeSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLimpModeSet is a log parse operation binding the contract event 0xcd24087b96e08a050a16febd67bd7a26feeb56ab1102259cde1fd9d3bebd6681.
//
// Solidity: event LimpModeSet(bool enabled)
func (_RollbackInbox *RollbackInboxFilterer) ParseLimpModeSet(log types.Log) (*RollbackInboxLimpModeSet, error) {
	event := new(RollbackInboxLimpModeSet)
	if err := _RollbackInbox.contract.UnpackLog(event, "LimpModeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the RollbackInbox contract.
type RollbackInboxOwnershipTransferStartedIterator struct {
	Event *RollbackInboxOwnershipTransferStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxOwnershipTransferStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxOwnershipTransferStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the RollbackInbox contract.
type RollbackInboxOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_RollbackInbox *RollbackInboxFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RollbackInboxOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxOwnershipTransferStartedIterator{contract: _RollbackInbox.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_RollbackInbox *RollbackInboxFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *RollbackInboxOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxOwnershipTransferStarted)
				if err := _RollbackInbox.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_RollbackInbox *RollbackInboxFilterer) ParseOwnershipTransferStarted(log types.Log) (*RollbackInboxOwnershipTransferStarted, error) {
	event := new(RollbackInboxOwnershipTransferStarted)
	if err := _RollbackInbox.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RollbackInbox contract.
type RollbackInboxOwnershipTransferredIterator struct {
	Event *RollbackInboxOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxOwnershipTransferred represents a OwnershipTransferred event raised by the RollbackInbox contract.
type RollbackInboxOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RollbackInbox *RollbackInboxFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RollbackInboxOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxOwnershipTransferredIterator{contract: _RollbackInbox.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RollbackInbox *RollbackInboxFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RollbackInboxOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxOwnershipTransferred)
				if err := _RollbackInbox.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RollbackInbox *RollbackInboxFilterer) ParseOwnershipTransferred(log types.Log) (*RollbackInboxOwnershipTransferred, error) {
	event := new(RollbackInboxOwnershipTransferred)
	if err := _RollbackInbox.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the RollbackInbox contract.
type RollbackInboxPausedIterator struct {
	Event *RollbackInboxPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxPaused represents a Paused event raised by the RollbackInbox contract.
type RollbackInboxPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RollbackInbox *RollbackInboxFilterer) FilterPaused(opts *bind.FilterOpts) (*RollbackInboxPausedIterator, error) {

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RollbackInboxPausedIterator{contract: _RollbackInbox.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RollbackInbox *RollbackInboxFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RollbackInboxPaused) (event.Subscription, error) {

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxPaused)
				if err := _RollbackInbox.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RollbackInbox *RollbackInboxFilterer) ParsePaused(log types.Log) (*RollbackInboxPaused, error) {
	event := new(RollbackInboxPaused)
	if err := _RollbackInbox.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxProposedIterator is returned from FilterProposed and is used to iterate over the raw logs and unpacked data for Proposed events raised by the RollbackInbox contract.
type RollbackInboxProposedIterator struct {
	Event *RollbackInboxProposed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxProposed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxProposed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxProposed represents a Proposed event raised by the RollbackInbox contract.
type RollbackInboxProposed struct {
	Id                             *big.Int
	Proposer                       common.Address
	ParentProposalHash             [32]byte
	EndOfSubmissionWindowTimestamp *big.Int
	BasefeeSharingPctg             uint8
	Sources                        []IInboxDerivationSource
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterProposed is a free log retrieval operation binding the contract event 0x7c4c4523e17533e451df15762a093e0693a2cd8b279fe54c6cd3777ed5771213.
//
// Solidity: event Proposed(uint48 indexed id, address indexed proposer, bytes32 parentProposalHash, uint48 endOfSubmissionWindowTimestamp, uint8 basefeeSharingPctg, (bool,(bytes32[],uint24,uint48))[] sources)
func (_RollbackInbox *RollbackInboxFilterer) FilterProposed(opts *bind.FilterOpts, id []*big.Int, proposer []common.Address) (*RollbackInboxProposedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "Proposed", idRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxProposedIterator{contract: _RollbackInbox.contract, event: "Proposed", logs: logs, sub: sub}, nil
}

// WatchProposed is a free log subscription operation binding the contract event 0x7c4c4523e17533e451df15762a093e0693a2cd8b279fe54c6cd3777ed5771213.
//
// Solidity: event Proposed(uint48 indexed id, address indexed proposer, bytes32 parentProposalHash, uint48 endOfSubmissionWindowTimestamp, uint8 basefeeSharingPctg, (bool,(bytes32[],uint24,uint48))[] sources)
func (_RollbackInbox *RollbackInboxFilterer) WatchProposed(opts *bind.WatchOpts, sink chan<- *RollbackInboxProposed, id []*big.Int, proposer []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "Proposed", idRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxProposed)
				if err := _RollbackInbox.contract.UnpackLog(event, "Proposed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProposed is a log parse operation binding the contract event 0x7c4c4523e17533e451df15762a093e0693a2cd8b279fe54c6cd3777ed5771213.
//
// Solidity: event Proposed(uint48 indexed id, address indexed proposer, bytes32 parentProposalHash, uint48 endOfSubmissionWindowTimestamp, uint8 basefeeSharingPctg, (bool,(bytes32[],uint24,uint48))[] sources)
func (_RollbackInbox *RollbackInboxFilterer) ParseProposed(log types.Log) (*RollbackInboxProposed, error) {
	event := new(RollbackInboxProposed)
	if err := _RollbackInbox.contract.UnpackLog(event, "Proposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxProvedIterator is returned from FilterProved and is used to iterate over the raw logs and unpacked data for Proved events raised by the RollbackInbox contract.
type RollbackInboxProvedIterator struct {
	Event *RollbackInboxProved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxProvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxProved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxProved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxProvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxProvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxProved represents a Proved event raised by the RollbackInbox contract.
type RollbackInboxProved struct {
	FirstProposalId    *big.Int
	FirstNewProposalId *big.Int
	LastProposalId     *big.Int
	ActualProver       common.Address
	CheckpointSynced   bool
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterProved is a free log retrieval operation binding the contract event 0x7ca0f1e30099488c4ee24e86a6b2c6802e9add6d530919af7aa17db3bcc1cff1.
//
// Solidity: event Proved(uint48 firstProposalId, uint48 firstNewProposalId, uint48 lastProposalId, address indexed actualProver, bool checkpointSynced)
func (_RollbackInbox *RollbackInboxFilterer) FilterProved(opts *bind.FilterOpts, actualProver []common.Address) (*RollbackInboxProvedIterator, error) {

	var actualProverRule []interface{}
	for _, actualProverItem := range actualProver {
		actualProverRule = append(actualProverRule, actualProverItem)
	}

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "Proved", actualProverRule)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxProvedIterator{contract: _RollbackInbox.contract, event: "Proved", logs: logs, sub: sub}, nil
}

// WatchProved is a free log subscription operation binding the contract event 0x7ca0f1e30099488c4ee24e86a6b2c6802e9add6d530919af7aa17db3bcc1cff1.
//
// Solidity: event Proved(uint48 firstProposalId, uint48 firstNewProposalId, uint48 lastProposalId, address indexed actualProver, bool checkpointSynced)
func (_RollbackInbox *RollbackInboxFilterer) WatchProved(opts *bind.WatchOpts, sink chan<- *RollbackInboxProved, actualProver []common.Address) (event.Subscription, error) {

	var actualProverRule []interface{}
	for _, actualProverItem := range actualProver {
		actualProverRule = append(actualProverRule, actualProverItem)
	}

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "Proved", actualProverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxProved)
				if err := _RollbackInbox.contract.UnpackLog(event, "Proved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProved is a log parse operation binding the contract event 0x7ca0f1e30099488c4ee24e86a6b2c6802e9add6d530919af7aa17db3bcc1cff1.
//
// Solidity: event Proved(uint48 firstProposalId, uint48 firstNewProposalId, uint48 lastProposalId, address indexed actualProver, bool checkpointSynced)
func (_RollbackInbox *RollbackInboxFilterer) ParseProved(log types.Log) (*RollbackInboxProved, error) {
	event := new(RollbackInboxProved)
	if err := _RollbackInbox.contract.UnpackLog(event, "Proved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxRollbackedIterator is returned from FilterRollbacked and is used to iterate over the raw logs and unpacked data for Rollbacked events raised by the RollbackInbox contract.
type RollbackInboxRollbackedIterator struct {
	Event *RollbackInboxRollbacked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxRollbackedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxRollbacked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxRollbacked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxRollbackedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxRollbackedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxRollbacked represents a Rollbacked event raised by the RollbackInbox contract.
type RollbackInboxRollbacked struct {
	FirstProposalId *big.Int
	LastProposalId  *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRollbacked is a free log retrieval operation binding the contract event 0x7472922f91878c5cf55de616ca2ed9b7ec32df0b0b588483985979da1f9d988e.
//
// Solidity: event Rollbacked(uint256 firstProposalId, uint256 lastProposalId)
func (_RollbackInbox *RollbackInboxFilterer) FilterRollbacked(opts *bind.FilterOpts) (*RollbackInboxRollbackedIterator, error) {

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "Rollbacked")
	if err != nil {
		return nil, err
	}
	return &RollbackInboxRollbackedIterator{contract: _RollbackInbox.contract, event: "Rollbacked", logs: logs, sub: sub}, nil
}

// WatchRollbacked is a free log subscription operation binding the contract event 0x7472922f91878c5cf55de616ca2ed9b7ec32df0b0b588483985979da1f9d988e.
//
// Solidity: event Rollbacked(uint256 firstProposalId, uint256 lastProposalId)
func (_RollbackInbox *RollbackInboxFilterer) WatchRollbacked(opts *bind.WatchOpts, sink chan<- *RollbackInboxRollbacked) (event.Subscription, error) {

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "Rollbacked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxRollbacked)
				if err := _RollbackInbox.contract.UnpackLog(event, "Rollbacked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRollbacked is a log parse operation binding the contract event 0x7472922f91878c5cf55de616ca2ed9b7ec32df0b0b588483985979da1f9d988e.
//
// Solidity: event Rollbacked(uint256 firstProposalId, uint256 lastProposalId)
func (_RollbackInbox *RollbackInboxFilterer) ParseRollbacked(log types.Log) (*RollbackInboxRollbacked, error) {
	event := new(RollbackInboxRollbacked)
	if err := _RollbackInbox.contract.UnpackLog(event, "Rollbacked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the RollbackInbox contract.
type RollbackInboxUnpausedIterator struct {
	Event *RollbackInboxUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxUnpaused represents a Unpaused event raised by the RollbackInbox contract.
type RollbackInboxUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RollbackInbox *RollbackInboxFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RollbackInboxUnpausedIterator, error) {

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RollbackInboxUnpausedIterator{contract: _RollbackInbox.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RollbackInbox *RollbackInboxFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RollbackInboxUnpaused) (event.Subscription, error) {

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxUnpaused)
				if err := _RollbackInbox.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RollbackInbox *RollbackInboxFilterer) ParseUnpaused(log types.Log) (*RollbackInboxUnpaused, error) {
	event := new(RollbackInboxUnpaused)
	if err := _RollbackInbox.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollbackInboxUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the RollbackInbox contract.
type RollbackInboxUpgradedIterator struct {
	Event *RollbackInboxUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollbackInboxUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollbackInboxUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollbackInboxUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollbackInboxUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollbackInboxUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollbackInboxUpgraded represents a Upgraded event raised by the RollbackInbox contract.
type RollbackInboxUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RollbackInbox *RollbackInboxFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RollbackInboxUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RollbackInbox.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RollbackInboxUpgradedIterator{contract: _RollbackInbox.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RollbackInbox *RollbackInboxFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RollbackInboxUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RollbackInbox.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollbackInboxUpgraded)
				if err := _RollbackInbox.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RollbackInbox *RollbackInboxFilterer) ParseUpgraded(log types.Log) (*RollbackInboxUpgraded, error) {
	event := new(RollbackInboxUpgraded)
	if err := _RollbackInbox.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
