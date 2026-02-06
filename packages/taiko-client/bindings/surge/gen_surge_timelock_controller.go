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

// SurgeTimelockControllerMetaData contains all meta data concerning the SurgeTimelockController contract.
var SurgeTimelockControllerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_inbox\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_proofVerifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minFinalizationStreak\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"_minDelay\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_proposers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_executors\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"CANCELLER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EXECUTOR_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PROPOSER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TIMELOCK_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[{\"name\":\"_contracts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancel\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"_target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_predecessor\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"executeBatch\",\"inputs\":[{\"name\":\"_targets\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"_payloads\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"_predecessor\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"executeSetVerifierInstant\",\"inputs\":[{\"name\":\"_proofBitFlag\",\"type\":\"uint8\",\"internalType\":\"LibProofBitmap.ProofBitmap\"},{\"name\":\"_newVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getMinDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTimestamp\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashOperation\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"predecessor\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"hashOperationBatch\",\"inputs\":[{\"name\":\"targets\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"payloads\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"predecessor\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"inbox\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isOperation\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isOperationDone\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isOperationPending\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isOperationReady\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minFinalizationStreak\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onERC1155BatchReceived\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC1155Received\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC721Received\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proofVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"schedule\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"predecessor\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"delay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"scheduleBatch\",\"inputs\":[{\"name\":\"targets\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"payloads\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"predecessor\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"delay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CallExecuted\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"target\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CallSalt\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CallScheduled\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"target\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"predecessor\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"delay\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Cancelled\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinDelayChange\",\"inputs\":[{\"name\":\"oldDuration\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newDuration\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"FinalizationStreakDisrupted\",\"inputs\":[]}]",
}

// SurgeTimelockControllerABI is the input ABI used to generate the binding from.
// Deprecated: Use SurgeTimelockControllerMetaData.ABI instead.
var SurgeTimelockControllerABI = SurgeTimelockControllerMetaData.ABI

// SurgeTimelockController is an auto generated Go binding around an Ethereum contract.
type SurgeTimelockController struct {
	SurgeTimelockControllerCaller     // Read-only binding to the contract
	SurgeTimelockControllerTransactor // Write-only binding to the contract
	SurgeTimelockControllerFilterer   // Log filterer for contract events
}

// SurgeTimelockControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SurgeTimelockControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeTimelockControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SurgeTimelockControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeTimelockControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SurgeTimelockControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeTimelockControllerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SurgeTimelockControllerSession struct {
	Contract     *SurgeTimelockController // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SurgeTimelockControllerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SurgeTimelockControllerCallerSession struct {
	Contract *SurgeTimelockControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// SurgeTimelockControllerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SurgeTimelockControllerTransactorSession struct {
	Contract     *SurgeTimelockControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// SurgeTimelockControllerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SurgeTimelockControllerRaw struct {
	Contract *SurgeTimelockController // Generic contract binding to access the raw methods on
}

// SurgeTimelockControllerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SurgeTimelockControllerCallerRaw struct {
	Contract *SurgeTimelockControllerCaller // Generic read-only contract binding to access the raw methods on
}

// SurgeTimelockControllerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SurgeTimelockControllerTransactorRaw struct {
	Contract *SurgeTimelockControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSurgeTimelockController creates a new instance of SurgeTimelockController, bound to a specific deployed contract.
func NewSurgeTimelockController(address common.Address, backend bind.ContractBackend) (*SurgeTimelockController, error) {
	contract, err := bindSurgeTimelockController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockController{SurgeTimelockControllerCaller: SurgeTimelockControllerCaller{contract: contract}, SurgeTimelockControllerTransactor: SurgeTimelockControllerTransactor{contract: contract}, SurgeTimelockControllerFilterer: SurgeTimelockControllerFilterer{contract: contract}}, nil
}

// NewSurgeTimelockControllerCaller creates a new read-only instance of SurgeTimelockController, bound to a specific deployed contract.
func NewSurgeTimelockControllerCaller(address common.Address, caller bind.ContractCaller) (*SurgeTimelockControllerCaller, error) {
	contract, err := bindSurgeTimelockController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerCaller{contract: contract}, nil
}

// NewSurgeTimelockControllerTransactor creates a new write-only instance of SurgeTimelockController, bound to a specific deployed contract.
func NewSurgeTimelockControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*SurgeTimelockControllerTransactor, error) {
	contract, err := bindSurgeTimelockController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerTransactor{contract: contract}, nil
}

// NewSurgeTimelockControllerFilterer creates a new log filterer instance of SurgeTimelockController, bound to a specific deployed contract.
func NewSurgeTimelockControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*SurgeTimelockControllerFilterer, error) {
	contract, err := bindSurgeTimelockController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerFilterer{contract: contract}, nil
}

// bindSurgeTimelockController binds a generic wrapper to an already deployed contract.
func bindSurgeTimelockController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SurgeTimelockControllerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SurgeTimelockController *SurgeTimelockControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SurgeTimelockController.Contract.SurgeTimelockControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SurgeTimelockController *SurgeTimelockControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.SurgeTimelockControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SurgeTimelockController *SurgeTimelockControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.SurgeTimelockControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SurgeTimelockController *SurgeTimelockControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SurgeTimelockController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SurgeTimelockController *SurgeTimelockControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SurgeTimelockController *SurgeTimelockControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.contract.Transact(opts, method, params...)
}

// CANCELLERROLE is a free data retrieval call binding the contract method 0xb08e51c0.
//
// Solidity: function CANCELLER_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) CANCELLERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "CANCELLER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CANCELLERROLE is a free data retrieval call binding the contract method 0xb08e51c0.
//
// Solidity: function CANCELLER_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerSession) CANCELLERROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.CANCELLERROLE(&_SurgeTimelockController.CallOpts)
}

