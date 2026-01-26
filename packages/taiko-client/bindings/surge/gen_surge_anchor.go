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

// AnchorBlockParams is an auto generated low-level Go binding around an user-defined struct.
type AnchorBlockParams struct {
	AnchorBlockNumber *big.Int
	AnchorBlockHash   [32]byte
	AnchorStateRoot   [32]byte
	RawTxListHash     [32]byte
}

// AnchorBlockState is an auto generated low-level Go binding around an user-defined struct.
type AnchorBlockState struct {
	AnchorBlockNumber *big.Int
	AncestorsHash     [32]byte
}

// AnchorPreconfMetadata is an auto generated low-level Go binding around an user-defined struct.
type AnchorPreconfMetadata struct {
	AnchorBlockNumber         *big.Int
	SubmissionWindowEnd       *big.Int
	ParentSubmissionWindowEnd *big.Int
	RawTxListHash             [32]byte
	ParentRawTxListHash       [32]byte
}

// AnchorProposalParams is an auto generated low-level Go binding around an user-defined struct.
type AnchorProposalParams struct {
	SubmissionWindowEnd *big.Int
}

// ICheckpointStoreCheckpoint is an auto generated low-level Go binding around an user-defined struct.
type ICheckpointStoreCheckpoint struct {
	BlockNumber *big.Int
	BlockHash   [32]byte
	StateRoot   [32]byte
}

// SurgeAnchorMetaData contains all meta data concerning the SurgeAnchor contract.
var SurgeAnchorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_checkpointStore\",\"type\":\"address\",\"internalType\":\"contractICheckpointStore\"},{\"name\":\"_l1ChainId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ANCHOR_GAS_LIMIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"GOLDEN_TOUCH_ADDRESS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchorV4\",\"inputs\":[{\"name\":\"_checkpoint\",\"type\":\"tuple\",\"internalType\":\"structICheckpointStore.Checkpoint\",\"components\":[{\"name\":\"blockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchorV4WithSignalSlots\",\"inputs\":[{\"name\":\"_checkpoint\",\"type\":\"tuple\",\"internalType\":\"structICheckpointStore.Checkpoint\",\"components\":[{\"name\":\"blockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"_signalSlots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchorV5\",\"inputs\":[{\"name\":\"_proposalParams\",\"type\":\"tuple\",\"internalType\":\"structAnchor.ProposalParams\",\"components\":[{\"name\":\"submissionWindowEnd\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"name\":\"_blockParams\",\"type\":\"tuple\",\"internalType\":\"structAnchor.BlockParams\",\"components\":[{\"name\":\"anchorBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"anchorBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"anchorStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rawTxListHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blockHashes\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkpointStore\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractICheckpointStore\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlockState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAnchor.BlockState\",\"components\":[{\"name\":\"anchorBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"ancestorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPreconfMetadata\",\"inputs\":[{\"name\":\"_blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAnchor.PreconfMetadata\",\"components\":[{\"name\":\"anchorBlockNumber\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"submissionWindowEnd\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"parentSubmissionWindowEnd\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"rawTxListHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"parentRawTxListHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"impl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inNonReentrant\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"init\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"l1ChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminChanged\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Anchored\",\"inputs\":[{\"name\":\"prevAnchorBlockNumber\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"anchorBlockNumber\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"ancestorsHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BeaconUpgraded\",\"inputs\":[{\"name\":\"beacon\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ACCESS_DENIED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AncestorsHashMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ETH_TRANSFER_FAILED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FUNC_NOT_IMPLEMENTED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"INVALID_PAUSE_STATUS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBlockNumber\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidL1ChainId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidL2ChainId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSender\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"REENTRANT_CALL\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_ADDRESS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_VALUE\",\"inputs\":[]}]",
}

// SurgeAnchorABI is the input ABI used to generate the binding from.
// Deprecated: Use SurgeAnchorMetaData.ABI instead.
var SurgeAnchorABI = SurgeAnchorMetaData.ABI

// SurgeAnchor is an auto generated Go binding around an Ethereum contract.
type SurgeAnchor struct {
	SurgeAnchorCaller     // Read-only binding to the contract
	SurgeAnchorTransactor // Write-only binding to the contract
	SurgeAnchorFilterer   // Log filterer for contract events
}

