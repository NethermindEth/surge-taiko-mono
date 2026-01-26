// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package surge

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

// IForcedInclusionStoreForcedInclusion is an auto generated low-level Go binding around an user-defined struct.
type IForcedInclusionStoreForcedInclusion struct {
	FeeInGwei uint64
	BlobSlice LibBlobsBlobSlice
}

// IInboxCommitment is an auto generated low-level Go binding around an user-defined struct.
type IInboxCommitment struct {
	FirstProposalId              *big.Int
	FirstProposalParentBlockHash [32]byte
	LastProposalHash             [32]byte
	ActualProver                 common.Address
	EndBlockNumber               *big.Int
	EndStateRoot                 [32]byte
	Transitions                  []IInboxTransition
}

// IInboxConfig is an auto generated low-level Go binding around an user-defined struct.
type IInboxConfig struct {
	ProofVerifier                     common.Address
	ProposerChecker                   common.Address
	ProverWhitelist                   common.Address
	SignalService                     common.Address
	BondToken                         common.Address
	MinBond                           uint64
	LivenessBond                      uint64
	WithdrawalDelay                   *big.Int
	ProvingWindow                     *big.Int
	PermissionlessProvingDelay        *big.Int
	MaxProofSubmissionDelay           *big.Int
	RingBufferSize                    *big.Int
	BasefeeSharingPctg                uint8
	ForcedInclusionDelay              uint16
	ForcedInclusionFeeInGwei          uint64
	ForcedInclusionFeeDoubleThreshold uint64
	PermissionlessInclusionMultiplier uint8
}

// IInboxCoreState is an auto generated low-level Go binding around an user-defined struct.
type IInboxCoreState struct {
	NextProposalId          *big.Int
	LastProposalBlockId     *big.Int
	LastFinalizedProposalId *big.Int
	LastFinalizedTimestamp  *big.Int
	LastCheckpointTimestamp *big.Int
	LastFinalizedBlockHash  [32]byte
}

// IInboxDerivationSource is an auto generated low-level Go binding around an user-defined struct.
type IInboxDerivationSource struct {
	IsForcedInclusion bool
	BlobSlice         LibBlobsBlobSlice
}

// IInboxProposal is an auto generated low-level Go binding around an user-defined struct.
type IInboxProposal struct {
	Id                             *big.Int
	Timestamp                      *big.Int
	EndOfSubmissionWindowTimestamp *big.Int
	Proposer                       common.Address
	ParentProposalHash             [32]byte
	OriginBlockNumber              *big.Int
	OriginBlockHash                [32]byte
	BasefeeSharingPctg             uint8
	Sources                        []IInboxDerivationSource
	SignalSlotsHash                [32]byte
}

// IInboxProposeInput is an auto generated low-level Go binding around an user-defined struct.
type IInboxProposeInput struct {
	Deadline            *big.Int
	BlobReference       LibBlobsBlobReference
	NumForcedInclusions uint16
}

// IInboxProveInput is an auto generated low-level Go binding around an user-defined struct.
type IInboxProveInput struct {
	Commitment IInboxCommitment
}

// IInboxTransition is an auto generated low-level Go binding around an user-defined struct.
type IInboxTransition struct {
	Proposer  common.Address
	Timestamp *big.Int
	BlockHash [32]byte
}

// LibBlobsBlobReference is an auto generated low-level Go binding around an user-defined struct.
type LibBlobsBlobReference struct {
	BlobStartIndex uint16
	NumBlobs       uint16
	Offset         *big.Int
}

// LibBlobsBlobSlice is an auto generated low-level Go binding around an user-defined struct.
type LibBlobsBlobSlice struct {
	BlobHashes [][32]byte
	Offset     *big.Int
	Timestamp  *big.Int
}

// SurgeInboxClientMetaData contains all meta data concerning the SurgeInboxClient contract.
var SurgeInboxClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_proofVerifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_proposerChecker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_proverWhitelist\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_bondToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"activate\",\"inputs\":[{\"name\":\"_lastPacayaBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"activationTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelWithdrawal\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decodeProposeInput\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"input_\",\"type\":\"tuple\",\"internalType\":\"structIInbox.ProposeInput\",\"components\":[{\"name\":\"deadline\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blobReference\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobReference\",\"components\":[{\"name\":\"blobStartIndex\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"numBlobs\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]},{\"name\":\"numForcedInclusions\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"decodeProveInput\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"input_\",\"type\":\"tuple\",\"internalType\":\"structIInbox.ProveInput\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Commitment\",\"components\":[{\"name\":\"firstProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"firstProposalParentBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lastProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actualProver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"endBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"transitions\",\"type\":\"tuple[]\",\"internalType\":\"structIInbox.Transition[]\",\"components\":[{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositTo\",\"inputs\":[{\"name\":\"_recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"encodeProposeInput\",\"inputs\":[{\"name\":\"_input\",\"type\":\"tuple\",\"internalType\":\"structIInbox.ProposeInput\",\"components\":[{\"name\":\"deadline\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blobReference\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobReference\",\"components\":[{\"name\":\"blobStartIndex\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"numBlobs\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]},{\"name\":\"numForcedInclusions\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"outputs\":[{\"name\":\"encoded_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeProveInput\",\"inputs\":[{\"name\":\"_input\",\"type\":\"tuple\",\"internalType\":\"structIInbox.ProveInput\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Commitment\",\"components\":[{\"name\":\"firstProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"firstProposalParentBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lastProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actualProver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"endBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"transitions\",\"type\":\"tuple[]\",\"internalType\":\"structIInbox.Transition[]\",\"components\":[{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}]}],\"outputs\":[{\"name\":\"encoded_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getBond\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"bond_\",\"type\":\"tuple\",\"internalType\":\"structIBondManager.Bond\",\"components\":[{\"name\":\"balance\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"withdrawalRequestedAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"config_\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Config\",\"components\":[{\"name\":\"proofVerifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposerChecker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proverWhitelist\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signalService\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bondToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minBond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"livenessBond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"withdrawalDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"provingWindow\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"permissionlessProvingDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"maxProofSubmissionDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"ringBufferSize\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"forcedInclusionDelay\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"forcedInclusionFeeInGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"forcedInclusionFeeDoubleThreshold\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"permissionlessInclusionMultiplier\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCoreState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIInbox.CoreState\",\"components\":[{\"name\":\"nextProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastProposalBlockId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastFinalizedProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastFinalizedTimestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastCheckpointTimestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"lastFinalizedBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentForcedInclusionFee\",\"inputs\":[],\"outputs\":[{\"name\":\"feeInGwei_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getForcedInclusionState\",\"inputs\":[],\"outputs\":[{\"name\":\"head_\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"tail_\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getForcedInclusions\",\"inputs\":[{\"name\":\"_start\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"_maxCount\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[{\"name\":\"inclusions_\",\"type\":\"tuple[]\",\"internalType\":\"structIForcedInclusionStore.ForcedInclusion[]\",\"components\":[{\"name\":\"feeInGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProposalHash\",\"inputs\":[{\"name\":\"_proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashCommitment\",\"inputs\":[{\"name\":\"_commitment\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Commitment\",\"components\":[{\"name\":\"firstProposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"firstProposalParentBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lastProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actualProver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"endBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"transitions\",\"type\":\"tuple[]\",\"internalType\":\"structIInbox.Transition[]\",\"components\":[{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"hashProposal\",\"inputs\":[{\"name\":\"_proposal\",\"type\":\"tuple\",\"internalType\":\"structIInbox.Proposal\",\"components\":[{\"name\":\"id\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endOfSubmissionWindowTimestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"parentProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"originBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"originBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"sources\",\"type\":\"tuple[]\",\"internalType\":\"structIInbox.DerivationSource[]\",\"components\":[{\"name\":\"isForcedInclusion\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]},{\"name\":\"signalSlotsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"impl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inNonReentrant\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"init\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"propose\",\"inputs\":[{\"name\":\"_lookahead\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proposeWithProof\",\"inputs\":[{\"name\":\"_lookahead\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_signalSlots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"prove\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestWithdrawal\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"saveForcedInclusion\",\"inputs\":[{\"name\":\"_blobReference\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobReference\",\"components\":[{\"name\":\"blobStartIndex\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"numBlobs\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminChanged\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BeaconUpgraded\",\"inputs\":[{\"name\":\"beacon\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BondDeposited\",\"inputs\":[{\"name\":\"depositor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BondWithdrawn\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ForcedInclusionSaved\",\"inputs\":[{\"name\":\"forcedInclusion\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIForcedInclusionStore.ForcedInclusion\",\"components\":[{\"name\":\"feeInGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboxActivated\",\"inputs\":[{\"name\":\"lastPacayaBlockHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LivenessBondSettled\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"livenessBond\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"credited\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"slashed\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Proposed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint48\",\"indexed\":true,\"internalType\":\"uint48\"},{\"name\":\"proposer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"parentProposalHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"endOfSubmissionWindowTimestamp\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"sources\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structIInbox.DerivationSource[]\",\"components\":[{\"name\":\"isForcedInclusion\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"structLibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]},{\"name\":\"signalSlotsHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Proved\",\"inputs\":[{\"name\":\"firstProposalId\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"firstNewProposalId\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"lastProposalId\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"actualProver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalCancelled\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalRequested\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawableAt\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ACCESS_DENIED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ActivationRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BlobNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotProposeInCurrentBlock\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DeadlineExceeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ETH_TRANSFER_FAILED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyBatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FUNC_NOT_IMPLEMENTED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FirstProposalIdTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"INVALID_PAUSE_STATUS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectProposalCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBond\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientETH\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LastProposalAlreadyFinalized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LastProposalHashMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LastProposalIdTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthExceedsUint16\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustMaintainMinBond\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoBlobs\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoBondToWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoWithdrawalRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughCapacity\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParentBlockHashMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProverNotWhitelisted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"REENTRANT_CALL\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnprocessedForcedInclusionIsDue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawalAlreadyRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_ADDRESS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_VALUE\",\"inputs\":[]}]",
}

// SurgeInboxClientABI is the input ABI used to generate the binding from.
// Deprecated: Use SurgeInboxClientMetaData.ABI instead.
var SurgeInboxClientABI = SurgeInboxClientMetaData.ABI

// SurgeInboxClient is an auto generated Go binding around an Ethereum contract.
type SurgeInboxClient struct {
	SurgeInboxClientCaller     // Read-only binding to the contract
	SurgeInboxClientTransactor // Write-only binding to the contract
	SurgeInboxClientFilterer   // Log filterer for contract events
}

// SurgeInboxClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type SurgeInboxClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeInboxClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SurgeInboxClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeInboxClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SurgeInboxClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeInboxClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SurgeInboxClientSession struct {
	Contract     *SurgeInboxClient // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SurgeInboxClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SurgeInboxClientCallerSession struct {
	Contract *SurgeInboxClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SurgeInboxClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SurgeInboxClientTransactorSession struct {
	Contract     *SurgeInboxClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SurgeInboxClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type SurgeInboxClientRaw struct {
	Contract *SurgeInboxClient // Generic contract binding to access the raw methods on
}

// SurgeInboxClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SurgeInboxClientCallerRaw struct {
	Contract *SurgeInboxClientCaller // Generic read-only contract binding to access the raw methods on
}

// SurgeInboxClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SurgeInboxClientTransactorRaw struct {
	Contract *SurgeInboxClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSurgeInboxClient creates a new instance of SurgeInboxClient, bound to a specific deployed contract.
func NewSurgeInboxClient(address common.Address, backend bind.ContractBackend) (*SurgeInboxClient, error) {
	contract, err := bindSurgeInboxClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClient{SurgeInboxClientCaller: SurgeInboxClientCaller{contract: contract}, SurgeInboxClientTransactor: SurgeInboxClientTransactor{contract: contract}, SurgeInboxClientFilterer: SurgeInboxClientFilterer{contract: contract}}, nil
}

// NewSurgeInboxClientCaller creates a new read-only instance of SurgeInboxClient, bound to a specific deployed contract.
func NewSurgeInboxClientCaller(address common.Address, caller bind.ContractCaller) (*SurgeInboxClientCaller, error) {
	contract, err := bindSurgeInboxClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientCaller{contract: contract}, nil
}

// NewSurgeInboxClientTransactor creates a new write-only instance of SurgeInboxClient, bound to a specific deployed contract.
func NewSurgeInboxClientTransactor(address common.Address, transactor bind.ContractTransactor) (*SurgeInboxClientTransactor, error) {
	contract, err := bindSurgeInboxClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientTransactor{contract: contract}, nil
}

// NewSurgeInboxClientFilterer creates a new log filterer instance of SurgeInboxClient, bound to a specific deployed contract.
func NewSurgeInboxClientFilterer(address common.Address, filterer bind.ContractFilterer) (*SurgeInboxClientFilterer, error) {
	contract, err := bindSurgeInboxClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientFilterer{contract: contract}, nil
}

// bindSurgeInboxClient binds a generic wrapper to an already deployed contract.
func bindSurgeInboxClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SurgeInboxClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SurgeInboxClient *SurgeInboxClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SurgeInboxClient.Contract.SurgeInboxClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SurgeInboxClient *SurgeInboxClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.SurgeInboxClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SurgeInboxClient *SurgeInboxClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.SurgeInboxClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SurgeInboxClient *SurgeInboxClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SurgeInboxClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SurgeInboxClient *SurgeInboxClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SurgeInboxClient *SurgeInboxClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.contract.Transact(opts, method, params...)
}