// CANCELLERROLE is a free data retrieval call binding the contract method 0xb08e51c0.
//
// Solidity: function CANCELLER_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) CANCELLERROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.CANCELLERROLE(&_SurgeTimelockController.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.DEFAULTADMINROLE(&_SurgeTimelockController.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.DEFAULTADMINROLE(&_SurgeTimelockController.CallOpts)
}

// EXECUTORROLE is a free data retrieval call binding the contract method 0x07bd0265.
//
// Solidity: function EXECUTOR_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) EXECUTORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "EXECUTOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EXECUTORROLE is a free data retrieval call binding the contract method 0x07bd0265.
//
// Solidity: function EXECUTOR_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerSession) EXECUTORROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.EXECUTORROLE(&_SurgeTimelockController.CallOpts)
}

// EXECUTORROLE is a free data retrieval call binding the contract method 0x07bd0265.
//
// Solidity: function EXECUTOR_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) EXECUTORROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.EXECUTORROLE(&_SurgeTimelockController.CallOpts)
}

// PROPOSERROLE is a free data retrieval call binding the contract method 0x8f61f4f5.
//
// Solidity: function PROPOSER_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) PROPOSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "PROPOSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PROPOSERROLE is a free data retrieval call binding the contract method 0x8f61f4f5.
//
// Solidity: function PROPOSER_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerSession) PROPOSERROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.PROPOSERROLE(&_SurgeTimelockController.CallOpts)
}

// PROPOSERROLE is a free data retrieval call binding the contract method 0x8f61f4f5.
//
// Solidity: function PROPOSER_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) PROPOSERROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.PROPOSERROLE(&_SurgeTimelockController.CallOpts)
}

// TIMELOCKADMINROLE is a free data retrieval call binding the contract method 0x0d3cf6fc.
//
// Solidity: function TIMELOCK_ADMIN_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) TIMELOCKADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "TIMELOCK_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TIMELOCKADMINROLE is a free data retrieval call binding the contract method 0x0d3cf6fc.
//
// Solidity: function TIMELOCK_ADMIN_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerSession) TIMELOCKADMINROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.TIMELOCKADMINROLE(&_SurgeTimelockController.CallOpts)
}

// TIMELOCKADMINROLE is a free data retrieval call binding the contract method 0x0d3cf6fc.
//
// Solidity: function TIMELOCK_ADMIN_ROLE() view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) TIMELOCKADMINROLE() ([32]byte, error) {
	return _SurgeTimelockController.Contract.TIMELOCKADMINROLE(&_SurgeTimelockController.CallOpts)
}

// GetMinDelay is a free data retrieval call binding the contract method 0xf27a0c92.
//
// Solidity: function getMinDelay() view returns(uint256)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) GetMinDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "getMinDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinDelay is a free data retrieval call binding the contract method 0xf27a0c92.
//
// Solidity: function getMinDelay() view returns(uint256)
func (_SurgeTimelockController *SurgeTimelockControllerSession) GetMinDelay() (*big.Int, error) {
	return _SurgeTimelockController.Contract.GetMinDelay(&_SurgeTimelockController.CallOpts)
}

// GetMinDelay is a free data retrieval call binding the contract method 0xf27a0c92.
//
// Solidity: function getMinDelay() view returns(uint256)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) GetMinDelay() (*big.Int, error) {
	return _SurgeTimelockController.Contract.GetMinDelay(&_SurgeTimelockController.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _SurgeTimelockController.Contract.GetRoleAdmin(&_SurgeTimelockController.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _SurgeTimelockController.Contract.GetRoleAdmin(&_SurgeTimelockController.CallOpts, role)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xd45c4435.
//
// Solidity: function getTimestamp(bytes32 id) view returns(uint256)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) GetTimestamp(opts *bind.CallOpts, id [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "getTimestamp", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xd45c4435.
//
// Solidity: function getTimestamp(bytes32 id) view returns(uint256)
func (_SurgeTimelockController *SurgeTimelockControllerSession) GetTimestamp(id [32]byte) (*big.Int, error) {
	return _SurgeTimelockController.Contract.GetTimestamp(&_SurgeTimelockController.CallOpts, id)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xd45c4435.
//
// Solidity: function getTimestamp(bytes32 id) view returns(uint256)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) GetTimestamp(id [32]byte) (*big.Int, error) {
	return _SurgeTimelockController.Contract.GetTimestamp(&_SurgeTimelockController.CallOpts, id)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _SurgeTimelockController.Contract.HasRole(&_SurgeTimelockController.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _SurgeTimelockController.Contract.HasRole(&_SurgeTimelockController.CallOpts, role, account)
}

// HashOperation is a free data retrieval call binding the contract method 0x8065657f.
//
// Solidity: function hashOperation(address target, uint256 value, bytes data, bytes32 predecessor, bytes32 salt) pure returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) HashOperation(opts *bind.CallOpts, target common.Address, value *big.Int, data []byte, predecessor [32]byte, salt [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "hashOperation", target, value, data, predecessor, salt)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashOperation is a free data retrieval call binding the contract method 0x8065657f.
//
// Solidity: function hashOperation(address target, uint256 value, bytes data, bytes32 predecessor, bytes32 salt) pure returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerSession) HashOperation(target common.Address, value *big.Int, data []byte, predecessor [32]byte, salt [32]byte) ([32]byte, error) {
	return _SurgeTimelockController.Contract.HashOperation(&_SurgeTimelockController.CallOpts, target, value, data, predecessor, salt)
}

// HashOperation is a free data retrieval call binding the contract method 0x8065657f.
//
// Solidity: function hashOperation(address target, uint256 value, bytes data, bytes32 predecessor, bytes32 salt) pure returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) HashOperation(target common.Address, value *big.Int, data []byte, predecessor [32]byte, salt [32]byte) ([32]byte, error) {
	return _SurgeTimelockController.Contract.HashOperation(&_SurgeTimelockController.CallOpts, target, value, data, predecessor, salt)
}

// HashOperationBatch is a free data retrieval call binding the contract method 0xb1c5f427.
//
// Solidity: function hashOperationBatch(address[] targets, uint256[] values, bytes[] payloads, bytes32 predecessor, bytes32 salt) pure returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) HashOperationBatch(opts *bind.CallOpts, targets []common.Address, values []*big.Int, payloads [][]byte, predecessor [32]byte, salt [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "hashOperationBatch", targets, values, payloads, predecessor, salt)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashOperationBatch is a free data retrieval call binding the contract method 0xb1c5f427.
//
// Solidity: function hashOperationBatch(address[] targets, uint256[] values, bytes[] payloads, bytes32 predecessor, bytes32 salt) pure returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerSession) HashOperationBatch(targets []common.Address, values []*big.Int, payloads [][]byte, predecessor [32]byte, salt [32]byte) ([32]byte, error) {
	return _SurgeTimelockController.Contract.HashOperationBatch(&_SurgeTimelockController.CallOpts, targets, values, payloads, predecessor, salt)
}