// SurgeAnchorCaller is an auto generated read-only Go binding around an Ethereum contract.
type SurgeAnchorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeAnchorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SurgeAnchorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeAnchorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SurgeAnchorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeAnchorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SurgeAnchorSession struct {
	Contract     *SurgeAnchor      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SurgeAnchorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SurgeAnchorCallerSession struct {
	Contract *SurgeAnchorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SurgeAnchorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SurgeAnchorTransactorSession struct {
	Contract     *SurgeAnchorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SurgeAnchorRaw is an auto generated low-level Go binding around an Ethereum contract.
type SurgeAnchorRaw struct {
	Contract *SurgeAnchor // Generic contract binding to access the raw methods on
}

// SurgeAnchorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SurgeAnchorCallerRaw struct {
	Contract *SurgeAnchorCaller // Generic read-only contract binding to access the raw methods on
}

// SurgeAnchorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SurgeAnchorTransactorRaw struct {
	Contract *SurgeAnchorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSurgeAnchor creates a new instance of SurgeAnchor, bound to a specific deployed contract.
func NewSurgeAnchor(address common.Address, backend bind.ContractBackend) (*SurgeAnchor, error) {
	contract, err := bindSurgeAnchor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SurgeAnchor{SurgeAnchorCaller: SurgeAnchorCaller{contract: contract}, SurgeAnchorTransactor: SurgeAnchorTransactor{contract: contract}, SurgeAnchorFilterer: SurgeAnchorFilterer{contract: contract}}, nil
}

// NewSurgeAnchorCaller creates a new read-only instance of SurgeAnchor, bound to a specific deployed contract.
func NewSurgeAnchorCaller(address common.Address, caller bind.ContractCaller) (*SurgeAnchorCaller, error) {
	contract, err := bindSurgeAnchor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorCaller{contract: contract}, nil
}

// NewSurgeAnchorTransactor creates a new write-only instance of SurgeAnchor, bound to a specific deployed contract.
func NewSurgeAnchorTransactor(address common.Address, transactor bind.ContractTransactor) (*SurgeAnchorTransactor, error) {
	contract, err := bindSurgeAnchor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorTransactor{contract: contract}, nil
}

// NewSurgeAnchorFilterer creates a new log filterer instance of SurgeAnchor, bound to a specific deployed contract.
func NewSurgeAnchorFilterer(address common.Address, filterer bind.ContractFilterer) (*SurgeAnchorFilterer, error) {
	contract, err := bindSurgeAnchor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorFilterer{contract: contract}, nil
}

// bindSurgeAnchor binds a generic wrapper to an already deployed contract.
func bindSurgeAnchor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SurgeAnchorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SurgeAnchor *SurgeAnchorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SurgeAnchor.Contract.SurgeAnchorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SurgeAnchor *SurgeAnchorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.SurgeAnchorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SurgeAnchor *SurgeAnchorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.SurgeAnchorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SurgeAnchor *SurgeAnchorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SurgeAnchor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SurgeAnchor *SurgeAnchorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SurgeAnchor *SurgeAnchorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.contract.Transact(opts, method, params...)
}

// ANCHORGASLIMIT is a free data retrieval call binding the contract method 0xc46e3a66.
//
// Solidity: function ANCHOR_GAS_LIMIT() view returns(uint64)
func (_SurgeAnchor *SurgeAnchorCaller) ANCHORGASLIMIT(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "ANCHOR_GAS_LIMIT")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ANCHORGASLIMIT is a free data retrieval call binding the contract method 0xc46e3a66.
//
// Solidity: function ANCHOR_GAS_LIMIT() view returns(uint64)
func (_SurgeAnchor *SurgeAnchorSession) ANCHORGASLIMIT() (uint64, error) {
	return _SurgeAnchor.Contract.ANCHORGASLIMIT(&_SurgeAnchor.CallOpts)
}

// ANCHORGASLIMIT is a free data retrieval call binding the contract method 0xc46e3a66.
//
// Solidity: function ANCHOR_GAS_LIMIT() view returns(uint64)
func (_SurgeAnchor *SurgeAnchorCallerSession) ANCHORGASLIMIT() (uint64, error) {
	return _SurgeAnchor.Contract.ANCHORGASLIMIT(&_SurgeAnchor.CallOpts)
}

// GOLDENTOUCHADDRESS is a free data retrieval call binding the contract method 0x9ee512f2.
//
// Solidity: function GOLDEN_TOUCH_ADDRESS() view returns(address)
func (_SurgeAnchor *SurgeAnchorCaller) GOLDENTOUCHADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "GOLDEN_TOUCH_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GOLDENTOUCHADDRESS is a free data retrieval call binding the contract method 0x9ee512f2.
//
// Solidity: function GOLDEN_TOUCH_ADDRESS() view returns(address)
func (_SurgeAnchor *SurgeAnchorSession) GOLDENTOUCHADDRESS() (common.Address, error) {
	return _SurgeAnchor.Contract.GOLDENTOUCHADDRESS(&_SurgeAnchor.CallOpts)
}

// GOLDENTOUCHADDRESS is a free data retrieval call binding the contract method 0x9ee512f2.
//
// Solidity: function GOLDEN_TOUCH_ADDRESS() view returns(address)
func (_SurgeAnchor *SurgeAnchorCallerSession) GOLDENTOUCHADDRESS() (common.Address, error) {
	return _SurgeAnchor.Contract.GOLDENTOUCHADDRESS(&_SurgeAnchor.CallOpts)
}

// BlockHashes is a free data retrieval call binding the contract method 0x34cdf78d.
//
// Solidity: function blockHashes(uint256 blockNumber) view returns(bytes32 blockHash)
func (_SurgeAnchor *SurgeAnchorCaller) BlockHashes(opts *bind.CallOpts, blockNumber *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "blockHashes", blockNumber)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockHashes is a free data retrieval call binding the contract method 0x34cdf78d.
//
// Solidity: function blockHashes(uint256 blockNumber) view returns(bytes32 blockHash)
func (_SurgeAnchor *SurgeAnchorSession) BlockHashes(blockNumber *big.Int) ([32]byte, error) {
	return _SurgeAnchor.Contract.BlockHashes(&_SurgeAnchor.CallOpts, blockNumber)
}

// BlockHashes is a free data retrieval call binding the contract method 0x34cdf78d.
//
// Solidity: function blockHashes(uint256 blockNumber) view returns(bytes32 blockHash)
func (_SurgeAnchor *SurgeAnchorCallerSession) BlockHashes(blockNumber *big.Int) ([32]byte, error) {
	return _SurgeAnchor.Contract.BlockHashes(&_SurgeAnchor.CallOpts, blockNumber)
}

// CheckpointStore is a free data retrieval call binding the contract method 0x955a7244.
//
// Solidity: function checkpointStore() view returns(address)
func (_SurgeAnchor *SurgeAnchorCaller) CheckpointStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "checkpointStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CheckpointStore is a free data retrieval call binding the contract method 0x955a7244.
//
// Solidity: function checkpointStore() view returns(address)
func (_SurgeAnchor *SurgeAnchorSession) CheckpointStore() (common.Address, error) {
	return _SurgeAnchor.Contract.CheckpointStore(&_SurgeAnchor.CallOpts)
}

// CheckpointStore is a free data retrieval call binding the contract method 0x955a7244.
//
// Solidity: function checkpointStore() view returns(address)
func (_SurgeAnchor *SurgeAnchorCallerSession) CheckpointStore() (common.Address, error) {
	return _SurgeAnchor.Contract.CheckpointStore(&_SurgeAnchor.CallOpts)
}

// GetBlockState is a free data retrieval call binding the contract method 0x0f439bd9.
//
// Solidity: function getBlockState() view returns((uint48,bytes32))
func (_SurgeAnchor *SurgeAnchorCaller) GetBlockState(opts *bind.CallOpts) (AnchorBlockState, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "getBlockState")

	if err != nil {
		return *new(AnchorBlockState), err
	}

	out0 := *abi.ConvertType(out[0], new(AnchorBlockState)).(*AnchorBlockState)

	return out0, err

}

