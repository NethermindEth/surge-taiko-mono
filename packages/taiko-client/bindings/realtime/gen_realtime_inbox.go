// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package realtime

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

// ICheckpointStoreCheckpoint is an auto generated low-level Go binding around an user-defined struct.
type ICheckpointStoreCheckpoint struct {
	BlockNumber *big.Int
	BlockHash   [32]byte
	StateRoot   [32]byte
}

// IInboxDerivationSource is an auto generated low-level Go binding around an user-defined struct.
type IInboxDerivationSource struct {
	IsForcedInclusion bool
	BlobSlice         LibBlobsBlobSlice
}

// IRealTimeInboxConfig is an auto generated low-level Go binding around an user-defined struct.
type IRealTimeInboxConfig struct {
	ProofVerifier      common.Address
	SignalService      common.Address
	BasefeeSharingPctg uint8
}

// LibBlobsBlobSlice is an auto generated low-level Go binding around an user-defined struct.
type LibBlobsBlobSlice struct {
	BlobHashes [][32]byte
	Offset     *big.Int
	Timestamp  *big.Int
}

// RealTimeInboxClientMetaData contains all meta data concerning the RealTimeInboxClient contract.
var RealTimeInboxClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_config\",\"type\":\"tuple\",\"internalType\":\"struct IRealTimeInbox.Config\",\"components\":[{\"name\":\"proofVerifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signalService\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"activate\",\"inputs\":[{\"name\":\"_genesisProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decodeProposeInput\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"input_\",\"type\":\"tuple\",\"internalType\":\"struct IRealTimeInbox.ProposeInput\",\"components\":[{\"name\":\"blobReference\",\"type\":\"tuple\",\"internalType\":\"struct LibBlobs.BlobReference\",\"components\":[{\"name\":\"blobStartIndex\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"numBlobs\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]},{\"name\":\"signalSlots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"maxAnchorBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeProposeInput\",\"inputs\":[{\"name\":\"_input\",\"type\":\"tuple\",\"internalType\":\"struct IRealTimeInbox.ProposeInput\",\"components\":[{\"name\":\"blobReference\",\"type\":\"tuple\",\"internalType\":\"struct LibBlobs.BlobReference\",\"components\":[{\"name\":\"blobStartIndex\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"numBlobs\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]},{\"name\":\"signalSlots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"maxAnchorBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}],\"outputs\":[{\"name\":\"encoded_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"config_\",\"type\":\"tuple\",\"internalType\":\"struct IRealTimeInbox.Config\",\"components\":[{\"name\":\"proofVerifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signalService\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLastProposalHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashCommitment\",\"inputs\":[{\"name\":\"_commitment\",\"type\":\"tuple\",\"internalType\":\"struct IRealTimeInbox.Commitment\",\"components\":[{\"name\":\"proposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"checkpoint\",\"type\":\"tuple\",\"internalType\":\"struct ICheckpointStore.Checkpoint\",\"components\":[{\"name\":\"blockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"hashProposal\",\"inputs\":[{\"name\":\"_proposal\",\"type\":\"tuple\",\"internalType\":\"struct IRealTimeInbox.Proposal\",\"components\":[{\"name\":\"parentProposalHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maxAnchorBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"maxAnchorBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"sources\",\"type\":\"tuple[]\",\"internalType\":\"struct IInbox.DerivationSource[]\",\"components\":[{\"name\":\"isForcedInclusion\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"struct LibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]},{\"name\":\"signalSlotsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"hashSignalSlots\",\"inputs\":[{\"name\":\"_signalSlots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"impl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inNonReentrant\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"init\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastProposalHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"propose\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_checkpoint\",\"type\":\"tuple\",\"internalType\":\"struct ICheckpointStore.Checkpoint\",\"components\":[{\"name\":\"blockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"_proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Activated\",\"inputs\":[{\"name\":\"genesisProposalHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AdminChanged\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BeaconUpgraded\",\"inputs\":[{\"name\":\"beacon\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProposedAndProved\",\"inputs\":[{\"name\":\"proposalHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"parentProposalHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"maxAnchorBlockNumber\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"basefeeSharingPctg\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"sources\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"struct IInbox.DerivationSource[]\",\"components\":[{\"name\":\"isForcedInclusion\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blobSlice\",\"type\":\"tuple\",\"internalType\":\"struct LibBlobs.BlobSlice\",\"components\":[{\"name\":\"blobHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"offset\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}]},{\"name\":\"signalSlots\",\"type\":\"bytes32[]\",\"indexed\":false,\"internalType\":\"bytes32[]\"},{\"name\":\"checkpoint\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"struct ICheckpointStore.Checkpoint\",\"components\":[{\"name\":\"blockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ACCESS_DENIED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyActivated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BlobNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FUNC_NOT_IMPLEMENTED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"INVALID_PAUSE_STATUS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidGenesisHash\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxAnchorBlockTooOld\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoBlobs\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotActivated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"REENTRANT_CALL\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignalSlotNotSent\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZERO_ADDRESS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_VALUE\",\"inputs\":[]}]",
}