// HashOperationBatch is a free data retrieval call binding the contract method 0xb1c5f427.
//
// Solidity: function hashOperationBatch(address[] targets, uint256[] values, bytes[] payloads, bytes32 predecessor, bytes32 salt) pure returns(bytes32)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) HashOperationBatch(targets []common.Address, values []*big.Int, payloads [][]byte, predecessor [32]byte, salt [32]byte) ([32]byte, error) {
	return _SurgeTimelockController.Contract.HashOperationBatch(&_SurgeTimelockController.CallOpts, targets, values, payloads, predecessor, salt)
}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) Inbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "inbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_SurgeTimelockController *SurgeTimelockControllerSession) Inbox() (common.Address, error) {
	return _SurgeTimelockController.Contract.Inbox(&_SurgeTimelockController.CallOpts)
}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) Inbox() (common.Address, error) {
	return _SurgeTimelockController.Contract.Inbox(&_SurgeTimelockController.CallOpts)
}

// IsOperation is a free data retrieval call binding the contract method 0x31d50750.
//
// Solidity: function isOperation(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) IsOperation(opts *bind.CallOpts, id [32]byte) (bool, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "isOperation", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperation is a free data retrieval call binding the contract method 0x31d50750.
//
// Solidity: function isOperation(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerSession) IsOperation(id [32]byte) (bool, error) {
	return _SurgeTimelockController.Contract.IsOperation(&_SurgeTimelockController.CallOpts, id)
}

// IsOperation is a free data retrieval call binding the contract method 0x31d50750.
//
// Solidity: function isOperation(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) IsOperation(id [32]byte) (bool, error) {
	return _SurgeTimelockController.Contract.IsOperation(&_SurgeTimelockController.CallOpts, id)
}

// IsOperationDone is a free data retrieval call binding the contract method 0x2ab0f529.
//
// Solidity: function isOperationDone(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) IsOperationDone(opts *bind.CallOpts, id [32]byte) (bool, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "isOperationDone", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperationDone is a free data retrieval call binding the contract method 0x2ab0f529.
//
// Solidity: function isOperationDone(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerSession) IsOperationDone(id [32]byte) (bool, error) {
	return _SurgeTimelockController.Contract.IsOperationDone(&_SurgeTimelockController.CallOpts, id)
}

// IsOperationDone is a free data retrieval call binding the contract method 0x2ab0f529.
//
// Solidity: function isOperationDone(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) IsOperationDone(id [32]byte) (bool, error) {
	return _SurgeTimelockController.Contract.IsOperationDone(&_SurgeTimelockController.CallOpts, id)
}

// IsOperationPending is a free data retrieval call binding the contract method 0x584b153e.
//
// Solidity: function isOperationPending(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) IsOperationPending(opts *bind.CallOpts, id [32]byte) (bool, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "isOperationPending", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperationPending is a free data retrieval call binding the contract method 0x584b153e.
//
// Solidity: function isOperationPending(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerSession) IsOperationPending(id [32]byte) (bool, error) {
	return _SurgeTimelockController.Contract.IsOperationPending(&_SurgeTimelockController.CallOpts, id)
}

// IsOperationPending is a free data retrieval call binding the contract method 0x584b153e.
//
// Solidity: function isOperationPending(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) IsOperationPending(id [32]byte) (bool, error) {
	return _SurgeTimelockController.Contract.IsOperationPending(&_SurgeTimelockController.CallOpts, id)
}

// IsOperationReady is a free data retrieval call binding the contract method 0x13bc9f20.
//
// Solidity: function isOperationReady(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) IsOperationReady(opts *bind.CallOpts, id [32]byte) (bool, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "isOperationReady", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperationReady is a free data retrieval call binding the contract method 0x13bc9f20.
//
// Solidity: function isOperationReady(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerSession) IsOperationReady(id [32]byte) (bool, error) {
	return _SurgeTimelockController.Contract.IsOperationReady(&_SurgeTimelockController.CallOpts, id)
}

// IsOperationReady is a free data retrieval call binding the contract method 0x13bc9f20.
//
// Solidity: function isOperationReady(bytes32 id) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) IsOperationReady(id [32]byte) (bool, error) {
	return _SurgeTimelockController.Contract.IsOperationReady(&_SurgeTimelockController.CallOpts, id)
}

