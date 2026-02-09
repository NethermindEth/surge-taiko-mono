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

// LibProofBitmapMetaData contains all meta data concerning the LibProofBitmap contract.
var LibProofBitmapMetaData = &bind.MetaData{
	ABI: "[]",
}

// LibProofBitmapABI is the input ABI used to generate the binding from.
// Deprecated: Use LibProofBitmapMetaData.ABI instead.
var LibProofBitmapABI = LibProofBitmapMetaData.ABI

// LibProofBitmap is an auto generated Go binding around an Ethereum contract.
type LibProofBitmap struct {
	LibProofBitmapCaller     // Read-only binding to the contract
	LibProofBitmapTransactor // Write-only binding to the contract
	LibProofBitmapFilterer   // Log filterer for contract events
}

// LibProofBitmapCaller is an auto generated read-only Go binding around an Ethereum contract.
type LibProofBitmapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibProofBitmapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LibProofBitmapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibProofBitmapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LibProofBitmapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibProofBitmapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LibProofBitmapSession struct {
	Contract     *LibProofBitmap   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LibProofBitmapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LibProofBitmapCallerSession struct {
	Contract *LibProofBitmapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// LibProofBitmapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LibProofBitmapTransactorSession struct {
	Contract     *LibProofBitmapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// LibProofBitmapRaw is an auto generated low-level Go binding around an Ethereum contract.
type LibProofBitmapRaw struct {
	Contract *LibProofBitmap // Generic contract binding to access the raw methods on
}

// LibProofBitmapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LibProofBitmapCallerRaw struct {
	Contract *LibProofBitmapCaller // Generic read-only contract binding to access the raw methods on
}

// LibProofBitmapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LibProofBitmapTransactorRaw struct {
	Contract *LibProofBitmapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLibProofBitmap creates a new instance of LibProofBitmap, bound to a specific deployed contract.
func NewLibProofBitmap(address common.Address, backend bind.ContractBackend) (*LibProofBitmap, error) {
	contract, err := bindLibProofBitmap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LibProofBitmap{LibProofBitmapCaller: LibProofBitmapCaller{contract: contract}, LibProofBitmapTransactor: LibProofBitmapTransactor{contract: contract}, LibProofBitmapFilterer: LibProofBitmapFilterer{contract: contract}}, nil
}

// NewLibProofBitmapCaller creates a new read-only instance of LibProofBitmap, bound to a specific deployed contract.
func NewLibProofBitmapCaller(address common.Address, caller bind.ContractCaller) (*LibProofBitmapCaller, error) {
	contract, err := bindLibProofBitmap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LibProofBitmapCaller{contract: contract}, nil
}

// NewLibProofBitmapTransactor creates a new write-only instance of LibProofBitmap, bound to a specific deployed contract.
func NewLibProofBitmapTransactor(address common.Address, transactor bind.ContractTransactor) (*LibProofBitmapTransactor, error) {
	contract, err := bindLibProofBitmap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LibProofBitmapTransactor{contract: contract}, nil
}

// NewLibProofBitmapFilterer creates a new log filterer instance of LibProofBitmap, bound to a specific deployed contract.
func NewLibProofBitmapFilterer(address common.Address, filterer bind.ContractFilterer) (*LibProofBitmapFilterer, error) {
	contract, err := bindLibProofBitmap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LibProofBitmapFilterer{contract: contract}, nil
}

// bindLibProofBitmap binds a generic wrapper to an already deployed contract.
func bindLibProofBitmap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LibProofBitmapMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibProofBitmap *LibProofBitmapRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LibProofBitmap.Contract.LibProofBitmapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibProofBitmap *LibProofBitmapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibProofBitmap.Contract.LibProofBitmapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibProofBitmap *LibProofBitmapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibProofBitmap.Contract.LibProofBitmapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibProofBitmap *LibProofBitmapCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LibProofBitmap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibProofBitmap *LibProofBitmapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibProofBitmap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibProofBitmap *LibProofBitmapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibProofBitmap.Contract.contract.Transact(opts, method, params...)
}