// ActivationTimestamp is a free data retrieval call binding the contract method 0x0423c7de.
//
// Solidity: function activationTimestamp() view returns(uint48)
func (_SurgeInboxClient *SurgeInboxClientCaller) ActivationTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "activationTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActivationTimestamp is a free data retrieval call binding the contract method 0x0423c7de.
//
// Solidity: function activationTimestamp() view returns(uint48)
func (_SurgeInboxClient *SurgeInboxClientSession) ActivationTimestamp() (*big.Int, error) {
	return _SurgeInboxClient.Contract.ActivationTimestamp(&_SurgeInboxClient.CallOpts)
}

// ActivationTimestamp is a free data retrieval call binding the contract method 0x0423c7de.
//
// Solidity: function activationTimestamp() view returns(uint48)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) ActivationTimestamp() (*big.Int, error) {
	return _SurgeInboxClient.Contract.ActivationTimestamp(&_SurgeInboxClient.CallOpts)
}

// DecodeProposeInput is a free data retrieval call binding the contract method 0xafb63ad4.
//
// Solidity: function decodeProposeInput(bytes _data) pure returns((uint48,(uint16,uint16,uint24),uint16) input_)
func (_SurgeInboxClient *SurgeInboxClientCaller) DecodeProposeInput(opts *bind.CallOpts, _data []byte) (IInboxProposeInput, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "decodeProposeInput", _data)

	if err != nil {
		return *new(IInboxProposeInput), err
	}

	out0 := *abi.ConvertType(out[0], new(IInboxProposeInput)).(*IInboxProposeInput)

	return out0, err

}

// DecodeProposeInput is a free data retrieval call binding the contract method 0xafb63ad4.
//
// Solidity: function decodeProposeInput(bytes _data) pure returns((uint48,(uint16,uint16,uint24),uint16) input_)
func (_SurgeInboxClient *SurgeInboxClientSession) DecodeProposeInput(_data []byte) (IInboxProposeInput, error) {
	return _SurgeInboxClient.Contract.DecodeProposeInput(&_SurgeInboxClient.CallOpts, _data)
}

// DecodeProposeInput is a free data retrieval call binding the contract method 0xafb63ad4.
//
// Solidity: function decodeProposeInput(bytes _data) pure returns((uint48,(uint16,uint16,uint24),uint16) input_)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) DecodeProposeInput(_data []byte) (IInboxProposeInput, error) {
	return _SurgeInboxClient.Contract.DecodeProposeInput(&_SurgeInboxClient.CallOpts, _data)
}

// DecodeProveInput is a free data retrieval call binding the contract method 0xedbacd44.
//
// Solidity: function decodeProveInput(bytes _data) pure returns(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,uint48,bytes32)[])) input_)
func (_SurgeInboxClient *SurgeInboxClientCaller) DecodeProveInput(opts *bind.CallOpts, _data []byte) (IInboxProveInput, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "decodeProveInput", _data)

	if err != nil {
		return *new(IInboxProveInput), err
	}

	out0 := *abi.ConvertType(out[0], new(IInboxProveInput)).(*IInboxProveInput)

	return out0, err

}

// DecodeProveInput is a free data retrieval call binding the contract method 0xedbacd44.
//
// Solidity: function decodeProveInput(bytes _data) pure returns(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,uint48,bytes32)[])) input_)
func (_SurgeInboxClient *SurgeInboxClientSession) DecodeProveInput(_data []byte) (IInboxProveInput, error) {
	return _SurgeInboxClient.Contract.DecodeProveInput(&_SurgeInboxClient.CallOpts, _data)
}

// DecodeProveInput is a free data retrieval call binding the contract method 0xedbacd44.
//
// Solidity: function decodeProveInput(bytes _data) pure returns(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,uint48,bytes32)[])) input_)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) DecodeProveInput(_data []byte) (IInboxProveInput, error) {
	return _SurgeInboxClient.Contract.DecodeProveInput(&_SurgeInboxClient.CallOpts, _data)
}