// MinFinalizationStreak is a free data retrieval call binding the contract method 0xae6b8a04.
//
// Solidity: function minFinalizationStreak() view returns(uint48)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) MinFinalizationStreak(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "minFinalizationStreak")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinFinalizationStreak is a free data retrieval call binding the contract method 0xae6b8a04.
//
// Solidity: function minFinalizationStreak() view returns(uint48)
func (_SurgeTimelockController *SurgeTimelockControllerSession) MinFinalizationStreak() (*big.Int, error) {
	return _SurgeTimelockController.Contract.MinFinalizationStreak(&_SurgeTimelockController.CallOpts)
}

// MinFinalizationStreak is a free data retrieval call binding the contract method 0xae6b8a04.
//
// Solidity: function minFinalizationStreak() view returns(uint48)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) MinFinalizationStreak() (*big.Int, error) {
	return _SurgeTimelockController.Contract.MinFinalizationStreak(&_SurgeTimelockController.CallOpts)
}

// ProofVerifier is a free data retrieval call binding the contract method 0x7fa417b3.
//
// Solidity: function proofVerifier() view returns(address)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) ProofVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "proofVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProofVerifier is a free data retrieval call binding the contract method 0x7fa417b3.
//
// Solidity: function proofVerifier() view returns(address)
func (_SurgeTimelockController *SurgeTimelockControllerSession) ProofVerifier() (common.Address, error) {
	return _SurgeTimelockController.Contract.ProofVerifier(&_SurgeTimelockController.CallOpts)
}

// ProofVerifier is a free data retrieval call binding the contract method 0x7fa417b3.
//
// Solidity: function proofVerifier() view returns(address)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) ProofVerifier() (common.Address, error) {
	return _SurgeTimelockController.Contract.ProofVerifier(&_SurgeTimelockController.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _SurgeTimelockController.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SurgeTimelockController.Contract.SupportsInterface(&_SurgeTimelockController.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SurgeTimelockController *SurgeTimelockControllerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SurgeTimelockController.Contract.SupportsInterface(&_SurgeTimelockController.CallOpts, interfaceId)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0xc15460c6.
//
// Solidity: function acceptOwnership(address[] _contracts) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) AcceptOwnership(opts *bind.TransactOpts, _contracts []common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "acceptOwnership", _contracts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0xc15460c6.
//
// Solidity: function acceptOwnership(address[] _contracts) returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) AcceptOwnership(_contracts []common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.AcceptOwnership(&_SurgeTimelockController.TransactOpts, _contracts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0xc15460c6.
//
// Solidity: function acceptOwnership(address[] _contracts) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) AcceptOwnership(_contracts []common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.AcceptOwnership(&_SurgeTimelockController.TransactOpts, _contracts)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 id) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) Cancel(opts *bind.TransactOpts, id [32]byte) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "cancel", id)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 id) returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) Cancel(id [32]byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.Cancel(&_SurgeTimelockController.TransactOpts, id)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 id) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) Cancel(id [32]byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.Cancel(&_SurgeTimelockController.TransactOpts, id)
}

// Execute is a paid mutator transaction binding the contract method 0x134008d3.
//
// Solidity: function execute(address _target, uint256 _value, bytes _payload, bytes32 _predecessor, bytes32 _salt) payable returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) Execute(opts *bind.TransactOpts, _target common.Address, _value *big.Int, _payload []byte, _predecessor [32]byte, _salt [32]byte) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "execute", _target, _value, _payload, _predecessor, _salt)
}

// Execute is a paid mutator transaction binding the contract method 0x134008d3.
//
// Solidity: function execute(address _target, uint256 _value, bytes _payload, bytes32 _predecessor, bytes32 _salt) payable returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) Execute(_target common.Address, _value *big.Int, _payload []byte, _predecessor [32]byte, _salt [32]byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.Execute(&_SurgeTimelockController.TransactOpts, _target, _value, _payload, _predecessor, _salt)
}

// Execute is a paid mutator transaction binding the contract method 0x134008d3.
//
// Solidity: function execute(address _target, uint256 _value, bytes _payload, bytes32 _predecessor, bytes32 _salt) payable returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) Execute(_target common.Address, _value *big.Int, _payload []byte, _predecessor [32]byte, _salt [32]byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.Execute(&_SurgeTimelockController.TransactOpts, _target, _value, _payload, _predecessor, _salt)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0xe38335e5.
//
// Solidity: function executeBatch(address[] _targets, uint256[] _values, bytes[] _payloads, bytes32 _predecessor, bytes32 _salt) payable returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) ExecuteBatch(opts *bind.TransactOpts, _targets []common.Address, _values []*big.Int, _payloads [][]byte, _predecessor [32]byte, _salt [32]byte) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "executeBatch", _targets, _values, _payloads, _predecessor, _salt)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0xe38335e5.
//
// Solidity: function executeBatch(address[] _targets, uint256[] _values, bytes[] _payloads, bytes32 _predecessor, bytes32 _salt) payable returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) ExecuteBatch(_targets []common.Address, _values []*big.Int, _payloads [][]byte, _predecessor [32]byte, _salt [32]byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.ExecuteBatch(&_SurgeTimelockController.TransactOpts, _targets, _values, _payloads, _predecessor, _salt)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0xe38335e5.
//
// Solidity: function executeBatch(address[] _targets, uint256[] _values, bytes[] _payloads, bytes32 _predecessor, bytes32 _salt) payable returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) ExecuteBatch(_targets []common.Address, _values []*big.Int, _payloads [][]byte, _predecessor [32]byte, _salt [32]byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.ExecuteBatch(&_SurgeTimelockController.TransactOpts, _targets, _values, _payloads, _predecessor, _salt)
}

