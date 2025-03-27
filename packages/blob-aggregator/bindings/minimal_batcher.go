// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// Call represents a call structure for batch execution
type Call struct {
	Target common.Address
	Value  *big.Int
	Data   []byte
}

// MinimalBatcherMetaData contains all meta data concerning the MinimalBatcher contract.
var MinimalBatcherMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"executeBatch\",\"inputs\":[{\"name\":\"calls\",\"type\":\"tuple[]\",\"internalType\":\"struct MinimalBatcher.Call[]\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"error\",\"name\":\"CALL_FAILED\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"INVALID_ETHER_AMOUNT\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NOT_AUTHORIZED\",\"inputs\":[]}]",
}

// MinimalBatcherABI is the input ABI used to generate the binding from.
// Deprecated: Use MinimalBatcherMetaData.ABI instead.
var MinimalBatcherABI = MinimalBatcherMetaData.ABI

// MinimalBatcher is an auto generated Go binding around an Ethereum contract.
type MinimalBatcher struct {
	MinimalBatcherCaller     // Read-only binding to the contract
	MinimalBatcherTransactor // Write-only binding to the contract
	MinimalBatcherFilterer   // Log filterer for contract events
}

// MinimalBatcherCaller is an auto generated read-only Go binding around an Ethereum contract.
type MinimalBatcherCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MinimalBatcherTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MinimalBatcherTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MinimalBatcherFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MinimalBatcherFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MinimalBatcherSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MinimalBatcherSession struct {
	Contract     *MinimalBatcher   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction options to use throughout this session
}

// MinimalBatcherCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MinimalBatcherCallerSession struct {
	Contract *MinimalBatcherCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// MinimalBatcherTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MinimalBatcherTransactorSession struct {
	Contract     *MinimalBatcherTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction options to use throughout this session
}

// MinimalBatcherRaw is an auto generated low-level Go binding around an Ethereum contract.
type MinimalBatcherRaw struct {
	Contract *MinimalBatcher // Generic contract binding to access the raw methods on
}

// MinimalBatcherCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MinimalBatcherCallerRaw struct {
	Contract *MinimalBatcherCaller // Generic read-only contract binding to access the raw methods on
}

// MinimalBatcherTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MinimalBatcherTransactorRaw struct {
	Contract *MinimalBatcherTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMinimalBatcher creates a new instance of MinimalBatcher, bound to a specific deployed contract.
func NewMinimalBatcher(address common.Address, backend bind.ContractBackend) (*MinimalBatcher, error) {
	contract, err := bindMinimalBatcher(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MinimalBatcher{MinimalBatcherCaller: MinimalBatcherCaller{contract: contract}, MinimalBatcherTransactor: MinimalBatcherTransactor{contract: contract}, MinimalBatcherFilterer: MinimalBatcherFilterer{contract: contract}}, nil
}

// NewMinimalBatcherCaller creates a new read-only instance of MinimalBatcher, bound to a specific deployed contract.
func NewMinimalBatcherCaller(address common.Address, caller bind.ContractCaller) (*MinimalBatcherCaller, error) {
	contract, err := bindMinimalBatcher(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MinimalBatcherCaller{contract: contract}, nil
}

// NewMinimalBatcherTransactor creates a new write-only instance of MinimalBatcher, bound to a specific deployed contract.
func NewMinimalBatcherTransactor(address common.Address, transactor bind.ContractTransactor) (*MinimalBatcherTransactor, error) {
	contract, err := bindMinimalBatcher(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MinimalBatcherTransactor{contract: contract}, nil
}

// NewMinimalBatcherFilterer creates a new log filterer instance of MinimalBatcher, bound to a specific deployed contract.
func NewMinimalBatcherFilterer(address common.Address, filterer bind.ContractFilterer) (*MinimalBatcherFilterer, error) {
	contract, err := bindMinimalBatcher(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MinimalBatcherFilterer{contract: contract}, nil
}

// bindMinimalBatcher binds a generic wrapper to an already deployed contract.
func bindMinimalBatcher(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MinimalBatcherMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MinimalBatcher *MinimalBatcherRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MinimalBatcher.Contract.MinimalBatcherCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MinimalBatcher *MinimalBatcherRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MinimalBatcher.Contract.MinimalBatcherTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MinimalBatcher *MinimalBatcherRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MinimalBatcher.Contract.MinimalBatcherTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MinimalBatcher *MinimalBatcherCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MinimalBatcher.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MinimalBatcher *MinimalBatcherTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MinimalBatcher.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MinimalBatcher *MinimalBatcherTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MinimalBatcher.Contract.contract.Transact(opts, method, params...)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0x34fcd5be.
//
// Solidity: function executeBatch((address,uint256,bytes)[] calls) payable returns()
func (_MinimalBatcher *MinimalBatcherTransactor) ExecuteBatch(opts *bind.TransactOpts, calls []Call) (*types.Transaction, error) {
	return _MinimalBatcher.contract.Transact(opts, "executeBatch", calls)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0x34fcd5be.
//
// Solidity: function executeBatch((address,uint256,bytes)[] calls) payable returns()
func (_MinimalBatcher *MinimalBatcherSession) ExecuteBatch(calls []Call) (*types.Transaction, error) {
	return _MinimalBatcher.Contract.ExecuteBatch(&_MinimalBatcher.TransactOpts, calls)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0x34fcd5be.
//
// Solidity: function executeBatch((address,uint256,bytes)[] calls) payable returns()
func (_MinimalBatcher *MinimalBatcherTransactorSession) ExecuteBatch(calls []Call) (*types.Transaction, error) {
	return _MinimalBatcher.Contract.ExecuteBatch(&_MinimalBatcher.TransactOpts, calls)
}