// EncodeProposeInput is a free data retrieval call binding the contract method 0x1275a673.
//
// Solidity: function encodeProposeInput((uint48,(uint16,uint16,uint24),uint16) _input) pure returns(bytes encoded_)
func (_SurgeInboxClient *SurgeInboxClientCaller) EncodeProposeInput(opts *bind.CallOpts, _input IInboxProposeInput) ([]byte, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "encodeProposeInput", _input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeProposeInput is a free data retrieval call binding the contract method 0x1275a673.
//
// Solidity: function encodeProposeInput((uint48,(uint16,uint16,uint24),uint16) _input) pure returns(bytes encoded_)
func (_SurgeInboxClient *SurgeInboxClientSession) EncodeProposeInput(_input IInboxProposeInput) ([]byte, error) {
	return _SurgeInboxClient.Contract.EncodeProposeInput(&_SurgeInboxClient.CallOpts, _input)
}

// EncodeProposeInput is a free data retrieval call binding the contract method 0x1275a673.
//
// Solidity: function encodeProposeInput((uint48,(uint16,uint16,uint24),uint16) _input) pure returns(bytes encoded_)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) EncodeProposeInput(_input IInboxProposeInput) ([]byte, error) {
	return _SurgeInboxClient.Contract.EncodeProposeInput(&_SurgeInboxClient.CallOpts, _input)
}

// EncodeProveInput is a free data retrieval call binding the contract method 0xc5954597.
//
// Solidity: function encodeProveInput(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,uint48,bytes32)[])) _input) pure returns(bytes encoded_)
func (_SurgeInboxClient *SurgeInboxClientCaller) EncodeProveInput(opts *bind.CallOpts, _input IInboxProveInput) ([]byte, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "encodeProveInput", _input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeProveInput is a free data retrieval call binding the contract method 0xc5954597.
//
// Solidity: function encodeProveInput(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,uint48,bytes32)[])) _input) pure returns(bytes encoded_)
func (_SurgeInboxClient *SurgeInboxClientSession) EncodeProveInput(_input IInboxProveInput) ([]byte, error) {
	return _SurgeInboxClient.Contract.EncodeProveInput(&_SurgeInboxClient.CallOpts, _input)
}

// EncodeProveInput is a free data retrieval call binding the contract method 0xc5954597.
//
// Solidity: function encodeProveInput(((uint48,bytes32,bytes32,address,uint48,bytes32,(address,uint48,bytes32)[])) _input) pure returns(bytes encoded_)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) EncodeProveInput(_input IInboxProveInput) ([]byte, error) {
	return _SurgeInboxClient.Contract.EncodeProveInput(&_SurgeInboxClient.CallOpts, _input)
}

// GetBond is a free data retrieval call binding the contract method 0x0d8912f3.
//
// Solidity: function getBond(address _address) view returns((uint64,uint48) bond_)
func (_SurgeInboxClient *SurgeInboxClientCaller) GetBond(opts *bind.CallOpts, _address common.Address) (IBondManagerBond, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "getBond", _address)

	if err != nil {
		return *new(IBondManagerBond), err
	}

	out0 := *abi.ConvertType(out[0], new(IBondManagerBond)).(*IBondManagerBond)

	return out0, err

}

// GetBond is a free data retrieval call binding the contract method 0x0d8912f3.
//
// Solidity: function getBond(address _address) view returns((uint64,uint48) bond_)
func (_SurgeInboxClient *SurgeInboxClientSession) GetBond(_address common.Address) (IBondManagerBond, error) {
	return _SurgeInboxClient.Contract.GetBond(&_SurgeInboxClient.CallOpts, _address)
}

// GetBond is a free data retrieval call binding the contract method 0x0d8912f3.
//
// Solidity: function getBond(address _address) view returns((uint64,uint48) bond_)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) GetBond(_address common.Address) (IBondManagerBond, error) {
	return _SurgeInboxClient.Contract.GetBond(&_SurgeInboxClient.CallOpts, _address)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,address,address,address,uint64,uint64,uint48,uint48,uint48,uint48,uint48,uint8,uint16,uint64,uint64,uint8) config_)
func (_SurgeInboxClient *SurgeInboxClientCaller) GetConfig(opts *bind.CallOpts) (IInboxConfig, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(IInboxConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IInboxConfig)).(*IInboxConfig)

	return out0, err

}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,address,address,address,uint64,uint64,uint48,uint48,uint48,uint48,uint48,uint8,uint16,uint64,uint64,uint8) config_)
func (_SurgeInboxClient *SurgeInboxClientSession) GetConfig() (IInboxConfig, error) {
	return _SurgeInboxClient.Contract.GetConfig(&_SurgeInboxClient.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,address,address,address,uint64,uint64,uint48,uint48,uint48,uint48,uint48,uint8,uint16,uint64,uint64,uint8) config_)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) GetConfig() (IInboxConfig, error) {
	return _SurgeInboxClient.Contract.GetConfig(&_SurgeInboxClient.CallOpts)
}

// GetCoreState is a free data retrieval call binding the contract method 0x6aa6a01a.
//
// Solidity: function getCoreState() view returns((uint48,uint48,uint48,uint48,uint48,bytes32))
func (_SurgeInboxClient *SurgeInboxClientCaller) GetCoreState(opts *bind.CallOpts) (IInboxCoreState, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "getCoreState")

	if err != nil {
		return *new(IInboxCoreState), err
	}

	out0 := *abi.ConvertType(out[0], new(IInboxCoreState)).(*IInboxCoreState)

	return out0, err

}

// GetCoreState is a free data retrieval call binding the contract method 0x6aa6a01a.
//
// Solidity: function getCoreState() view returns((uint48,uint48,uint48,uint48,uint48,bytes32))
func (_SurgeInboxClient *SurgeInboxClientSession) GetCoreState() (IInboxCoreState, error) {
	return _SurgeInboxClient.Contract.GetCoreState(&_SurgeInboxClient.CallOpts)
}

// GetCoreState is a free data retrieval call binding the contract method 0x6aa6a01a.
//
// Solidity: function getCoreState() view returns((uint48,uint48,uint48,uint48,uint48,bytes32))
func (_SurgeInboxClient *SurgeInboxClientCallerSession) GetCoreState() (IInboxCoreState, error) {
	return _SurgeInboxClient.Contract.GetCoreState(&_SurgeInboxClient.CallOpts)
}

// GetCurrentForcedInclusionFee is a free data retrieval call binding the contract method 0xe3053335.
//
// Solidity: function getCurrentForcedInclusionFee() view returns(uint64 feeInGwei_)
func (_SurgeInboxClient *SurgeInboxClientCaller) GetCurrentForcedInclusionFee(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "getCurrentForcedInclusionFee")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetCurrentForcedInclusionFee is a free data retrieval call binding the contract method 0xe3053335.
//
// Solidity: function getCurrentForcedInclusionFee() view returns(uint64 feeInGwei_)
func (_SurgeInboxClient *SurgeInboxClientSession) GetCurrentForcedInclusionFee() (uint64, error) {
	return _SurgeInboxClient.Contract.GetCurrentForcedInclusionFee(&_SurgeInboxClient.CallOpts)
}

// GetCurrentForcedInclusionFee is a free data retrieval call binding the contract method 0xe3053335.
//
// Solidity: function getCurrentForcedInclusionFee() view returns(uint64 feeInGwei_)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) GetCurrentForcedInclusionFee() (uint64, error) {
	return _SurgeInboxClient.Contract.GetCurrentForcedInclusionFee(&_SurgeInboxClient.CallOpts)
}

// GetForcedInclusionState is a free data retrieval call binding the contract method 0x5ccc1718.
//
// Solidity: function getForcedInclusionState() view returns(uint48 head_, uint48 tail_)
func (_SurgeInboxClient *SurgeInboxClientCaller) GetForcedInclusionState(opts *bind.CallOpts) (struct {
	Head *big.Int
	Tail *big.Int
}, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "getForcedInclusionState")

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
func (_SurgeInboxClient *SurgeInboxClientSession) GetForcedInclusionState() (struct {
	Head *big.Int
	Tail *big.Int
}, error) {
	return _SurgeInboxClient.Contract.GetForcedInclusionState(&_SurgeInboxClient.CallOpts)
}

// GetForcedInclusionState is a free data retrieval call binding the contract method 0x5ccc1718.
//
// Solidity: function getForcedInclusionState() view returns(uint48 head_, uint48 tail_)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) GetForcedInclusionState() (struct {
	Head *big.Int
	Tail *big.Int
}, error) {
	return _SurgeInboxClient.Contract.GetForcedInclusionState(&_SurgeInboxClient.CallOpts)
}