// ExecuteSetVerifierInstant is a paid mutator transaction binding the contract method 0x1b05f6ac.
//
// Solidity: function executeSetVerifierInstant(uint8 _proofBitFlag, address _newVerifier) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) ExecuteSetVerifierInstant(opts *bind.TransactOpts, _proofBitFlag uint8, _newVerifier common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "executeSetVerifierInstant", _proofBitFlag, _newVerifier)
}

// ExecuteSetVerifierInstant is a paid mutator transaction binding the contract method 0x1b05f6ac.
//
// Solidity: function executeSetVerifierInstant(uint8 _proofBitFlag, address _newVerifier) returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) ExecuteSetVerifierInstant(_proofBitFlag uint8, _newVerifier common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.ExecuteSetVerifierInstant(&_SurgeTimelockController.TransactOpts, _proofBitFlag, _newVerifier)
}

// ExecuteSetVerifierInstant is a paid mutator transaction binding the contract method 0x1b05f6ac.
//
// Solidity: function executeSetVerifierInstant(uint8 _proofBitFlag, address _newVerifier) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) ExecuteSetVerifierInstant(_proofBitFlag uint8, _newVerifier common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.ExecuteSetVerifierInstant(&_SurgeTimelockController.TransactOpts, _proofBitFlag, _newVerifier)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.GrantRole(&_SurgeTimelockController.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.GrantRole(&_SurgeTimelockController.TransactOpts, role, account)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) OnERC1155BatchReceived(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "onERC1155BatchReceived", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_SurgeTimelockController *SurgeTimelockControllerSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.OnERC1155BatchReceived(&_SurgeTimelockController.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.OnERC1155BatchReceived(&_SurgeTimelockController.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) OnERC1155Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "onERC1155Received", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_SurgeTimelockController *SurgeTimelockControllerSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.OnERC1155Received(&_SurgeTimelockController.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.OnERC1155Received(&_SurgeTimelockController.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) OnERC721Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "onERC721Received", arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_SurgeTimelockController *SurgeTimelockControllerSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.OnERC721Received(&_SurgeTimelockController.TransactOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.OnERC721Received(&_SurgeTimelockController.TransactOpts, arg0, arg1, arg2, arg3)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.RenounceRole(&_SurgeTimelockController.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.RenounceRole(&_SurgeTimelockController.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.RevokeRole(&_SurgeTimelockController.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.RevokeRole(&_SurgeTimelockController.TransactOpts, role, account)
}

// Schedule is a paid mutator transaction binding the contract method 0x01d5062a.
//
// Solidity: function schedule(address target, uint256 value, bytes data, bytes32 predecessor, bytes32 salt, uint256 delay) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) Schedule(opts *bind.TransactOpts, target common.Address, value *big.Int, data []byte, predecessor [32]byte, salt [32]byte, delay *big.Int) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "schedule", target, value, data, predecessor, salt, delay)
}

// Schedule is a paid mutator transaction binding the contract method 0x01d5062a.
//
// Solidity: function schedule(address target, uint256 value, bytes data, bytes32 predecessor, bytes32 salt, uint256 delay) returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) Schedule(target common.Address, value *big.Int, data []byte, predecessor [32]byte, salt [32]byte, delay *big.Int) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.Schedule(&_SurgeTimelockController.TransactOpts, target, value, data, predecessor, salt, delay)
}

// Schedule is a paid mutator transaction binding the contract method 0x01d5062a.
//
// Solidity: function schedule(address target, uint256 value, bytes data, bytes32 predecessor, bytes32 salt, uint256 delay) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) Schedule(target common.Address, value *big.Int, data []byte, predecessor [32]byte, salt [32]byte, delay *big.Int) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.Schedule(&_SurgeTimelockController.TransactOpts, target, value, data, predecessor, salt, delay)
}

// ScheduleBatch is a paid mutator transaction binding the contract method 0x8f2a0bb0.
//
// Solidity: function scheduleBatch(address[] targets, uint256[] values, bytes[] payloads, bytes32 predecessor, bytes32 salt, uint256 delay) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) ScheduleBatch(opts *bind.TransactOpts, targets []common.Address, values []*big.Int, payloads [][]byte, predecessor [32]byte, salt [32]byte, delay *big.Int) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "scheduleBatch", targets, values, payloads, predecessor, salt, delay)
}

// ScheduleBatch is a paid mutator transaction binding the contract method 0x8f2a0bb0.
//
// Solidity: function scheduleBatch(address[] targets, uint256[] values, bytes[] payloads, bytes32 predecessor, bytes32 salt, uint256 delay) returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) ScheduleBatch(targets []common.Address, values []*big.Int, payloads [][]byte, predecessor [32]byte, salt [32]byte, delay *big.Int) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.ScheduleBatch(&_SurgeTimelockController.TransactOpts, targets, values, payloads, predecessor, salt, delay)
}

// ScheduleBatch is a paid mutator transaction binding the contract method 0x8f2a0bb0.
//
// Solidity: function scheduleBatch(address[] targets, uint256[] values, bytes[] payloads, bytes32 predecessor, bytes32 salt, uint256 delay) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) ScheduleBatch(targets []common.Address, values []*big.Int, payloads [][]byte, predecessor [32]byte, salt [32]byte, delay *big.Int) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.ScheduleBatch(&_SurgeTimelockController.TransactOpts, targets, values, payloads, predecessor, salt, delay)
}

// UpdateDelay is a paid mutator transaction binding the contract method 0x64d62353.
//
// Solidity: function updateDelay(uint256 newDelay) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) UpdateDelay(opts *bind.TransactOpts, newDelay *big.Int) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.Transact(opts, "updateDelay", newDelay)
}