// RealTimeInboxClientABI is the input ABI used to generate the binding from.
// Deprecated: Use RealTimeInboxClientMetaData.ABI instead.
var RealTimeInboxClientABI = RealTimeInboxClientMetaData.ABI

// RealTimeInboxClient is an auto generated Go binding around an Ethereum contract.
type RealTimeInboxClient struct {
	RealTimeInboxClientCaller     // Read-only binding to the contract
	RealTimeInboxClientTransactor // Write-only binding to the contract
	RealTimeInboxClientFilterer   // Log filterer for contract events
}

// RealTimeInboxClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type RealTimeInboxClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RealTimeInboxClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RealTimeInboxClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RealTimeInboxClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RealTimeInboxClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RealTimeInboxClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RealTimeInboxClientSession struct {
	Contract     *RealTimeInboxClient // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RealTimeInboxClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RealTimeInboxClientCallerSession struct {
	Contract *RealTimeInboxClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// RealTimeInboxClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RealTimeInboxClientTransactorSession struct {
	Contract     *RealTimeInboxClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// RealTimeInboxClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type RealTimeInboxClientRaw struct {
	Contract *RealTimeInboxClient // Generic contract binding to access the raw methods on
}

// RealTimeInboxClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RealTimeInboxClientCallerRaw struct {
	Contract *RealTimeInboxClientCaller // Generic read-only contract binding to access the raw methods on
}

// RealTimeInboxClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RealTimeInboxClientTransactorRaw struct {
	Contract *RealTimeInboxClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRealTimeInboxClient creates a new instance of RealTimeInboxClient, bound to a specific deployed contract.
func NewRealTimeInboxClient(address common.Address, backend bind.ContractBackend) (*RealTimeInboxClient, error) {
	contract, err := bindRealTimeInboxClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RealTimeInboxClient{RealTimeInboxClientCaller: RealTimeInboxClientCaller{contract: contract}, RealTimeInboxClientTransactor: RealTimeInboxClientTransactor{contract: contract}, RealTimeInboxClientFilterer: RealTimeInboxClientFilterer{contract: contract}}, nil
}

// NewRealTimeInboxClientCaller creates a new read-only instance of RealTimeInboxClient, bound to a specific deployed contract.
func NewRealTimeInboxClientCaller(address common.Address, caller bind.ContractCaller) (*RealTimeInboxClientCaller, error) {
	contract, err := bindRealTimeInboxClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RealTimeInboxClientCaller{contract: contract}, nil
}

// NewRealTimeInboxClientTransactor creates a new write-only instance of RealTimeInboxClient, bound to a specific deployed contract.
func NewRealTimeInboxClientTransactor(address common.Address, transactor bind.ContractTransactor) (*RealTimeInboxClientTransactor, error) {
	contract, err := bindRealTimeInboxClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RealTimeInboxClientTransactor{contract: contract}, nil
}

// NewRealTimeInboxClientFilterer creates a new log filterer instance of RealTimeInboxClient, bound to a specific deployed contract.
func NewRealTimeInboxClientFilterer(address common.Address, filterer bind.ContractFilterer) (*RealTimeInboxClientFilterer, error) {
	contract, err := bindRealTimeInboxClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RealTimeInboxClientFilterer{contract: contract}, nil
}

// bindRealTimeInboxClient binds a generic wrapper to an already deployed contract.
func bindRealTimeInboxClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RealTimeInboxClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RealTimeInboxClient *RealTimeInboxClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RealTimeInboxClient.Contract.RealTimeInboxClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RealTimeInboxClient *RealTimeInboxClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RealTimeInboxClient.Contract.RealTimeInboxClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RealTimeInboxClient *RealTimeInboxClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RealTimeInboxClient.Contract.RealTimeInboxClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RealTimeInboxClient *RealTimeInboxClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RealTimeInboxClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RealTimeInboxClient *RealTimeInboxClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RealTimeInboxClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RealTimeInboxClient *RealTimeInboxClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RealTimeInboxClient.Contract.contract.Transact(opts, method, params...)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,uint8) config_)
func (_RealTimeInboxClient *RealTimeInboxClientCaller) GetConfig(opts *bind.CallOpts) (IRealTimeInboxConfig, error) {
	var out []interface{}
	err := _RealTimeInboxClient.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(IRealTimeInboxConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IRealTimeInboxConfig)).(*IRealTimeInboxConfig)

	return out0, err
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,uint8) config_)
func (_RealTimeInboxClient *RealTimeInboxClientSession) GetConfig() (IRealTimeInboxConfig, error) {
	return _RealTimeInboxClient.Contract.GetConfig(&_RealTimeInboxClient.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,uint8) config_)
func (_RealTimeInboxClient *RealTimeInboxClientCallerSession) GetConfig() (IRealTimeInboxConfig, error) {
	return _RealTimeInboxClient.Contract.GetConfig(&_RealTimeInboxClient.CallOpts)
}