// GetForcedInclusions is a free data retrieval call binding the contract method 0x40df9866.
//
// Solidity: function getForcedInclusions(uint48 _start, uint48 _maxCount) view returns((uint64,(bytes32[],uint24,uint48))[] inclusions_)
func (_SurgeInboxClient *SurgeInboxClientCaller) GetForcedInclusions(opts *bind.CallOpts, _start *big.Int, _maxCount *big.Int) ([]IForcedInclusionStoreForcedInclusion, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "getForcedInclusions", _start, _maxCount)

	if err != nil {
		return *new([]IForcedInclusionStoreForcedInclusion), err
	}

	out0 := *abi.ConvertType(out[0], new([]IForcedInclusionStoreForcedInclusion)).(*[]IForcedInclusionStoreForcedInclusion)

	return out0, err

}

// GetForcedInclusions is a free data retrieval call binding the contract method 0x40df9866.
//
// Solidity: function getForcedInclusions(uint48 _start, uint48 _maxCount) view returns((uint64,(bytes32[],uint24,uint48))[] inclusions_)
func (_SurgeInboxClient *SurgeInboxClientSession) GetForcedInclusions(_start *big.Int, _maxCount *big.Int) ([]IForcedInclusionStoreForcedInclusion, error) {
	return _SurgeInboxClient.Contract.GetForcedInclusions(&_SurgeInboxClient.CallOpts, _start, _maxCount)
}

// GetForcedInclusions is a free data retrieval call binding the contract method 0x40df9866.
//
// Solidity: function getForcedInclusions(uint48 _start, uint48 _maxCount) view returns((uint64,(bytes32[],uint24,uint48))[] inclusions_)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) GetForcedInclusions(_start *big.Int, _maxCount *big.Int) ([]IForcedInclusionStoreForcedInclusion, error) {
	return _SurgeInboxClient.Contract.GetForcedInclusions(&_SurgeInboxClient.CallOpts, _start, _maxCount)
}

// GetProposalHash is a free data retrieval call binding the contract method 0xa834725a.
//
// Solidity: function getProposalHash(uint256 _proposalId) view returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientCaller) GetProposalHash(opts *bind.CallOpts, _proposalId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "getProposalHash", _proposalId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetProposalHash is a free data retrieval call binding the contract method 0xa834725a.
//
// Solidity: function getProposalHash(uint256 _proposalId) view returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientSession) GetProposalHash(_proposalId *big.Int) ([32]byte, error) {
	return _SurgeInboxClient.Contract.GetProposalHash(&_SurgeInboxClient.CallOpts, _proposalId)
}

// GetProposalHash is a free data retrieval call binding the contract method 0xa834725a.
//
// Solidity: function getProposalHash(uint256 _proposalId) view returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) GetProposalHash(_proposalId *big.Int) ([32]byte, error) {
	return _SurgeInboxClient.Contract.GetProposalHash(&_SurgeInboxClient.CallOpts, _proposalId)
}

// HashCommitment is a free data retrieval call binding the contract method 0xf954ab92.
//
// Solidity: function hashCommitment((uint48,bytes32,bytes32,address,uint48,bytes32,(address,uint48,bytes32)[]) _commitment) pure returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientCaller) HashCommitment(opts *bind.CallOpts, _commitment IInboxCommitment) ([32]byte, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "hashCommitment", _commitment)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashCommitment is a free data retrieval call binding the contract method 0xf954ab92.
//
// Solidity: function hashCommitment((uint48,bytes32,bytes32,address,uint48,bytes32,(address,uint48,bytes32)[]) _commitment) pure returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientSession) HashCommitment(_commitment IInboxCommitment) ([32]byte, error) {
	return _SurgeInboxClient.Contract.HashCommitment(&_SurgeInboxClient.CallOpts, _commitment)
}

// HashCommitment is a free data retrieval call binding the contract method 0xf954ab92.
//
// Solidity: function hashCommitment((uint48,bytes32,bytes32,address,uint48,bytes32,(address,uint48,bytes32)[]) _commitment) pure returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) HashCommitment(_commitment IInboxCommitment) ([32]byte, error) {
	return _SurgeInboxClient.Contract.HashCommitment(&_SurgeInboxClient.CallOpts, _commitment)
}

// HashProposal is a free data retrieval call binding the contract method 0xcc8eabd7.
//
// Solidity: function hashProposal((uint48,uint48,uint48,address,bytes32,uint48,bytes32,uint8,(bool,(bytes32[],uint24,uint48))[],bytes32) _proposal) pure returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientCaller) HashProposal(opts *bind.CallOpts, _proposal IInboxProposal) ([32]byte, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "hashProposal", _proposal)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashProposal is a free data retrieval call binding the contract method 0xcc8eabd7.
//
// Solidity: function hashProposal((uint48,uint48,uint48,address,bytes32,uint48,bytes32,uint8,(bool,(bytes32[],uint24,uint48))[],bytes32) _proposal) pure returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientSession) HashProposal(_proposal IInboxProposal) ([32]byte, error) {
	return _SurgeInboxClient.Contract.HashProposal(&_SurgeInboxClient.CallOpts, _proposal)
}