// GetBlockState is a free data retrieval call binding the contract method 0x0f439bd9.
//
// Solidity: function getBlockState() view returns((uint48,bytes32))
func (_SurgeAnchor *SurgeAnchorSession) GetBlockState() (AnchorBlockState, error) {
	return _SurgeAnchor.Contract.GetBlockState(&_SurgeAnchor.CallOpts)
}

// GetBlockState is a free data retrieval call binding the contract method 0x0f439bd9.
//
// Solidity: function getBlockState() view returns((uint48,bytes32))
func (_SurgeAnchor *SurgeAnchorCallerSession) GetBlockState() (AnchorBlockState, error) {
	return _SurgeAnchor.Contract.GetBlockState(&_SurgeAnchor.CallOpts)
}

// GetPreconfMetadata is a free data retrieval call binding the contract method 0x260c5344.
//
// Solidity: function getPreconfMetadata(uint256 _blockNumber) view returns((uint48,uint48,uint48,bytes32,bytes32))
func (_SurgeAnchor *SurgeAnchorCaller) GetPreconfMetadata(opts *bind.CallOpts, _blockNumber *big.Int) (AnchorPreconfMetadata, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "getPreconfMetadata", _blockNumber)

	if err != nil {
		return *new(AnchorPreconfMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new(AnchorPreconfMetadata)).(*AnchorPreconfMetadata)

	return out0, err

}