// GetLastProposalHash is a free data retrieval call binding the contract method 0xabc0be3a.
//
// Solidity: function getLastProposalHash() view returns(bytes32)
func (_RealTimeInboxClient *RealTimeInboxClientCaller) GetLastProposalHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RealTimeInboxClient.contract.Call(opts, &out, "getLastProposalHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// GetLastProposalHash is a free data retrieval call binding the contract method 0xabc0be3a.
//
// Solidity: function getLastProposalHash() view returns(bytes32)
func (_RealTimeInboxClient *RealTimeInboxClientSession) GetLastProposalHash() ([32]byte, error) {
	return _RealTimeInboxClient.Contract.GetLastProposalHash(&_RealTimeInboxClient.CallOpts)
}

// GetLastProposalHash is a free data retrieval call binding the contract method 0xabc0be3a.
//
// Solidity: function getLastProposalHash() view returns(bytes32)
func (_RealTimeInboxClient *RealTimeInboxClientCallerSession) GetLastProposalHash() ([32]byte, error) {
	return _RealTimeInboxClient.Contract.GetLastProposalHash(&_RealTimeInboxClient.CallOpts)
}

// RealTimeInboxClientProposedAndProvedIterator is returned from FilterProposedAndProved and is used to iterate over the raw logs and unpacked data for ProposedAndProved events raised by the RealTimeInboxClient contract.
type RealTimeInboxClientProposedAndProvedIterator struct {
	Event *RealTimeInboxClientProposedAndProved // Event containing the contract specifics and raw log

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
func (it *RealTimeInboxClientProposedAndProvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RealTimeInboxClientProposedAndProved)
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
		it.Event = new(RealTimeInboxClientProposedAndProved)
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
func (it *RealTimeInboxClientProposedAndProvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RealTimeInboxClientProposedAndProvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RealTimeInboxClientProposedAndProved represents a ProposedAndProved event raised by the RealTimeInboxClient contract.
type RealTimeInboxClientProposedAndProved struct {
	ProposalHash       [32]byte
	ParentProposalHash [32]byte
	MaxAnchorBlockNumber *big.Int
	BasefeeSharingPctg uint8
	Sources            []IInboxDerivationSource
	SignalSlots        [][32]byte
	Checkpoint         ICheckpointStoreCheckpoint
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterProposedAndProved is a free log retrieval operation binding the contract event.
//
// Solidity: event ProposedAndProved(bytes32 indexed proposalHash, bytes32 parentProposalHash, uint48 maxAnchorBlockNumber, uint8 basefeeSharingPctg, (bool,(bytes32[],uint24,uint48))[] sources, bytes32[] signalSlots, (uint48,bytes32,bytes32) checkpoint)
func (_RealTimeInboxClient *RealTimeInboxClientFilterer) FilterProposedAndProved(opts *bind.FilterOpts, proposalHash [][32]byte) (*RealTimeInboxClientProposedAndProvedIterator, error) {

	var proposalHashRule []interface{}
	for _, proposalHashItem := range proposalHash {
		proposalHashRule = append(proposalHashRule, proposalHashItem)
	}

	logs, sub, err := _RealTimeInboxClient.contract.FilterLogs(opts, "ProposedAndProved", proposalHashRule)
	if err != nil {
		return nil, err
	}
	return &RealTimeInboxClientProposedAndProvedIterator{contract: _RealTimeInboxClient.contract, event: "ProposedAndProved", logs: logs, sub: sub}, nil
}

// WatchProposedAndProved is a free log subscription operation binding the contract event.
//
// Solidity: event ProposedAndProved(bytes32 indexed proposalHash, bytes32 parentProposalHash, uint48 maxAnchorBlockNumber, uint8 basefeeSharingPctg, (bool,(bytes32[],uint24,uint48))[] sources, bytes32[] signalSlots, (uint48,bytes32,bytes32) checkpoint)
func (_RealTimeInboxClient *RealTimeInboxClientFilterer) WatchProposedAndProved(opts *bind.WatchOpts, sink chan<- *RealTimeInboxClientProposedAndProved, proposalHash [][32]byte) (event.Subscription, error) {

	var proposalHashRule []interface{}
	for _, proposalHashItem := range proposalHash {
		proposalHashRule = append(proposalHashRule, proposalHashItem)
	}

	logs, sub, err := _RealTimeInboxClient.contract.WatchLogs(opts, "ProposedAndProved", proposalHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RealTimeInboxClientProposedAndProved)
				if err := _RealTimeInboxClient.contract.UnpackLog(event, "ProposedAndProved", log); err != nil {
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

// ParseProposedAndProved is a log parse operation binding the contract event.
//
// Solidity: event ProposedAndProved(bytes32 indexed proposalHash, bytes32 parentProposalHash, uint48 maxAnchorBlockNumber, uint8 basefeeSharingPctg, (bool,(bytes32[],uint24,uint48))[] sources, bytes32[] signalSlots, (uint48,bytes32,bytes32) checkpoint)
func (_RealTimeInboxClient *RealTimeInboxClientFilterer) ParseProposedAndProved(log types.Log) (*RealTimeInboxClientProposedAndProved, error) {
	event := new(RealTimeInboxClientProposedAndProved)
	if err := _RealTimeInboxClient.contract.UnpackLog(event, "ProposedAndProved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