// UpdateDelay is a paid mutator transaction binding the contract method 0x64d62353.
//
// Solidity: function updateDelay(uint256 newDelay) returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) UpdateDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.UpdateDelay(&_SurgeTimelockController.TransactOpts, newDelay)
}

// UpdateDelay is a paid mutator transaction binding the contract method 0x64d62353.
//
// Solidity: function updateDelay(uint256 newDelay) returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) UpdateDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.UpdateDelay(&_SurgeTimelockController.TransactOpts, newDelay)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeTimelockController.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SurgeTimelockController *SurgeTimelockControllerSession) Receive() (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.Receive(&_SurgeTimelockController.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SurgeTimelockController *SurgeTimelockControllerTransactorSession) Receive() (*types.Transaction, error) {
	return _SurgeTimelockController.Contract.Receive(&_SurgeTimelockController.TransactOpts)
}

// SurgeTimelockControllerCallExecutedIterator is returned from FilterCallExecuted and is used to iterate over the raw logs and unpacked data for CallExecuted events raised by the SurgeTimelockController contract.
type SurgeTimelockControllerCallExecutedIterator struct {
	Event *SurgeTimelockControllerCallExecuted // Event containing the contract specifics and raw log

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
func (it *SurgeTimelockControllerCallExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeTimelockControllerCallExecuted)
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
		it.Event = new(SurgeTimelockControllerCallExecuted)
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
func (it *SurgeTimelockControllerCallExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeTimelockControllerCallExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeTimelockControllerCallExecuted represents a CallExecuted event raised by the SurgeTimelockController contract.
type SurgeTimelockControllerCallExecuted struct {
	Id     [32]byte
	Index  *big.Int
	Target common.Address
	Value  *big.Int
	Data   []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCallExecuted is a free log retrieval operation binding the contract event 0xc2617efa69bab66782fa219543714338489c4e9e178271560a91b82c3f612b58.
//
// Solidity: event CallExecuted(bytes32 indexed id, uint256 indexed index, address target, uint256 value, bytes data)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) FilterCallExecuted(opts *bind.FilterOpts, id [][32]byte, index []*big.Int) (*SurgeTimelockControllerCallExecutedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.FilterLogs(opts, "CallExecuted", idRule, indexRule)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerCallExecutedIterator{contract: _SurgeTimelockController.contract, event: "CallExecuted", logs: logs, sub: sub}, nil
}

// WatchCallExecuted is a free log subscription operation binding the contract event 0xc2617efa69bab66782fa219543714338489c4e9e178271560a91b82c3f612b58.
//
// Solidity: event CallExecuted(bytes32 indexed id, uint256 indexed index, address target, uint256 value, bytes data)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) WatchCallExecuted(opts *bind.WatchOpts, sink chan<- *SurgeTimelockControllerCallExecuted, id [][32]byte, index []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.WatchLogs(opts, "CallExecuted", idRule, indexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeTimelockControllerCallExecuted)
				if err := _SurgeTimelockController.contract.UnpackLog(event, "CallExecuted", log); err != nil {
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

// ParseCallExecuted is a log parse operation binding the contract event 0xc2617efa69bab66782fa219543714338489c4e9e178271560a91b82c3f612b58.
//
// Solidity: event CallExecuted(bytes32 indexed id, uint256 indexed index, address target, uint256 value, bytes data)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) ParseCallExecuted(log types.Log) (*SurgeTimelockControllerCallExecuted, error) {
	event := new(SurgeTimelockControllerCallExecuted)
	if err := _SurgeTimelockController.contract.UnpackLog(event, "CallExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeTimelockControllerCallSaltIterator is returned from FilterCallSalt and is used to iterate over the raw logs and unpacked data for CallSalt events raised by the SurgeTimelockController contract.
type SurgeTimelockControllerCallSaltIterator struct {
	Event *SurgeTimelockControllerCallSalt // Event containing the contract specifics and raw log

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
func (it *SurgeTimelockControllerCallSaltIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeTimelockControllerCallSalt)
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
		it.Event = new(SurgeTimelockControllerCallSalt)
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
func (it *SurgeTimelockControllerCallSaltIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeTimelockControllerCallSaltIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeTimelockControllerCallSalt represents a CallSalt event raised by the SurgeTimelockController contract.
type SurgeTimelockControllerCallSalt struct {
	Id   [32]byte
	Salt [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterCallSalt is a free log retrieval operation binding the contract event 0x20fda5fd27a1ea7bf5b9567f143ac5470bb059374a27e8f67cb44f946f6d0387.
//
// Solidity: event CallSalt(bytes32 indexed id, bytes32 salt)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) FilterCallSalt(opts *bind.FilterOpts, id [][32]byte) (*SurgeTimelockControllerCallSaltIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.FilterLogs(opts, "CallSalt", idRule)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerCallSaltIterator{contract: _SurgeTimelockController.contract, event: "CallSalt", logs: logs, sub: sub}, nil
}

// WatchCallSalt is a free log subscription operation binding the contract event 0x20fda5fd27a1ea7bf5b9567f143ac5470bb059374a27e8f67cb44f946f6d0387.
//
// Solidity: event CallSalt(bytes32 indexed id, bytes32 salt)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) WatchCallSalt(opts *bind.WatchOpts, sink chan<- *SurgeTimelockControllerCallSalt, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.WatchLogs(opts, "CallSalt", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeTimelockControllerCallSalt)
				if err := _SurgeTimelockController.contract.UnpackLog(event, "CallSalt", log); err != nil {
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

// ParseCallSalt is a log parse operation binding the contract event 0x20fda5fd27a1ea7bf5b9567f143ac5470bb059374a27e8f67cb44f946f6d0387.
//
// Solidity: event CallSalt(bytes32 indexed id, bytes32 salt)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) ParseCallSalt(log types.Log) (*SurgeTimelockControllerCallSalt, error) {
	event := new(SurgeTimelockControllerCallSalt)
	if err := _SurgeTimelockController.contract.UnpackLog(event, "CallSalt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeTimelockControllerCallScheduledIterator is returned from FilterCallScheduled and is used to iterate over the raw logs and unpacked data for CallScheduled events raised by the SurgeTimelockController contract.
type SurgeTimelockControllerCallScheduledIterator struct {
	Event *SurgeTimelockControllerCallScheduled // Event containing the contract specifics and raw log

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
func (it *SurgeTimelockControllerCallScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeTimelockControllerCallScheduled)
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
		it.Event = new(SurgeTimelockControllerCallScheduled)
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
func (it *SurgeTimelockControllerCallScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeTimelockControllerCallScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeTimelockControllerCallScheduled represents a CallScheduled event raised by the SurgeTimelockController contract.
type SurgeTimelockControllerCallScheduled struct {
	Id          [32]byte
	Index       *big.Int
	Target      common.Address
	Value       *big.Int
	Data        []byte
	Predecessor [32]byte
	Delay       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCallScheduled is a free log retrieval operation binding the contract event 0x4cf4410cc57040e44862ef0f45f3dd5a5e02db8eb8add648d4b0e236f1d07dca.
//
// Solidity: event CallScheduled(bytes32 indexed id, uint256 indexed index, address target, uint256 value, bytes data, bytes32 predecessor, uint256 delay)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) FilterCallScheduled(opts *bind.FilterOpts, id [][32]byte, index []*big.Int) (*SurgeTimelockControllerCallScheduledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.FilterLogs(opts, "CallScheduled", idRule, indexRule)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerCallScheduledIterator{contract: _SurgeTimelockController.contract, event: "CallScheduled", logs: logs, sub: sub}, nil
}

// WatchCallScheduled is a free log subscription operation binding the contract event 0x4cf4410cc57040e44862ef0f45f3dd5a5e02db8eb8add648d4b0e236f1d07dca.
//
// Solidity: event CallScheduled(bytes32 indexed id, uint256 indexed index, address target, uint256 value, bytes data, bytes32 predecessor, uint256 delay)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) WatchCallScheduled(opts *bind.WatchOpts, sink chan<- *SurgeTimelockControllerCallScheduled, id [][32]byte, index []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.WatchLogs(opts, "CallScheduled", idRule, indexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeTimelockControllerCallScheduled)
				if err := _SurgeTimelockController.contract.UnpackLog(event, "CallScheduled", log); err != nil {
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

// ParseCallScheduled is a log parse operation binding the contract event 0x4cf4410cc57040e44862ef0f45f3dd5a5e02db8eb8add648d4b0e236f1d07dca.
//
// Solidity: event CallScheduled(bytes32 indexed id, uint256 indexed index, address target, uint256 value, bytes data, bytes32 predecessor, uint256 delay)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) ParseCallScheduled(log types.Log) (*SurgeTimelockControllerCallScheduled, error) {
	event := new(SurgeTimelockControllerCallScheduled)
	if err := _SurgeTimelockController.contract.UnpackLog(event, "CallScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeTimelockControllerCancelledIterator is returned from FilterCancelled and is used to iterate over the raw logs and unpacked data for Cancelled events raised by the SurgeTimelockController contract.
type SurgeTimelockControllerCancelledIterator struct {
	Event *SurgeTimelockControllerCancelled // Event containing the contract specifics and raw log

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
func (it *SurgeTimelockControllerCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeTimelockControllerCancelled)
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
		it.Event = new(SurgeTimelockControllerCancelled)
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
func (it *SurgeTimelockControllerCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeTimelockControllerCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeTimelockControllerCancelled represents a Cancelled event raised by the SurgeTimelockController contract.
type SurgeTimelockControllerCancelled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCancelled is a free log retrieval operation binding the contract event 0xbaa1eb22f2a492ba1a5fea61b8df4d27c6c8b5f3971e63bb58fa14ff72eedb70.
//
// Solidity: event Cancelled(bytes32 indexed id)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) FilterCancelled(opts *bind.FilterOpts, id [][32]byte) (*SurgeTimelockControllerCancelledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.FilterLogs(opts, "Cancelled", idRule)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerCancelledIterator{contract: _SurgeTimelockController.contract, event: "Cancelled", logs: logs, sub: sub}, nil
}

// WatchCancelled is a free log subscription operation binding the contract event 0xbaa1eb22f2a492ba1a5fea61b8df4d27c6c8b5f3971e63bb58fa14ff72eedb70.
//
// Solidity: event Cancelled(bytes32 indexed id)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) WatchCancelled(opts *bind.WatchOpts, sink chan<- *SurgeTimelockControllerCancelled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.WatchLogs(opts, "Cancelled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeTimelockControllerCancelled)
				if err := _SurgeTimelockController.contract.UnpackLog(event, "Cancelled", log); err != nil {
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

// ParseCancelled is a log parse operation binding the contract event 0xbaa1eb22f2a492ba1a5fea61b8df4d27c6c8b5f3971e63bb58fa14ff72eedb70.
//
// Solidity: event Cancelled(bytes32 indexed id)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) ParseCancelled(log types.Log) (*SurgeTimelockControllerCancelled, error) {
	event := new(SurgeTimelockControllerCancelled)
	if err := _SurgeTimelockController.contract.UnpackLog(event, "Cancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeTimelockControllerMinDelayChangeIterator is returned from FilterMinDelayChange and is used to iterate over the raw logs and unpacked data for MinDelayChange events raised by the SurgeTimelockController contract.
type SurgeTimelockControllerMinDelayChangeIterator struct {
	Event *SurgeTimelockControllerMinDelayChange // Event containing the contract specifics and raw log

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
func (it *SurgeTimelockControllerMinDelayChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeTimelockControllerMinDelayChange)
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
		it.Event = new(SurgeTimelockControllerMinDelayChange)
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
func (it *SurgeTimelockControllerMinDelayChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeTimelockControllerMinDelayChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeTimelockControllerMinDelayChange represents a MinDelayChange event raised by the SurgeTimelockController contract.
type SurgeTimelockControllerMinDelayChange struct {
	OldDuration *big.Int
	NewDuration *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMinDelayChange is a free log retrieval operation binding the contract event 0x11c24f4ead16507c69ac467fbd5e4eed5fb5c699626d2cc6d66421df253886d5.
//
// Solidity: event MinDelayChange(uint256 oldDuration, uint256 newDuration)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) FilterMinDelayChange(opts *bind.FilterOpts) (*SurgeTimelockControllerMinDelayChangeIterator, error) {

	logs, sub, err := _SurgeTimelockController.contract.FilterLogs(opts, "MinDelayChange")
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerMinDelayChangeIterator{contract: _SurgeTimelockController.contract, event: "MinDelayChange", logs: logs, sub: sub}, nil
}

// WatchMinDelayChange is a free log subscription operation binding the contract event 0x11c24f4ead16507c69ac467fbd5e4eed5fb5c699626d2cc6d66421df253886d5.
//
// Solidity: event MinDelayChange(uint256 oldDuration, uint256 newDuration)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) WatchMinDelayChange(opts *bind.WatchOpts, sink chan<- *SurgeTimelockControllerMinDelayChange) (event.Subscription, error) {

	logs, sub, err := _SurgeTimelockController.contract.WatchLogs(opts, "MinDelayChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeTimelockControllerMinDelayChange)
				if err := _SurgeTimelockController.contract.UnpackLog(event, "MinDelayChange", log); err != nil {
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

// ParseMinDelayChange is a log parse operation binding the contract event 0x11c24f4ead16507c69ac467fbd5e4eed5fb5c699626d2cc6d66421df253886d5.
//
// Solidity: event MinDelayChange(uint256 oldDuration, uint256 newDuration)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) ParseMinDelayChange(log types.Log) (*SurgeTimelockControllerMinDelayChange, error) {
	event := new(SurgeTimelockControllerMinDelayChange)
	if err := _SurgeTimelockController.contract.UnpackLog(event, "MinDelayChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeTimelockControllerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the SurgeTimelockController contract.
type SurgeTimelockControllerRoleAdminChangedIterator struct {
	Event *SurgeTimelockControllerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *SurgeTimelockControllerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeTimelockControllerRoleAdminChanged)
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
		it.Event = new(SurgeTimelockControllerRoleAdminChanged)
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
func (it *SurgeTimelockControllerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeTimelockControllerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeTimelockControllerRoleAdminChanged represents a RoleAdminChanged event raised by the SurgeTimelockController contract.
type SurgeTimelockControllerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*SurgeTimelockControllerRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerRoleAdminChangedIterator{contract: _SurgeTimelockController.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *SurgeTimelockControllerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeTimelockControllerRoleAdminChanged)
				if err := _SurgeTimelockController.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) ParseRoleAdminChanged(log types.Log) (*SurgeTimelockControllerRoleAdminChanged, error) {
	event := new(SurgeTimelockControllerRoleAdminChanged)
	if err := _SurgeTimelockController.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeTimelockControllerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the SurgeTimelockController contract.
type SurgeTimelockControllerRoleGrantedIterator struct {
	Event *SurgeTimelockControllerRoleGranted // Event containing the contract specifics and raw log

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
func (it *SurgeTimelockControllerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeTimelockControllerRoleGranted)
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
		it.Event = new(SurgeTimelockControllerRoleGranted)
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
func (it *SurgeTimelockControllerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeTimelockControllerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeTimelockControllerRoleGranted represents a RoleGranted event raised by the SurgeTimelockController contract.
type SurgeTimelockControllerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*SurgeTimelockControllerRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerRoleGrantedIterator{contract: _SurgeTimelockController.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *SurgeTimelockControllerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeTimelockControllerRoleGranted)
				if err := _SurgeTimelockController.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) ParseRoleGranted(log types.Log) (*SurgeTimelockControllerRoleGranted, error) {
	event := new(SurgeTimelockControllerRoleGranted)
	if err := _SurgeTimelockController.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeTimelockControllerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the SurgeTimelockController contract.
type SurgeTimelockControllerRoleRevokedIterator struct {
	Event *SurgeTimelockControllerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *SurgeTimelockControllerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeTimelockControllerRoleRevoked)
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
		it.Event = new(SurgeTimelockControllerRoleRevoked)
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
func (it *SurgeTimelockControllerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeTimelockControllerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeTimelockControllerRoleRevoked represents a RoleRevoked event raised by the SurgeTimelockController contract.
type SurgeTimelockControllerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*SurgeTimelockControllerRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &SurgeTimelockControllerRoleRevokedIterator{contract: _SurgeTimelockController.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *SurgeTimelockControllerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SurgeTimelockController.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeTimelockControllerRoleRevoked)
				if err := _SurgeTimelockController.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_SurgeTimelockController *SurgeTimelockControllerFilterer) ParseRoleRevoked(log types.Log) (*SurgeTimelockControllerRoleRevoked, error) {
	event := new(SurgeTimelockControllerRoleRevoked)
	if err := _SurgeTimelockController.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