// HashProposal is a free data retrieval call binding the contract method 0xcc8eabd7.
//
// Solidity: function hashProposal((uint48,uint48,uint48,address,bytes32,uint48,bytes32,uint8,(bool,(bytes32[],uint24,uint48))[],bytes32) _proposal) pure returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) HashProposal(_proposal IInboxProposal) ([32]byte, error) {
	return _SurgeInboxClient.Contract.HashProposal(&_SurgeInboxClient.CallOpts, _proposal)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientCaller) Impl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "impl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientSession) Impl() (common.Address, error) {
	return _SurgeInboxClient.Contract.Impl(&_SurgeInboxClient.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) Impl() (common.Address, error) {
	return _SurgeInboxClient.Contract.Impl(&_SurgeInboxClient.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_SurgeInboxClient *SurgeInboxClientCaller) InNonReentrant(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "inNonReentrant")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_SurgeInboxClient *SurgeInboxClientSession) InNonReentrant() (bool, error) {
	return _SurgeInboxClient.Contract.InNonReentrant(&_SurgeInboxClient.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) InNonReentrant() (bool, error) {
	return _SurgeInboxClient.Contract.InNonReentrant(&_SurgeInboxClient.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientSession) Owner() (common.Address, error) {
	return _SurgeInboxClient.Contract.Owner(&_SurgeInboxClient.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) Owner() (common.Address, error) {
	return _SurgeInboxClient.Contract.Owner(&_SurgeInboxClient.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SurgeInboxClient *SurgeInboxClientCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SurgeInboxClient *SurgeInboxClientSession) Paused() (bool, error) {
	return _SurgeInboxClient.Contract.Paused(&_SurgeInboxClient.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) Paused() (bool, error) {
	return _SurgeInboxClient.Contract.Paused(&_SurgeInboxClient.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientSession) PendingOwner() (common.Address, error) {
	return _SurgeInboxClient.Contract.PendingOwner(&_SurgeInboxClient.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) PendingOwner() (common.Address, error) {
	return _SurgeInboxClient.Contract.PendingOwner(&_SurgeInboxClient.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientSession) ProxiableUUID() ([32]byte, error) {
	return _SurgeInboxClient.Contract.ProxiableUUID(&_SurgeInboxClient.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) ProxiableUUID() ([32]byte, error) {
	return _SurgeInboxClient.Contract.ProxiableUUID(&_SurgeInboxClient.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientCaller) Resolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeInboxClient.contract.Call(opts, &out, "resolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientSession) Resolver() (common.Address, error) {
	return _SurgeInboxClient.Contract.Resolver(&_SurgeInboxClient.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_SurgeInboxClient *SurgeInboxClientCallerSession) Resolver() (common.Address, error) {
	return _SurgeInboxClient.Contract.Resolver(&_SurgeInboxClient.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SurgeInboxClient *SurgeInboxClientSession) AcceptOwnership() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.AcceptOwnership(&_SurgeInboxClient.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.AcceptOwnership(&_SurgeInboxClient.TransactOpts)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _lastPacayaBlockHash) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) Activate(opts *bind.TransactOpts, _lastPacayaBlockHash [32]byte) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "activate", _lastPacayaBlockHash)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _lastPacayaBlockHash) returns()
func (_SurgeInboxClient *SurgeInboxClientSession) Activate(_lastPacayaBlockHash [32]byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Activate(&_SurgeInboxClient.TransactOpts, _lastPacayaBlockHash)
}

// Activate is a paid mutator transaction binding the contract method 0x59db6e85.
//
// Solidity: function activate(bytes32 _lastPacayaBlockHash) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) Activate(_lastPacayaBlockHash [32]byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Activate(&_SurgeInboxClient.TransactOpts, _lastPacayaBlockHash)
}

// CancelWithdrawal is a paid mutator transaction binding the contract method 0x22611280.
//
// Solidity: function cancelWithdrawal() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) CancelWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "cancelWithdrawal")
}

// CancelWithdrawal is a paid mutator transaction binding the contract method 0x22611280.
//
// Solidity: function cancelWithdrawal() returns()
func (_SurgeInboxClient *SurgeInboxClientSession) CancelWithdrawal() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.CancelWithdrawal(&_SurgeInboxClient.TransactOpts)
}

// CancelWithdrawal is a paid mutator transaction binding the contract method 0x22611280.
//
// Solidity: function cancelWithdrawal() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) CancelWithdrawal() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.CancelWithdrawal(&_SurgeInboxClient.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x13765838.
//
// Solidity: function deposit(uint64 _amount) payable returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) Deposit(opts *bind.TransactOpts, _amount uint64) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "deposit", _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x13765838.
//
// Solidity: function deposit(uint64 _amount) payable returns()
func (_SurgeInboxClient *SurgeInboxClientSession) Deposit(_amount uint64) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Deposit(&_SurgeInboxClient.TransactOpts, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x13765838.
//
// Solidity: function deposit(uint64 _amount) payable returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) Deposit(_amount uint64) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Deposit(&_SurgeInboxClient.TransactOpts, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xefba83c9.
//
// Solidity: function depositTo(address _recipient, uint64 _amount) payable returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) DepositTo(opts *bind.TransactOpts, _recipient common.Address, _amount uint64) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "depositTo", _recipient, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xefba83c9.
//
// Solidity: function depositTo(address _recipient, uint64 _amount) payable returns()
func (_SurgeInboxClient *SurgeInboxClientSession) DepositTo(_recipient common.Address, _amount uint64) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.DepositTo(&_SurgeInboxClient.TransactOpts, _recipient, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xefba83c9.
//
// Solidity: function depositTo(address _recipient, uint64 _amount) payable returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) DepositTo(_recipient common.Address, _amount uint64) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.DepositTo(&_SurgeInboxClient.TransactOpts, _recipient, _amount)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) Init(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "init", _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_SurgeInboxClient *SurgeInboxClientSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Init(&_SurgeInboxClient.TransactOpts, _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Init(&_SurgeInboxClient.TransactOpts, _owner)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SurgeInboxClient *SurgeInboxClientSession) Pause() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Pause(&_SurgeInboxClient.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) Pause() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Pause(&_SurgeInboxClient.TransactOpts)
}

// Propose is a paid mutator transaction binding the contract method 0x9791e644.
//
// Solidity: function propose(bytes _lookahead, bytes _data) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) Propose(opts *bind.TransactOpts, _lookahead []byte, _data []byte) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "propose", _lookahead, _data)
}

// Propose is a paid mutator transaction binding the contract method 0x9791e644.
//
// Solidity: function propose(bytes _lookahead, bytes _data) returns()
func (_SurgeInboxClient *SurgeInboxClientSession) Propose(_lookahead []byte, _data []byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Propose(&_SurgeInboxClient.TransactOpts, _lookahead, _data)
}

// Propose is a paid mutator transaction binding the contract method 0x9791e644.
//
// Solidity: function propose(bytes _lookahead, bytes _data) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) Propose(_lookahead []byte, _data []byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Propose(&_SurgeInboxClient.TransactOpts, _lookahead, _data)
}

// ProposeWithProof is a paid mutator transaction binding the contract method 0xa2040d21.
//
// Solidity: function proposeWithProof(bytes _lookahead, bytes _data, bytes _proof, bytes32[] _signalSlots) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) ProposeWithProof(opts *bind.TransactOpts, _lookahead []byte, _data []byte, _proof []byte, _signalSlots [][32]byte) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "proposeWithProof", _lookahead, _data, _proof, _signalSlots)
}

// ProposeWithProof is a paid mutator transaction binding the contract method 0xa2040d21.
//
// Solidity: function proposeWithProof(bytes _lookahead, bytes _data, bytes _proof, bytes32[] _signalSlots) returns()
func (_SurgeInboxClient *SurgeInboxClientSession) ProposeWithProof(_lookahead []byte, _data []byte, _proof []byte, _signalSlots [][32]byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.ProposeWithProof(&_SurgeInboxClient.TransactOpts, _lookahead, _data, _proof, _signalSlots)
}

// ProposeWithProof is a paid mutator transaction binding the contract method 0xa2040d21.
//
// Solidity: function proposeWithProof(bytes _lookahead, bytes _data, bytes _proof, bytes32[] _signalSlots) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) ProposeWithProof(_lookahead []byte, _data []byte, _proof []byte, _signalSlots [][32]byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.ProposeWithProof(&_SurgeInboxClient.TransactOpts, _lookahead, _data, _proof, _signalSlots)
}

// Prove is a paid mutator transaction binding the contract method 0xea191743.
//
// Solidity: function prove(bytes _data, bytes _proof) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) Prove(opts *bind.TransactOpts, _data []byte, _proof []byte) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "prove", _data, _proof)
}

// Prove is a paid mutator transaction binding the contract method 0xea191743.
//
// Solidity: function prove(bytes _data, bytes _proof) returns()
func (_SurgeInboxClient *SurgeInboxClientSession) Prove(_data []byte, _proof []byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Prove(&_SurgeInboxClient.TransactOpts, _data, _proof)
}

// Prove is a paid mutator transaction binding the contract method 0xea191743.
//
// Solidity: function prove(bytes _data, bytes _proof) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) Prove(_data []byte, _proof []byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Prove(&_SurgeInboxClient.TransactOpts, _data, _proof)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SurgeInboxClient *SurgeInboxClientSession) RenounceOwnership() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.RenounceOwnership(&_SurgeInboxClient.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.RenounceOwnership(&_SurgeInboxClient.TransactOpts)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0xdbaf2145.
//
// Solidity: function requestWithdrawal() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) RequestWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "requestWithdrawal")
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0xdbaf2145.
//
// Solidity: function requestWithdrawal() returns()
func (_SurgeInboxClient *SurgeInboxClientSession) RequestWithdrawal() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.RequestWithdrawal(&_SurgeInboxClient.TransactOpts)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0xdbaf2145.
//
// Solidity: function requestWithdrawal() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) RequestWithdrawal() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.RequestWithdrawal(&_SurgeInboxClient.TransactOpts)
}

// SaveForcedInclusion is a paid mutator transaction binding the contract method 0xdf596d9e.
//
// Solidity: function saveForcedInclusion((uint16,uint16,uint24) _blobReference) payable returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) SaveForcedInclusion(opts *bind.TransactOpts, _blobReference LibBlobsBlobReference) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "saveForcedInclusion", _blobReference)
}

// SaveForcedInclusion is a paid mutator transaction binding the contract method 0xdf596d9e.
//
// Solidity: function saveForcedInclusion((uint16,uint16,uint24) _blobReference) payable returns()
func (_SurgeInboxClient *SurgeInboxClientSession) SaveForcedInclusion(_blobReference LibBlobsBlobReference) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.SaveForcedInclusion(&_SurgeInboxClient.TransactOpts, _blobReference)
}