// GetPreconfMetadata is a free data retrieval call binding the contract method 0x260c5344.
//
// Solidity: function getPreconfMetadata(uint256 _blockNumber) view returns((uint48,uint48,uint48,bytes32,bytes32))
func (_SurgeAnchor *SurgeAnchorSession) GetPreconfMetadata(_blockNumber *big.Int) (AnchorPreconfMetadata, error) {
	return _SurgeAnchor.Contract.GetPreconfMetadata(&_SurgeAnchor.CallOpts, _blockNumber)
}

// GetPreconfMetadata is a free data retrieval call binding the contract method 0x260c5344.
//
// Solidity: function getPreconfMetadata(uint256 _blockNumber) view returns((uint48,uint48,uint48,bytes32,bytes32))
func (_SurgeAnchor *SurgeAnchorCallerSession) GetPreconfMetadata(_blockNumber *big.Int) (AnchorPreconfMetadata, error) {
	return _SurgeAnchor.Contract.GetPreconfMetadata(&_SurgeAnchor.CallOpts, _blockNumber)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_SurgeAnchor *SurgeAnchorCaller) Impl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "impl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_SurgeAnchor *SurgeAnchorSession) Impl() (common.Address, error) {
	return _SurgeAnchor.Contract.Impl(&_SurgeAnchor.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_SurgeAnchor *SurgeAnchorCallerSession) Impl() (common.Address, error) {
	return _SurgeAnchor.Contract.Impl(&_SurgeAnchor.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_SurgeAnchor *SurgeAnchorCaller) InNonReentrant(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "inNonReentrant")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_SurgeAnchor *SurgeAnchorSession) InNonReentrant() (bool, error) {
	return _SurgeAnchor.Contract.InNonReentrant(&_SurgeAnchor.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_SurgeAnchor *SurgeAnchorCallerSession) InNonReentrant() (bool, error) {
	return _SurgeAnchor.Contract.InNonReentrant(&_SurgeAnchor.CallOpts)
}

// L1ChainId is a free data retrieval call binding the contract method 0x12622e5b.
//
// Solidity: function l1ChainId() view returns(uint64)
func (_SurgeAnchor *SurgeAnchorCaller) L1ChainId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "l1ChainId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// L1ChainId is a free data retrieval call binding the contract method 0x12622e5b.
//
// Solidity: function l1ChainId() view returns(uint64)
func (_SurgeAnchor *SurgeAnchorSession) L1ChainId() (uint64, error) {
	return _SurgeAnchor.Contract.L1ChainId(&_SurgeAnchor.CallOpts)
}

// L1ChainId is a free data retrieval call binding the contract method 0x12622e5b.
//
// Solidity: function l1ChainId() view returns(uint64)
func (_SurgeAnchor *SurgeAnchorCallerSession) L1ChainId() (uint64, error) {
	return _SurgeAnchor.Contract.L1ChainId(&_SurgeAnchor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SurgeAnchor *SurgeAnchorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SurgeAnchor *SurgeAnchorSession) Owner() (common.Address, error) {
	return _SurgeAnchor.Contract.Owner(&_SurgeAnchor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SurgeAnchor *SurgeAnchorCallerSession) Owner() (common.Address, error) {
	return _SurgeAnchor.Contract.Owner(&_SurgeAnchor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SurgeAnchor *SurgeAnchorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SurgeAnchor *SurgeAnchorSession) Paused() (bool, error) {
	return _SurgeAnchor.Contract.Paused(&_SurgeAnchor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SurgeAnchor *SurgeAnchorCallerSession) Paused() (bool, error) {
	return _SurgeAnchor.Contract.Paused(&_SurgeAnchor.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_SurgeAnchor *SurgeAnchorCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_SurgeAnchor *SurgeAnchorSession) PendingOwner() (common.Address, error) {
	return _SurgeAnchor.Contract.PendingOwner(&_SurgeAnchor.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_SurgeAnchor *SurgeAnchorCallerSession) PendingOwner() (common.Address, error) {
	return _SurgeAnchor.Contract.PendingOwner(&_SurgeAnchor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SurgeAnchor *SurgeAnchorCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SurgeAnchor *SurgeAnchorSession) ProxiableUUID() ([32]byte, error) {
	return _SurgeAnchor.Contract.ProxiableUUID(&_SurgeAnchor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SurgeAnchor *SurgeAnchorCallerSession) ProxiableUUID() ([32]byte, error) {
	return _SurgeAnchor.Contract.ProxiableUUID(&_SurgeAnchor.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_SurgeAnchor *SurgeAnchorCaller) Resolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeAnchor.contract.Call(opts, &out, "resolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_SurgeAnchor *SurgeAnchorSession) Resolver() (common.Address, error) {
	return _SurgeAnchor.Contract.Resolver(&_SurgeAnchor.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_SurgeAnchor *SurgeAnchorCallerSession) Resolver() (common.Address, error) {
	return _SurgeAnchor.Contract.Resolver(&_SurgeAnchor.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SurgeAnchor *SurgeAnchorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SurgeAnchor *SurgeAnchorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SurgeAnchor.Contract.AcceptOwnership(&_SurgeAnchor.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SurgeAnchor.Contract.AcceptOwnership(&_SurgeAnchor.TransactOpts)
}

// AnchorV4 is a paid mutator transaction binding the contract method 0x523e6854.
//
// Solidity: function anchorV4((uint48,bytes32,bytes32) _checkpoint) returns()
func (_SurgeAnchor *SurgeAnchorTransactor) AnchorV4(opts *bind.TransactOpts, _checkpoint ICheckpointStoreCheckpoint) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "anchorV4", _checkpoint)
}

// AnchorV4 is a paid mutator transaction binding the contract method 0x523e6854.
//
// Solidity: function anchorV4((uint48,bytes32,bytes32) _checkpoint) returns()
func (_SurgeAnchor *SurgeAnchorSession) AnchorV4(_checkpoint ICheckpointStoreCheckpoint) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.AnchorV4(&_SurgeAnchor.TransactOpts, _checkpoint)
}

// AnchorV4 is a paid mutator transaction binding the contract method 0x523e6854.
//
// Solidity: function anchorV4((uint48,bytes32,bytes32) _checkpoint) returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) AnchorV4(_checkpoint ICheckpointStoreCheckpoint) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.AnchorV4(&_SurgeAnchor.TransactOpts, _checkpoint)
}

// AnchorV4WithSignalSlots is a paid mutator transaction binding the contract method 0xcc07826e.
//
// Solidity: function anchorV4WithSignalSlots((uint48,bytes32,bytes32) _checkpoint, bytes32[] _signalSlots) returns()
func (_SurgeAnchor *SurgeAnchorTransactor) AnchorV4WithSignalSlots(opts *bind.TransactOpts, _checkpoint ICheckpointStoreCheckpoint, _signalSlots [][32]byte) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "anchorV4WithSignalSlots", _checkpoint, _signalSlots)
}

// AnchorV4WithSignalSlots is a paid mutator transaction binding the contract method 0xcc07826e.
//
// Solidity: function anchorV4WithSignalSlots((uint48,bytes32,bytes32) _checkpoint, bytes32[] _signalSlots) returns()
func (_SurgeAnchor *SurgeAnchorSession) AnchorV4WithSignalSlots(_checkpoint ICheckpointStoreCheckpoint, _signalSlots [][32]byte) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.AnchorV4WithSignalSlots(&_SurgeAnchor.TransactOpts, _checkpoint, _signalSlots)
}

// AnchorV4WithSignalSlots is a paid mutator transaction binding the contract method 0xcc07826e.
//
// Solidity: function anchorV4WithSignalSlots((uint48,bytes32,bytes32) _checkpoint, bytes32[] _signalSlots) returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) AnchorV4WithSignalSlots(_checkpoint ICheckpointStoreCheckpoint, _signalSlots [][32]byte) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.AnchorV4WithSignalSlots(&_SurgeAnchor.TransactOpts, _checkpoint, _signalSlots)
}

// AnchorV5 is a paid mutator transaction binding the contract method 0x38da8410.
//
// Solidity: function anchorV5((uint48) _proposalParams, (uint48,bytes32,bytes32,bytes32) _blockParams) returns()
func (_SurgeAnchor *SurgeAnchorTransactor) AnchorV5(opts *bind.TransactOpts, _proposalParams AnchorProposalParams, _blockParams AnchorBlockParams) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "anchorV5", _proposalParams, _blockParams)
}

// AnchorV5 is a paid mutator transaction binding the contract method 0x38da8410.
//
// Solidity: function anchorV5((uint48) _proposalParams, (uint48,bytes32,bytes32,bytes32) _blockParams) returns()
func (_SurgeAnchor *SurgeAnchorSession) AnchorV5(_proposalParams AnchorProposalParams, _blockParams AnchorBlockParams) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.AnchorV5(&_SurgeAnchor.TransactOpts, _proposalParams, _blockParams)
}

// AnchorV5 is a paid mutator transaction binding the contract method 0x38da8410.
//
// Solidity: function anchorV5((uint48) _proposalParams, (uint48,bytes32,bytes32,bytes32) _blockParams) returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) AnchorV5(_proposalParams AnchorProposalParams, _blockParams AnchorBlockParams) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.AnchorV5(&_SurgeAnchor.TransactOpts, _proposalParams, _blockParams)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_SurgeAnchor *SurgeAnchorTransactor) Init(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "init", _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_SurgeAnchor *SurgeAnchorSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.Init(&_SurgeAnchor.TransactOpts, _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.Init(&_SurgeAnchor.TransactOpts, _owner)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SurgeAnchor *SurgeAnchorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SurgeAnchor *SurgeAnchorSession) Pause() (*types.Transaction, error) {
	return _SurgeAnchor.Contract.Pause(&_SurgeAnchor.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) Pause() (*types.Transaction, error) {
	return _SurgeAnchor.Contract.Pause(&_SurgeAnchor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SurgeAnchor *SurgeAnchorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SurgeAnchor *SurgeAnchorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SurgeAnchor.Contract.RenounceOwnership(&_SurgeAnchor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SurgeAnchor.Contract.RenounceOwnership(&_SurgeAnchor.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SurgeAnchor *SurgeAnchorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SurgeAnchor *SurgeAnchorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.TransferOwnership(&_SurgeAnchor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.TransferOwnership(&_SurgeAnchor.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SurgeAnchor *SurgeAnchorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SurgeAnchor *SurgeAnchorSession) Unpause() (*types.Transaction, error) {
	return _SurgeAnchor.Contract.Unpause(&_SurgeAnchor.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) Unpause() (*types.Transaction, error) {
	return _SurgeAnchor.Contract.Unpause(&_SurgeAnchor.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SurgeAnchor *SurgeAnchorTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SurgeAnchor *SurgeAnchorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.UpgradeTo(&_SurgeAnchor.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.UpgradeTo(&_SurgeAnchor.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SurgeAnchor *SurgeAnchorTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SurgeAnchor *SurgeAnchorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.UpgradeToAndCall(&_SurgeAnchor.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.UpgradeToAndCall(&_SurgeAnchor.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(address _token, address _to) returns()
func (_SurgeAnchor *SurgeAnchorTransactor) Withdraw(opts *bind.TransactOpts, _token common.Address, _to common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.contract.Transact(opts, "withdraw", _token, _to)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(address _token, address _to) returns()
func (_SurgeAnchor *SurgeAnchorSession) Withdraw(_token common.Address, _to common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.Withdraw(&_SurgeAnchor.TransactOpts, _token, _to)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(address _token, address _to) returns()
func (_SurgeAnchor *SurgeAnchorTransactorSession) Withdraw(_token common.Address, _to common.Address) (*types.Transaction, error) {
	return _SurgeAnchor.Contract.Withdraw(&_SurgeAnchor.TransactOpts, _token, _to)
}

// SurgeAnchorAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the SurgeAnchor contract.
type SurgeAnchorAdminChangedIterator struct {
	Event *SurgeAnchorAdminChanged // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorAdminChanged)
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
		it.Event = new(SurgeAnchorAdminChanged)
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
func (it *SurgeAnchorAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorAdminChanged represents a AdminChanged event raised by the SurgeAnchor contract.
type SurgeAnchorAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*SurgeAnchorAdminChangedIterator, error) {

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorAdminChangedIterator{contract: _SurgeAnchor.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *SurgeAnchorAdminChanged) (event.Subscription, error) {

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorAdminChanged)
				if err := _SurgeAnchor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_SurgeAnchor *SurgeAnchorFilterer) ParseAdminChanged(log types.Log) (*SurgeAnchorAdminChanged, error) {
	event := new(SurgeAnchorAdminChanged)
	if err := _SurgeAnchor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeAnchorAnchoredIterator is returned from FilterAnchored and is used to iterate over the raw logs and unpacked data for Anchored events raised by the SurgeAnchor contract.
type SurgeAnchorAnchoredIterator struct {
	Event *SurgeAnchorAnchored // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorAnchoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorAnchored)
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
		it.Event = new(SurgeAnchorAnchored)
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
func (it *SurgeAnchorAnchoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorAnchoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorAnchored represents a Anchored event raised by the SurgeAnchor contract.
type SurgeAnchorAnchored struct {
	PrevAnchorBlockNumber *big.Int
	AnchorBlockNumber     *big.Int
	AncestorsHash         [32]byte
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterAnchored is a free log retrieval operation binding the contract event 0x7efb645089c13cff3d84f896f1bf50411e06ae8453a771bd7e3ea45ffb97ac38.
//
// Solidity: event Anchored(uint48 prevAnchorBlockNumber, uint48 anchorBlockNumber, bytes32 ancestorsHash)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterAnchored(opts *bind.FilterOpts) (*SurgeAnchorAnchoredIterator, error) {

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "Anchored")
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorAnchoredIterator{contract: _SurgeAnchor.contract, event: "Anchored", logs: logs, sub: sub}, nil
}

// WatchAnchored is a free log subscription operation binding the contract event 0x7efb645089c13cff3d84f896f1bf50411e06ae8453a771bd7e3ea45ffb97ac38.
//
// Solidity: event Anchored(uint48 prevAnchorBlockNumber, uint48 anchorBlockNumber, bytes32 ancestorsHash)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchAnchored(opts *bind.WatchOpts, sink chan<- *SurgeAnchorAnchored) (event.Subscription, error) {

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "Anchored")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorAnchored)
				if err := _SurgeAnchor.contract.UnpackLog(event, "Anchored", log); err != nil {
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

// ParseAnchored is a log parse operation binding the contract event 0x7efb645089c13cff3d84f896f1bf50411e06ae8453a771bd7e3ea45ffb97ac38.
//
// Solidity: event Anchored(uint48 prevAnchorBlockNumber, uint48 anchorBlockNumber, bytes32 ancestorsHash)
func (_SurgeAnchor *SurgeAnchorFilterer) ParseAnchored(log types.Log) (*SurgeAnchorAnchored, error) {
	event := new(SurgeAnchorAnchored)
	if err := _SurgeAnchor.contract.UnpackLog(event, "Anchored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeAnchorBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the SurgeAnchor contract.
type SurgeAnchorBeaconUpgradedIterator struct {
	Event *SurgeAnchorBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorBeaconUpgraded)
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
		it.Event = new(SurgeAnchorBeaconUpgraded)
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
func (it *SurgeAnchorBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorBeaconUpgraded represents a BeaconUpgraded event raised by the SurgeAnchor contract.
type SurgeAnchorBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*SurgeAnchorBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorBeaconUpgradedIterator{contract: _SurgeAnchor.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *SurgeAnchorBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorBeaconUpgraded)
				if err := _SurgeAnchor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_SurgeAnchor *SurgeAnchorFilterer) ParseBeaconUpgraded(log types.Log) (*SurgeAnchorBeaconUpgraded, error) {
	event := new(SurgeAnchorBeaconUpgraded)
	if err := _SurgeAnchor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeAnchorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SurgeAnchor contract.
type SurgeAnchorInitializedIterator struct {
	Event *SurgeAnchorInitialized // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorInitialized)
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
		it.Event = new(SurgeAnchorInitialized)
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
func (it *SurgeAnchorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorInitialized represents a Initialized event raised by the SurgeAnchor contract.
type SurgeAnchorInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterInitialized(opts *bind.FilterOpts) (*SurgeAnchorInitializedIterator, error) {

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorInitializedIterator{contract: _SurgeAnchor.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SurgeAnchorInitialized) (event.Subscription, error) {

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorInitialized)
				if err := _SurgeAnchor.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_SurgeAnchor *SurgeAnchorFilterer) ParseInitialized(log types.Log) (*SurgeAnchorInitialized, error) {
	event := new(SurgeAnchorInitialized)
	if err := _SurgeAnchor.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeAnchorOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the SurgeAnchor contract.
type SurgeAnchorOwnershipTransferStartedIterator struct {
	Event *SurgeAnchorOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorOwnershipTransferStarted)
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
		it.Event = new(SurgeAnchorOwnershipTransferStarted)
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
func (it *SurgeAnchorOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the SurgeAnchor contract.
type SurgeAnchorOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SurgeAnchorOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorOwnershipTransferStartedIterator{contract: _SurgeAnchor.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *SurgeAnchorOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorOwnershipTransferStarted)
				if err := _SurgeAnchor.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_SurgeAnchor *SurgeAnchorFilterer) ParseOwnershipTransferStarted(log types.Log) (*SurgeAnchorOwnershipTransferStarted, error) {
	event := new(SurgeAnchorOwnershipTransferStarted)
	if err := _SurgeAnchor.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeAnchorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SurgeAnchor contract.
type SurgeAnchorOwnershipTransferredIterator struct {
	Event *SurgeAnchorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorOwnershipTransferred)
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
		it.Event = new(SurgeAnchorOwnershipTransferred)
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
func (it *SurgeAnchorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorOwnershipTransferred represents a OwnershipTransferred event raised by the SurgeAnchor contract.
type SurgeAnchorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SurgeAnchorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorOwnershipTransferredIterator{contract: _SurgeAnchor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SurgeAnchorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorOwnershipTransferred)
				if err := _SurgeAnchor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SurgeAnchor *SurgeAnchorFilterer) ParseOwnershipTransferred(log types.Log) (*SurgeAnchorOwnershipTransferred, error) {
	event := new(SurgeAnchorOwnershipTransferred)
	if err := _SurgeAnchor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeAnchorPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the SurgeAnchor contract.
type SurgeAnchorPausedIterator struct {
	Event *SurgeAnchorPaused // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorPaused)
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
		it.Event = new(SurgeAnchorPaused)
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
func (it *SurgeAnchorPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorPaused represents a Paused event raised by the SurgeAnchor contract.
type SurgeAnchorPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterPaused(opts *bind.FilterOpts) (*SurgeAnchorPausedIterator, error) {

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorPausedIterator{contract: _SurgeAnchor.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *SurgeAnchorPaused) (event.Subscription, error) {

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorPaused)
				if err := _SurgeAnchor.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_SurgeAnchor *SurgeAnchorFilterer) ParsePaused(log types.Log) (*SurgeAnchorPaused, error) {
	event := new(SurgeAnchorPaused)
	if err := _SurgeAnchor.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeAnchorUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the SurgeAnchor contract.
type SurgeAnchorUnpausedIterator struct {
	Event *SurgeAnchorUnpaused // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorUnpaused)
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
		it.Event = new(SurgeAnchorUnpaused)
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
func (it *SurgeAnchorUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorUnpaused represents a Unpaused event raised by the SurgeAnchor contract.
type SurgeAnchorUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*SurgeAnchorUnpausedIterator, error) {

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorUnpausedIterator{contract: _SurgeAnchor.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *SurgeAnchorUnpaused) (event.Subscription, error) {

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorUnpaused)
				if err := _SurgeAnchor.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_SurgeAnchor *SurgeAnchorFilterer) ParseUnpaused(log types.Log) (*SurgeAnchorUnpaused, error) {
	event := new(SurgeAnchorUnpaused)
	if err := _SurgeAnchor.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeAnchorUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the SurgeAnchor contract.
type SurgeAnchorUpgradedIterator struct {
	Event *SurgeAnchorUpgraded // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorUpgraded)
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
		it.Event = new(SurgeAnchorUpgraded)
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
func (it *SurgeAnchorUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorUpgraded represents a Upgraded event raised by the SurgeAnchor contract.
type SurgeAnchorUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SurgeAnchorUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorUpgradedIterator{contract: _SurgeAnchor.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SurgeAnchorUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorUpgraded)
				if err := _SurgeAnchor.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_SurgeAnchor *SurgeAnchorFilterer) ParseUpgraded(log types.Log) (*SurgeAnchorUpgraded, error) {
	event := new(SurgeAnchorUpgraded)
	if err := _SurgeAnchor.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeAnchorWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the SurgeAnchor contract.
type SurgeAnchorWithdrawnIterator struct {
	Event *SurgeAnchorWithdrawn // Event containing the contract specifics and raw log

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
func (it *SurgeAnchorWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeAnchorWithdrawn)
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
		it.Event = new(SurgeAnchorWithdrawn)
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
func (it *SurgeAnchorWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeAnchorWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeAnchorWithdrawn represents a Withdrawn event raised by the SurgeAnchor contract.
type SurgeAnchorWithdrawn struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address token, address to, uint256 amount)
func (_SurgeAnchor *SurgeAnchorFilterer) FilterWithdrawn(opts *bind.FilterOpts) (*SurgeAnchorWithdrawnIterator, error) {

	logs, sub, err := _SurgeAnchor.contract.FilterLogs(opts, "Withdrawn")
	if err != nil {
		return nil, err
	}
	return &SurgeAnchorWithdrawnIterator{contract: _SurgeAnchor.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address token, address to, uint256 amount)
func (_SurgeAnchor *SurgeAnchorFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *SurgeAnchorWithdrawn) (event.Subscription, error) {

	logs, sub, err := _SurgeAnchor.contract.WatchLogs(opts, "Withdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeAnchorWithdrawn)
				if err := _SurgeAnchor.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address token, address to, uint256 amount)
func (_SurgeAnchor *SurgeAnchorFilterer) ParseWithdrawn(log types.Log) (*SurgeAnchorWithdrawn, error) {
	event := new(SurgeAnchorWithdrawn)
	if err := _SurgeAnchor.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