// SaveForcedInclusion is a paid mutator transaction binding the contract method 0xdf596d9e.
//
// Solidity: function saveForcedInclusion((uint16,uint16,uint24) _blobReference) payable returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) SaveForcedInclusion(_blobReference LibBlobsBlobReference) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.SaveForcedInclusion(&_SurgeInboxClient.TransactOpts, _blobReference)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SurgeInboxClient *SurgeInboxClientSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.TransferOwnership(&_SurgeInboxClient.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.TransferOwnership(&_SurgeInboxClient.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SurgeInboxClient *SurgeInboxClientSession) Unpause() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Unpause(&_SurgeInboxClient.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) Unpause() (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Unpause(&_SurgeInboxClient.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SurgeInboxClient *SurgeInboxClientSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.UpgradeTo(&_SurgeInboxClient.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.UpgradeTo(&_SurgeInboxClient.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SurgeInboxClient *SurgeInboxClientSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.UpgradeToAndCall(&_SurgeInboxClient.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.UpgradeToAndCall(&_SurgeInboxClient.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd6dad060.
//
// Solidity: function withdraw(address _to, uint64 _amount) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactor) Withdraw(opts *bind.TransactOpts, _to common.Address, _amount uint64) (*types.Transaction, error) {
	return _SurgeInboxClient.contract.Transact(opts, "withdraw", _to, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd6dad060.
//
// Solidity: function withdraw(address _to, uint64 _amount) returns()
func (_SurgeInboxClient *SurgeInboxClientSession) Withdraw(_to common.Address, _amount uint64) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Withdraw(&_SurgeInboxClient.TransactOpts, _to, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd6dad060.
//
// Solidity: function withdraw(address _to, uint64 _amount) returns()
func (_SurgeInboxClient *SurgeInboxClientTransactorSession) Withdraw(_to common.Address, _amount uint64) (*types.Transaction, error) {
	return _SurgeInboxClient.Contract.Withdraw(&_SurgeInboxClient.TransactOpts, _to, _amount)
}

// SurgeInboxClientAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the SurgeInboxClient contract.
type SurgeInboxClientAdminChangedIterator struct {
	Event *SurgeInboxClientAdminChanged // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientAdminChanged)
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
		it.Event = new(SurgeInboxClientAdminChanged)
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
func (it *SurgeInboxClientAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientAdminChanged represents a AdminChanged event raised by the SurgeInboxClient contract.
type SurgeInboxClientAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*SurgeInboxClientAdminChangedIterator, error) {

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientAdminChangedIterator{contract: _SurgeInboxClient.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientAdminChanged) (event.Subscription, error) {

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientAdminChanged)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseAdminChanged(log types.Log) (*SurgeInboxClientAdminChanged, error) {
	event := new(SurgeInboxClientAdminChanged)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the SurgeInboxClient contract.
type SurgeInboxClientBeaconUpgradedIterator struct {
	Event *SurgeInboxClientBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientBeaconUpgraded)
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
		it.Event = new(SurgeInboxClientBeaconUpgraded)
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
func (it *SurgeInboxClientBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientBeaconUpgraded represents a BeaconUpgraded event raised by the SurgeInboxClient contract.
type SurgeInboxClientBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*SurgeInboxClientBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientBeaconUpgradedIterator{contract: _SurgeInboxClient.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientBeaconUpgraded)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseBeaconUpgraded(log types.Log) (*SurgeInboxClientBeaconUpgraded, error) {
	event := new(SurgeInboxClientBeaconUpgraded)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientBondDepositedIterator is returned from FilterBondDeposited and is used to iterate over the raw logs and unpacked data for BondDeposited events raised by the SurgeInboxClient contract.
type SurgeInboxClientBondDepositedIterator struct {
	Event *SurgeInboxClientBondDeposited // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientBondDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientBondDeposited)
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
		it.Event = new(SurgeInboxClientBondDeposited)
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
func (it *SurgeInboxClientBondDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientBondDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientBondDeposited represents a BondDeposited event raised by the SurgeInboxClient contract.
type SurgeInboxClientBondDeposited struct {
	Depositor common.Address
	Recipient common.Address
	Amount    uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBondDeposited is a free log retrieval operation binding the contract event 0xe5e95641fa87bdfef3ce0d39f0c9a37c200f3bf59f53623b3de21e03ed33e3d2.
//
// Solidity: event BondDeposited(address indexed depositor, address indexed recipient, uint64 amount)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterBondDeposited(opts *bind.FilterOpts, depositor []common.Address, recipient []common.Address) (*SurgeInboxClientBondDepositedIterator, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "BondDeposited", depositorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientBondDepositedIterator{contract: _SurgeInboxClient.contract, event: "BondDeposited", logs: logs, sub: sub}, nil
}

// WatchBondDeposited is a free log subscription operation binding the contract event 0xe5e95641fa87bdfef3ce0d39f0c9a37c200f3bf59f53623b3de21e03ed33e3d2.
//
// Solidity: event BondDeposited(address indexed depositor, address indexed recipient, uint64 amount)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchBondDeposited(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientBondDeposited, depositor []common.Address, recipient []common.Address) (event.Subscription, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "BondDeposited", depositorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientBondDeposited)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "BondDeposited", log); err != nil {
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

// ParseBondDeposited is a log parse operation binding the contract event 0xe5e95641fa87bdfef3ce0d39f0c9a37c200f3bf59f53623b3de21e03ed33e3d2.
//
// Solidity: event BondDeposited(address indexed depositor, address indexed recipient, uint64 amount)
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseBondDeposited(log types.Log) (*SurgeInboxClientBondDeposited, error) {
	event := new(SurgeInboxClientBondDeposited)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "BondDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientBondWithdrawnIterator is returned from FilterBondWithdrawn and is used to iterate over the raw logs and unpacked data for BondWithdrawn events raised by the SurgeInboxClient contract.
type SurgeInboxClientBondWithdrawnIterator struct {
	Event *SurgeInboxClientBondWithdrawn // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientBondWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientBondWithdrawn)
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
		it.Event = new(SurgeInboxClientBondWithdrawn)
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
func (it *SurgeInboxClientBondWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientBondWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientBondWithdrawn represents a BondWithdrawn event raised by the SurgeInboxClient contract.
type SurgeInboxClientBondWithdrawn struct {
	Account common.Address
	Amount  uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBondWithdrawn is a free log retrieval operation binding the contract event 0x3362c96009316515fccd3dd29c7036c305ad9e892d83dd5681845ac9edb0c9a8.
//
// Solidity: event BondWithdrawn(address indexed account, uint64 amount)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterBondWithdrawn(opts *bind.FilterOpts, account []common.Address) (*SurgeInboxClientBondWithdrawnIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "BondWithdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientBondWithdrawnIterator{contract: _SurgeInboxClient.contract, event: "BondWithdrawn", logs: logs, sub: sub}, nil
}

// WatchBondWithdrawn is a free log subscription operation binding the contract event 0x3362c96009316515fccd3dd29c7036c305ad9e892d83dd5681845ac9edb0c9a8.
//
// Solidity: event BondWithdrawn(address indexed account, uint64 amount)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchBondWithdrawn(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientBondWithdrawn, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "BondWithdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientBondWithdrawn)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "BondWithdrawn", log); err != nil {
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

// ParseBondWithdrawn is a log parse operation binding the contract event 0x3362c96009316515fccd3dd29c7036c305ad9e892d83dd5681845ac9edb0c9a8.
//
// Solidity: event BondWithdrawn(address indexed account, uint64 amount)
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseBondWithdrawn(log types.Log) (*SurgeInboxClientBondWithdrawn, error) {
	event := new(SurgeInboxClientBondWithdrawn)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "BondWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientForcedInclusionSavedIterator is returned from FilterForcedInclusionSaved and is used to iterate over the raw logs and unpacked data for ForcedInclusionSaved events raised by the SurgeInboxClient contract.
type SurgeInboxClientForcedInclusionSavedIterator struct {
	Event *SurgeInboxClientForcedInclusionSaved // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientForcedInclusionSavedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientForcedInclusionSaved)
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
		it.Event = new(SurgeInboxClientForcedInclusionSaved)
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
func (it *SurgeInboxClientForcedInclusionSavedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientForcedInclusionSavedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientForcedInclusionSaved represents a ForcedInclusionSaved event raised by the SurgeInboxClient contract.
type SurgeInboxClientForcedInclusionSaved struct {
	ForcedInclusion IForcedInclusionStoreForcedInclusion
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterForcedInclusionSaved is a free log retrieval operation binding the contract event 0x18c4fc1e6ac628dbb537b0375bf0efabf1ff2528af1ec22faa74d2da95c29471.
//
// Solidity: event ForcedInclusionSaved((uint64,(bytes32[],uint24,uint48)) forcedInclusion)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterForcedInclusionSaved(opts *bind.FilterOpts) (*SurgeInboxClientForcedInclusionSavedIterator, error) {

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "ForcedInclusionSaved")
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientForcedInclusionSavedIterator{contract: _SurgeInboxClient.contract, event: "ForcedInclusionSaved", logs: logs, sub: sub}, nil
}

// WatchForcedInclusionSaved is a free log subscription operation binding the contract event 0x18c4fc1e6ac628dbb537b0375bf0efabf1ff2528af1ec22faa74d2da95c29471.
//
// Solidity: event ForcedInclusionSaved((uint64,(bytes32[],uint24,uint48)) forcedInclusion)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchForcedInclusionSaved(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientForcedInclusionSaved) (event.Subscription, error) {

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "ForcedInclusionSaved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientForcedInclusionSaved)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "ForcedInclusionSaved", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseForcedInclusionSaved(log types.Log) (*SurgeInboxClientForcedInclusionSaved, error) {
	event := new(SurgeInboxClientForcedInclusionSaved)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "ForcedInclusionSaved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientInboxActivatedIterator is returned from FilterInboxActivated and is used to iterate over the raw logs and unpacked data for InboxActivated events raised by the SurgeInboxClient contract.
type SurgeInboxClientInboxActivatedIterator struct {
	Event *SurgeInboxClientInboxActivated // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientInboxActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientInboxActivated)
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
		it.Event = new(SurgeInboxClientInboxActivated)
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
func (it *SurgeInboxClientInboxActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientInboxActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientInboxActivated represents a InboxActivated event raised by the SurgeInboxClient contract.
type SurgeInboxClientInboxActivated struct {
	LastPacayaBlockHash [32]byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterInboxActivated is a free log retrieval operation binding the contract event 0xe4356761c97932c05c3ee0859fb1a5e4f91f7a1d7a3752c7d5a72d5cc6ecb2d2.
//
// Solidity: event InboxActivated(bytes32 lastPacayaBlockHash)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterInboxActivated(opts *bind.FilterOpts) (*SurgeInboxClientInboxActivatedIterator, error) {

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "InboxActivated")
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientInboxActivatedIterator{contract: _SurgeInboxClient.contract, event: "InboxActivated", logs: logs, sub: sub}, nil
}

// WatchInboxActivated is a free log subscription operation binding the contract event 0xe4356761c97932c05c3ee0859fb1a5e4f91f7a1d7a3752c7d5a72d5cc6ecb2d2.
//
// Solidity: event InboxActivated(bytes32 lastPacayaBlockHash)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchInboxActivated(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientInboxActivated) (event.Subscription, error) {

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "InboxActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientInboxActivated)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "InboxActivated", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseInboxActivated(log types.Log) (*SurgeInboxClientInboxActivated, error) {
	event := new(SurgeInboxClientInboxActivated)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "InboxActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SurgeInboxClient contract.
type SurgeInboxClientInitializedIterator struct {
	Event *SurgeInboxClientInitialized // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientInitialized)
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
		it.Event = new(SurgeInboxClientInitialized)
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
func (it *SurgeInboxClientInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientInitialized represents a Initialized event raised by the SurgeInboxClient contract.
type SurgeInboxClientInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterInitialized(opts *bind.FilterOpts) (*SurgeInboxClientInitializedIterator, error) {

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientInitializedIterator{contract: _SurgeInboxClient.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientInitialized) (event.Subscription, error) {

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientInitialized)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseInitialized(log types.Log) (*SurgeInboxClientInitialized, error) {
	event := new(SurgeInboxClientInitialized)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientLivenessBondSettledIterator is returned from FilterLivenessBondSettled and is used to iterate over the raw logs and unpacked data for LivenessBondSettled events raised by the SurgeInboxClient contract.
type SurgeInboxClientLivenessBondSettledIterator struct {
	Event *SurgeInboxClientLivenessBondSettled // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientLivenessBondSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientLivenessBondSettled)
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
		it.Event = new(SurgeInboxClientLivenessBondSettled)
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
func (it *SurgeInboxClientLivenessBondSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientLivenessBondSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientLivenessBondSettled represents a LivenessBondSettled event raised by the SurgeInboxClient contract.
type SurgeInboxClientLivenessBondSettled struct {
	Payer        common.Address
	Payee        common.Address
	LivenessBond uint64
	Credited     uint64
	Slashed      uint64
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLivenessBondSettled is a free log retrieval operation binding the contract event 0xaa22f5157944b5fa6846460e159d57ea9c3878e71fda274af372fa2ccf285aa0.
//
// Solidity: event LivenessBondSettled(address indexed payer, address indexed payee, uint64 livenessBond, uint64 credited, uint64 slashed)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterLivenessBondSettled(opts *bind.FilterOpts, payer []common.Address, payee []common.Address) (*SurgeInboxClientLivenessBondSettledIterator, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "LivenessBondSettled", payerRule, payeeRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientLivenessBondSettledIterator{contract: _SurgeInboxClient.contract, event: "LivenessBondSettled", logs: logs, sub: sub}, nil
}

// WatchLivenessBondSettled is a free log subscription operation binding the contract event 0xaa22f5157944b5fa6846460e159d57ea9c3878e71fda274af372fa2ccf285aa0.
//
// Solidity: event LivenessBondSettled(address indexed payer, address indexed payee, uint64 livenessBond, uint64 credited, uint64 slashed)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchLivenessBondSettled(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientLivenessBondSettled, payer []common.Address, payee []common.Address) (event.Subscription, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "LivenessBondSettled", payerRule, payeeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientLivenessBondSettled)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "LivenessBondSettled", log); err != nil {
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

// ParseLivenessBondSettled is a log parse operation binding the contract event 0xaa22f5157944b5fa6846460e159d57ea9c3878e71fda274af372fa2ccf285aa0.
//
// Solidity: event LivenessBondSettled(address indexed payer, address indexed payee, uint64 livenessBond, uint64 credited, uint64 slashed)
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseLivenessBondSettled(log types.Log) (*SurgeInboxClientLivenessBondSettled, error) {
	event := new(SurgeInboxClientLivenessBondSettled)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "LivenessBondSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the SurgeInboxClient contract.
type SurgeInboxClientOwnershipTransferStartedIterator struct {
	Event *SurgeInboxClientOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientOwnershipTransferStarted)
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
		it.Event = new(SurgeInboxClientOwnershipTransferStarted)
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
func (it *SurgeInboxClientOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the SurgeInboxClient contract.
type SurgeInboxClientOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SurgeInboxClientOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientOwnershipTransferStartedIterator{contract: _SurgeInboxClient.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientOwnershipTransferStarted)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseOwnershipTransferStarted(log types.Log) (*SurgeInboxClientOwnershipTransferStarted, error) {
	event := new(SurgeInboxClientOwnershipTransferStarted)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SurgeInboxClient contract.
type SurgeInboxClientOwnershipTransferredIterator struct {
	Event *SurgeInboxClientOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientOwnershipTransferred)
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
		it.Event = new(SurgeInboxClientOwnershipTransferred)
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
func (it *SurgeInboxClientOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientOwnershipTransferred represents a OwnershipTransferred event raised by the SurgeInboxClient contract.
type SurgeInboxClientOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SurgeInboxClientOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientOwnershipTransferredIterator{contract: _SurgeInboxClient.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientOwnershipTransferred)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseOwnershipTransferred(log types.Log) (*SurgeInboxClientOwnershipTransferred, error) {
	event := new(SurgeInboxClientOwnershipTransferred)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the SurgeInboxClient contract.
type SurgeInboxClientPausedIterator struct {
	Event *SurgeInboxClientPaused // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientPaused)
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
		it.Event = new(SurgeInboxClientPaused)
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
func (it *SurgeInboxClientPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientPaused represents a Paused event raised by the SurgeInboxClient contract.
type SurgeInboxClientPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterPaused(opts *bind.FilterOpts) (*SurgeInboxClientPausedIterator, error) {

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientPausedIterator{contract: _SurgeInboxClient.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientPaused) (event.Subscription, error) {

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientPaused)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParsePaused(log types.Log) (*SurgeInboxClientPaused, error) {
	event := new(SurgeInboxClientPaused)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientProposedIterator is returned from FilterProposed and is used to iterate over the raw logs and unpacked data for Proposed events raised by the SurgeInboxClient contract.
type SurgeInboxClientProposedIterator struct {
	Event *SurgeInboxClientProposed // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientProposed)
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
		it.Event = new(SurgeInboxClientProposed)
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
func (it *SurgeInboxClientProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientProposed represents a Proposed event raised by the SurgeInboxClient contract.
type SurgeInboxClientProposed struct {
	Id                             *big.Int
	Proposer                       common.Address
	ParentProposalHash             [32]byte
	EndOfSubmissionWindowTimestamp *big.Int
	BasefeeSharingPctg             uint8
	Sources                        []IInboxDerivationSource
	SignalSlotsHash                [32]byte
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterProposed is a free log retrieval operation binding the contract event 0xeb1c821f2b826a6b3ef64e0bd8ac4db546edee46f9d2b53ff08707e5bd7921f0.
//
// Solidity: event Proposed(uint48 indexed id, address indexed proposer, bytes32 parentProposalHash, uint48 endOfSubmissionWindowTimestamp, uint8 basefeeSharingPctg, (bool,(bytes32[],uint24,uint48))[] sources, bytes32 signalSlotsHash)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterProposed(opts *bind.FilterOpts, id []*big.Int, proposer []common.Address) (*SurgeInboxClientProposedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "Proposed", idRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientProposedIterator{contract: _SurgeInboxClient.contract, event: "Proposed", logs: logs, sub: sub}, nil
}

// WatchProposed is a free log subscription operation binding the contract event 0xeb1c821f2b826a6b3ef64e0bd8ac4db546edee46f9d2b53ff08707e5bd7921f0.
//
// Solidity: event Proposed(uint48 indexed id, address indexed proposer, bytes32 parentProposalHash, uint48 endOfSubmissionWindowTimestamp, uint8 basefeeSharingPctg, (bool,(bytes32[],uint24,uint48))[] sources, bytes32 signalSlotsHash)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchProposed(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientProposed, id []*big.Int, proposer []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "Proposed", idRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientProposed)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "Proposed", log); err != nil {
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

// ParseProposed is a log parse operation binding the contract event 0xeb1c821f2b826a6b3ef64e0bd8ac4db546edee46f9d2b53ff08707e5bd7921f0.
//
// Solidity: event Proposed(uint48 indexed id, address indexed proposer, bytes32 parentProposalHash, uint48 endOfSubmissionWindowTimestamp, uint8 basefeeSharingPctg, (bool,(bytes32[],uint24,uint48))[] sources, bytes32 signalSlotsHash)
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseProposed(log types.Log) (*SurgeInboxClientProposed, error) {
	event := new(SurgeInboxClientProposed)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "Proposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientProvedIterator is returned from FilterProved and is used to iterate over the raw logs and unpacked data for Proved events raised by the SurgeInboxClient contract.
type SurgeInboxClientProvedIterator struct {
	Event *SurgeInboxClientProved // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientProvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientProved)
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
		it.Event = new(SurgeInboxClientProved)
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
func (it *SurgeInboxClientProvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientProvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientProved represents a Proved event raised by the SurgeInboxClient contract.
type SurgeInboxClientProved struct {
	FirstProposalId    *big.Int
	FirstNewProposalId *big.Int
	LastProposalId     *big.Int
	ActualProver       common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterProved is a free log retrieval operation binding the contract event 0xa274dcaff3629ec7d69d144038e97732516ff306fcbf8a2bc9423d106779a2f0.
//
// Solidity: event Proved(uint48 firstProposalId, uint48 firstNewProposalId, uint48 lastProposalId, address indexed actualProver)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterProved(opts *bind.FilterOpts, actualProver []common.Address) (*SurgeInboxClientProvedIterator, error) {

	var actualProverRule []interface{}
	for _, actualProverItem := range actualProver {
		actualProverRule = append(actualProverRule, actualProverItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "Proved", actualProverRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientProvedIterator{contract: _SurgeInboxClient.contract, event: "Proved", logs: logs, sub: sub}, nil
}

// WatchProved is a free log subscription operation binding the contract event 0xa274dcaff3629ec7d69d144038e97732516ff306fcbf8a2bc9423d106779a2f0.
//
// Solidity: event Proved(uint48 firstProposalId, uint48 firstNewProposalId, uint48 lastProposalId, address indexed actualProver)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchProved(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientProved, actualProver []common.Address) (event.Subscription, error) {

	var actualProverRule []interface{}
	for _, actualProverItem := range actualProver {
		actualProverRule = append(actualProverRule, actualProverItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "Proved", actualProverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientProved)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "Proved", log); err != nil {
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

// ParseProved is a log parse operation binding the contract event 0xa274dcaff3629ec7d69d144038e97732516ff306fcbf8a2bc9423d106779a2f0.
//
// Solidity: event Proved(uint48 firstProposalId, uint48 firstNewProposalId, uint48 lastProposalId, address indexed actualProver)
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseProved(log types.Log) (*SurgeInboxClientProved, error) {
	event := new(SurgeInboxClientProved)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "Proved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the SurgeInboxClient contract.
type SurgeInboxClientUnpausedIterator struct {
	Event *SurgeInboxClientUnpaused // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientUnpaused)
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
		it.Event = new(SurgeInboxClientUnpaused)
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
func (it *SurgeInboxClientUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientUnpaused represents a Unpaused event raised by the SurgeInboxClient contract.
type SurgeInboxClientUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterUnpaused(opts *bind.FilterOpts) (*SurgeInboxClientUnpausedIterator, error) {

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientUnpausedIterator{contract: _SurgeInboxClient.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientUnpaused) (event.Subscription, error) {

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientUnpaused)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseUnpaused(log types.Log) (*SurgeInboxClientUnpaused, error) {
	event := new(SurgeInboxClientUnpaused)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the SurgeInboxClient contract.
type SurgeInboxClientUpgradedIterator struct {
	Event *SurgeInboxClientUpgraded // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientUpgraded)
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
		it.Event = new(SurgeInboxClientUpgraded)
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
func (it *SurgeInboxClientUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientUpgraded represents a Upgraded event raised by the SurgeInboxClient contract.
type SurgeInboxClientUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SurgeInboxClientUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientUpgradedIterator{contract: _SurgeInboxClient.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientUpgraded)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseUpgraded(log types.Log) (*SurgeInboxClientUpgraded, error) {
	event := new(SurgeInboxClientUpgraded)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientWithdrawalCancelledIterator is returned from FilterWithdrawalCancelled and is used to iterate over the raw logs and unpacked data for WithdrawalCancelled events raised by the SurgeInboxClient contract.
type SurgeInboxClientWithdrawalCancelledIterator struct {
	Event *SurgeInboxClientWithdrawalCancelled // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientWithdrawalCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientWithdrawalCancelled)
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
		it.Event = new(SurgeInboxClientWithdrawalCancelled)
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
func (it *SurgeInboxClientWithdrawalCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientWithdrawalCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientWithdrawalCancelled represents a WithdrawalCancelled event raised by the SurgeInboxClient contract.
type SurgeInboxClientWithdrawalCancelled struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalCancelled is a free log retrieval operation binding the contract event 0xc51fdb96728de385ec7859819e3997bc618362ef0dbca0ad051d856866cda3db.
//
// Solidity: event WithdrawalCancelled(address indexed account)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterWithdrawalCancelled(opts *bind.FilterOpts, account []common.Address) (*SurgeInboxClientWithdrawalCancelledIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "WithdrawalCancelled", accountRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientWithdrawalCancelledIterator{contract: _SurgeInboxClient.contract, event: "WithdrawalCancelled", logs: logs, sub: sub}, nil
}

// WatchWithdrawalCancelled is a free log subscription operation binding the contract event 0xc51fdb96728de385ec7859819e3997bc618362ef0dbca0ad051d856866cda3db.
//
// Solidity: event WithdrawalCancelled(address indexed account)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchWithdrawalCancelled(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientWithdrawalCancelled, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "WithdrawalCancelled", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientWithdrawalCancelled)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "WithdrawalCancelled", log); err != nil {
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

// ParseWithdrawalCancelled is a log parse operation binding the contract event 0xc51fdb96728de385ec7859819e3997bc618362ef0dbca0ad051d856866cda3db.
//
// Solidity: event WithdrawalCancelled(address indexed account)
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseWithdrawalCancelled(log types.Log) (*SurgeInboxClientWithdrawalCancelled, error) {
	event := new(SurgeInboxClientWithdrawalCancelled)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "WithdrawalCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeInboxClientWithdrawalRequestedIterator is returned from FilterWithdrawalRequested and is used to iterate over the raw logs and unpacked data for WithdrawalRequested events raised by the SurgeInboxClient contract.
type SurgeInboxClientWithdrawalRequestedIterator struct {
	Event *SurgeInboxClientWithdrawalRequested // Event containing the contract specifics and raw log

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
func (it *SurgeInboxClientWithdrawalRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeInboxClientWithdrawalRequested)
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
		it.Event = new(SurgeInboxClientWithdrawalRequested)
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
func (it *SurgeInboxClientWithdrawalRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeInboxClientWithdrawalRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeInboxClientWithdrawalRequested represents a WithdrawalRequested event raised by the SurgeInboxClient contract.
type SurgeInboxClientWithdrawalRequested struct {
	Account        common.Address
	WithdrawableAt *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalRequested is a free log retrieval operation binding the contract event 0x3bbe41cfdd142e0f9b2224dac18c6efd2a6966e35a9ec23ab57ce63a60b33604.
//
// Solidity: event WithdrawalRequested(address indexed account, uint48 withdrawableAt)
func (_SurgeInboxClient *SurgeInboxClientFilterer) FilterWithdrawalRequested(opts *bind.FilterOpts, account []common.Address) (*SurgeInboxClientWithdrawalRequestedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.FilterLogs(opts, "WithdrawalRequested", accountRule)
	if err != nil {
		return nil, err
	}
	return &SurgeInboxClientWithdrawalRequestedIterator{contract: _SurgeInboxClient.contract, event: "WithdrawalRequested", logs: logs, sub: sub}, nil
}

// WatchWithdrawalRequested is a free log subscription operation binding the contract event 0x3bbe41cfdd142e0f9b2224dac18c6efd2a6966e35a9ec23ab57ce63a60b33604.
//
// Solidity: event WithdrawalRequested(address indexed account, uint48 withdrawableAt)
func (_SurgeInboxClient *SurgeInboxClientFilterer) WatchWithdrawalRequested(opts *bind.WatchOpts, sink chan<- *SurgeInboxClientWithdrawalRequested, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SurgeInboxClient.contract.WatchLogs(opts, "WithdrawalRequested", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeInboxClientWithdrawalRequested)
				if err := _SurgeInboxClient.contract.UnpackLog(event, "WithdrawalRequested", log); err != nil {
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

// ParseWithdrawalRequested is a log parse operation binding the contract event 0x3bbe41cfdd142e0f9b2224dac18c6efd2a6966e35a9ec23ab57ce63a60b33604.
//
// Solidity: event WithdrawalRequested(address indexed account, uint48 withdrawableAt)
func (_SurgeInboxClient *SurgeInboxClientFilterer) ParseWithdrawalRequested(log types.Log) (*SurgeInboxClientWithdrawalRequested, error) {
	event := new(SurgeInboxClientWithdrawalRequested)
	if err := _SurgeInboxClient.contract.UnpackLog(event, "WithdrawalRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
