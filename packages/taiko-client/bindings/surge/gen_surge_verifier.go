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

// SurgeVerifierInternalVerifier is an auto generated low-level Go binding around an user-defined struct.
type SurgeVerifierInternalVerifier struct {
	Addr                common.Address
	AllowInstantUpgrade bool
}

// SurgeVerifierMetaData contains all meta data concerning the SurgeVerifier contract.
var SurgeVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_inbox\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_numProofsThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"RISC0_RETH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SP1_RETH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ZISK_RETH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getHistoricalMarking\",\"inputs\":[{\"name\":\"_proofBitFlag\",\"type\":\"uint8\",\"internalType\":\"LibProofBitmap.ProofBitmap\"},{\"name\":\"_proposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[{\"name\":\"markedForUpgrade_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getInternalVerifier\",\"inputs\":[{\"name\":\"_proofBitFlag\",\"type\":\"uint8\",\"internalType\":\"LibProofBitmap.ProofBitmap\"}],\"outputs\":[{\"name\":\"verifier_\",\"type\":\"tuple\",\"internalType\":\"structSurgeVerifier.InternalVerifier\",\"components\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowInstantUpgrade\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inbox\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"markVerifiersUpgradeable\",\"inputs\":[{\"name\":\"_proposalId\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"_proofBitmap\",\"type\":\"uint8\",\"internalType\":\"LibProofBitmap.ProofBitmap\"},{\"name\":\"_allowInstantUpgrade\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"numProofsThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setVerifier\",\"inputs\":[{\"name\":\"_proofBitFlag\",\"type\":\"uint8\",\"internalType\":\"LibProofBitmap.ProofBitmap\"},{\"name\":\"_verifierAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setVerifierInstant\",\"inputs\":[{\"name\":\"_proofBitFlag\",\"type\":\"uint8\",\"internalType\":\"LibProofBitmap.ProofBitmap\"},{\"name\":\"_verifierAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyProof\",\"inputs\":[{\"name\":\"_requiresThreshold\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"_transitionsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"mergedBitmap_\",\"type\":\"uint8\",\"internalType\":\"LibProofBitmap.ProofBitmap\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VerifierUpdated\",\"inputs\":[{\"name\":\"proofBitFlag\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"LibProofBitmap.ProofBitmap\"},{\"name\":\"oldVerifier\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newVerifier\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"Surge_AlreadyMarkedForProposalId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Surge_CallerIsNotInbox\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Surge_InstantUpgradeNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Surge_InvalidProofBitFlag\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Surge_NumProofsThresholdNotMet\",\"inputs\":[]}]",
}

// SurgeVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use SurgeVerifierMetaData.ABI instead.
var SurgeVerifierABI = SurgeVerifierMetaData.ABI

// SurgeVerifier is an auto generated Go binding around an Ethereum contract.
type SurgeVerifier struct {
	SurgeVerifierCaller     // Read-only binding to the contract
	SurgeVerifierTransactor // Write-only binding to the contract
	SurgeVerifierFilterer   // Log filterer for contract events
}

// SurgeVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type SurgeVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SurgeVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SurgeVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SurgeVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SurgeVerifierSession struct {
	Contract     *SurgeVerifier    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SurgeVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SurgeVerifierCallerSession struct {
	Contract *SurgeVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SurgeVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SurgeVerifierTransactorSession struct {
	Contract     *SurgeVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SurgeVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type SurgeVerifierRaw struct {
	Contract *SurgeVerifier // Generic contract binding to access the raw methods on
}

// SurgeVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SurgeVerifierCallerRaw struct {
	Contract *SurgeVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// SurgeVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SurgeVerifierTransactorRaw struct {
	Contract *SurgeVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSurgeVerifier creates a new instance of SurgeVerifier, bound to a specific deployed contract.
func NewSurgeVerifier(address common.Address, backend bind.ContractBackend) (*SurgeVerifier, error) {
	contract, err := bindSurgeVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SurgeVerifier{SurgeVerifierCaller: SurgeVerifierCaller{contract: contract}, SurgeVerifierTransactor: SurgeVerifierTransactor{contract: contract}, SurgeVerifierFilterer: SurgeVerifierFilterer{contract: contract}}, nil
}

// NewSurgeVerifierCaller creates a new read-only instance of SurgeVerifier, bound to a specific deployed contract.
func NewSurgeVerifierCaller(address common.Address, caller bind.ContractCaller) (*SurgeVerifierCaller, error) {
	contract, err := bindSurgeVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SurgeVerifierCaller{contract: contract}, nil
}

// NewSurgeVerifierTransactor creates a new write-only instance of SurgeVerifier, bound to a specific deployed contract.
func NewSurgeVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*SurgeVerifierTransactor, error) {
	contract, err := bindSurgeVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SurgeVerifierTransactor{contract: contract}, nil
}

// NewSurgeVerifierFilterer creates a new log filterer instance of SurgeVerifier, bound to a specific deployed contract.
func NewSurgeVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*SurgeVerifierFilterer, error) {
	contract, err := bindSurgeVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SurgeVerifierFilterer{contract: contract}, nil
}

// bindSurgeVerifier binds a generic wrapper to an already deployed contract.
func bindSurgeVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SurgeVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SurgeVerifier *SurgeVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SurgeVerifier.Contract.SurgeVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SurgeVerifier *SurgeVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.SurgeVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SurgeVerifier *SurgeVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.SurgeVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SurgeVerifier *SurgeVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SurgeVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SurgeVerifier *SurgeVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SurgeVerifier *SurgeVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.contract.Transact(opts, method, params...)
}

// RISC0RETH is a free data retrieval call binding the contract method 0x8a153cb3.
//
// Solidity: function RISC0_RETH() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierCaller) RISC0RETH(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "RISC0_RETH")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// RISC0RETH is a free data retrieval call binding the contract method 0x8a153cb3.
//
// Solidity: function RISC0_RETH() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierSession) RISC0RETH() (uint8, error) {
	return _SurgeVerifier.Contract.RISC0RETH(&_SurgeVerifier.CallOpts)
}

// RISC0RETH is a free data retrieval call binding the contract method 0x8a153cb3.
//
// Solidity: function RISC0_RETH() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierCallerSession) RISC0RETH() (uint8, error) {
	return _SurgeVerifier.Contract.RISC0RETH(&_SurgeVerifier.CallOpts)
}

// SP1RETH is a free data retrieval call binding the contract method 0x89b2e184.
//
// Solidity: function SP1_RETH() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierCaller) SP1RETH(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "SP1_RETH")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SP1RETH is a free data retrieval call binding the contract method 0x89b2e184.
//
// Solidity: function SP1_RETH() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierSession) SP1RETH() (uint8, error) {
	return _SurgeVerifier.Contract.SP1RETH(&_SurgeVerifier.CallOpts)
}

// SP1RETH is a free data retrieval call binding the contract method 0x89b2e184.
//
// Solidity: function SP1_RETH() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierCallerSession) SP1RETH() (uint8, error) {
	return _SurgeVerifier.Contract.SP1RETH(&_SurgeVerifier.CallOpts)
}

// ZISKRETH is a free data retrieval call binding the contract method 0xca889f91.
//
// Solidity: function ZISK_RETH() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierCaller) ZISKRETH(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "ZISK_RETH")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// ZISKRETH is a free data retrieval call binding the contract method 0xca889f91.
//
// Solidity: function ZISK_RETH() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierSession) ZISKRETH() (uint8, error) {
	return _SurgeVerifier.Contract.ZISKRETH(&_SurgeVerifier.CallOpts)
}

// ZISKRETH is a free data retrieval call binding the contract method 0xca889f91.
//
// Solidity: function ZISK_RETH() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierCallerSession) ZISKRETH() (uint8, error) {
	return _SurgeVerifier.Contract.ZISKRETH(&_SurgeVerifier.CallOpts)
}

// GetHistoricalMarking is a free data retrieval call binding the contract method 0x12a90ac6.
//
// Solidity: function getHistoricalMarking(uint8 _proofBitFlag, uint48 _proposalId) view returns(bool markedForUpgrade_)
func (_SurgeVerifier *SurgeVerifierCaller) GetHistoricalMarking(opts *bind.CallOpts, _proofBitFlag uint8, _proposalId *big.Int) (bool, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "getHistoricalMarking", _proofBitFlag, _proposalId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetHistoricalMarking is a free data retrieval call binding the contract method 0x12a90ac6.
//
// Solidity: function getHistoricalMarking(uint8 _proofBitFlag, uint48 _proposalId) view returns(bool markedForUpgrade_)
func (_SurgeVerifier *SurgeVerifierSession) GetHistoricalMarking(_proofBitFlag uint8, _proposalId *big.Int) (bool, error) {
	return _SurgeVerifier.Contract.GetHistoricalMarking(&_SurgeVerifier.CallOpts, _proofBitFlag, _proposalId)
}

// GetHistoricalMarking is a free data retrieval call binding the contract method 0x12a90ac6.
//
// Solidity: function getHistoricalMarking(uint8 _proofBitFlag, uint48 _proposalId) view returns(bool markedForUpgrade_)
func (_SurgeVerifier *SurgeVerifierCallerSession) GetHistoricalMarking(_proofBitFlag uint8, _proposalId *big.Int) (bool, error) {
	return _SurgeVerifier.Contract.GetHistoricalMarking(&_SurgeVerifier.CallOpts, _proofBitFlag, _proposalId)
}

// GetInternalVerifier is a free data retrieval call binding the contract method 0xe8ff17b0.
//
// Solidity: function getInternalVerifier(uint8 _proofBitFlag) view returns((address,bool) verifier_)
func (_SurgeVerifier *SurgeVerifierCaller) GetInternalVerifier(opts *bind.CallOpts, _proofBitFlag uint8) (SurgeVerifierInternalVerifier, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "getInternalVerifier", _proofBitFlag)

	if err != nil {
		return *new(SurgeVerifierInternalVerifier), err
	}

	out0 := *abi.ConvertType(out[0], new(SurgeVerifierInternalVerifier)).(*SurgeVerifierInternalVerifier)

	return out0, err

}

// GetInternalVerifier is a free data retrieval call binding the contract method 0xe8ff17b0.
//
// Solidity: function getInternalVerifier(uint8 _proofBitFlag) view returns((address,bool) verifier_)
func (_SurgeVerifier *SurgeVerifierSession) GetInternalVerifier(_proofBitFlag uint8) (SurgeVerifierInternalVerifier, error) {
	return _SurgeVerifier.Contract.GetInternalVerifier(&_SurgeVerifier.CallOpts, _proofBitFlag)
}

// GetInternalVerifier is a free data retrieval call binding the contract method 0xe8ff17b0.
//
// Solidity: function getInternalVerifier(uint8 _proofBitFlag) view returns((address,bool) verifier_)
func (_SurgeVerifier *SurgeVerifierCallerSession) GetInternalVerifier(_proofBitFlag uint8) (SurgeVerifierInternalVerifier, error) {
	return _SurgeVerifier.Contract.GetInternalVerifier(&_SurgeVerifier.CallOpts, _proofBitFlag)
}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_SurgeVerifier *SurgeVerifierCaller) Inbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "inbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_SurgeVerifier *SurgeVerifierSession) Inbox() (common.Address, error) {
	return _SurgeVerifier.Contract.Inbox(&_SurgeVerifier.CallOpts)
}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_SurgeVerifier *SurgeVerifierCallerSession) Inbox() (common.Address, error) {
	return _SurgeVerifier.Contract.Inbox(&_SurgeVerifier.CallOpts)
}

// NumProofsThreshold is a free data retrieval call binding the contract method 0xf69e2330.
//
// Solidity: function numProofsThreshold() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierCaller) NumProofsThreshold(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "numProofsThreshold")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// NumProofsThreshold is a free data retrieval call binding the contract method 0xf69e2330.
//
// Solidity: function numProofsThreshold() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierSession) NumProofsThreshold() (uint8, error) {
	return _SurgeVerifier.Contract.NumProofsThreshold(&_SurgeVerifier.CallOpts)
}

// NumProofsThreshold is a free data retrieval call binding the contract method 0xf69e2330.
//
// Solidity: function numProofsThreshold() view returns(uint8)
func (_SurgeVerifier *SurgeVerifierCallerSession) NumProofsThreshold() (uint8, error) {
	return _SurgeVerifier.Contract.NumProofsThreshold(&_SurgeVerifier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SurgeVerifier *SurgeVerifierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SurgeVerifier *SurgeVerifierSession) Owner() (common.Address, error) {
	return _SurgeVerifier.Contract.Owner(&_SurgeVerifier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SurgeVerifier *SurgeVerifierCallerSession) Owner() (common.Address, error) {
	return _SurgeVerifier.Contract.Owner(&_SurgeVerifier.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_SurgeVerifier *SurgeVerifierCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_SurgeVerifier *SurgeVerifierSession) PendingOwner() (common.Address, error) {
	return _SurgeVerifier.Contract.PendingOwner(&_SurgeVerifier.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_SurgeVerifier *SurgeVerifierCallerSession) PendingOwner() (common.Address, error) {
	return _SurgeVerifier.Contract.PendingOwner(&_SurgeVerifier.CallOpts)
}

// VerifyProof is a free data retrieval call binding the contract method 0x0d0eb2f7.
//
// Solidity: function verifyProof(bool _requiresThreshold, bytes32 _transitionsHash, bytes _proof) view returns(uint8 mergedBitmap_)
func (_SurgeVerifier *SurgeVerifierCaller) VerifyProof(opts *bind.CallOpts, _requiresThreshold bool, _transitionsHash [32]byte, _proof []byte) (uint8, error) {
	var out []interface{}
	err := _SurgeVerifier.contract.Call(opts, &out, "verifyProof", _requiresThreshold, _transitionsHash, _proof)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// VerifyProof is a free data retrieval call binding the contract method 0x0d0eb2f7.
//
// Solidity: function verifyProof(bool _requiresThreshold, bytes32 _transitionsHash, bytes _proof) view returns(uint8 mergedBitmap_)
func (_SurgeVerifier *SurgeVerifierSession) VerifyProof(_requiresThreshold bool, _transitionsHash [32]byte, _proof []byte) (uint8, error) {
	return _SurgeVerifier.Contract.VerifyProof(&_SurgeVerifier.CallOpts, _requiresThreshold, _transitionsHash, _proof)
}

// VerifyProof is a free data retrieval call binding the contract method 0x0d0eb2f7.
//
// Solidity: function verifyProof(bool _requiresThreshold, bytes32 _transitionsHash, bytes _proof) view returns(uint8 mergedBitmap_)
func (_SurgeVerifier *SurgeVerifierCallerSession) VerifyProof(_requiresThreshold bool, _transitionsHash [32]byte, _proof []byte) (uint8, error) {
	return _SurgeVerifier.Contract.VerifyProof(&_SurgeVerifier.CallOpts, _requiresThreshold, _transitionsHash, _proof)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SurgeVerifier *SurgeVerifierTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeVerifier.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SurgeVerifier *SurgeVerifierSession) AcceptOwnership() (*types.Transaction, error) {
	return _SurgeVerifier.Contract.AcceptOwnership(&_SurgeVerifier.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SurgeVerifier *SurgeVerifierTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SurgeVerifier.Contract.AcceptOwnership(&_SurgeVerifier.TransactOpts)
}

// MarkVerifiersUpgradeable is a paid mutator transaction binding the contract method 0xb3fd159e.
//
// Solidity: function markVerifiersUpgradeable(uint48 _proposalId, uint8 _proofBitmap, bool _allowInstantUpgrade) returns()
func (_SurgeVerifier *SurgeVerifierTransactor) MarkVerifiersUpgradeable(opts *bind.TransactOpts, _proposalId *big.Int, _proofBitmap uint8, _allowInstantUpgrade bool) (*types.Transaction, error) {
	return _SurgeVerifier.contract.Transact(opts, "markVerifiersUpgradeable", _proposalId, _proofBitmap, _allowInstantUpgrade)
}

// MarkVerifiersUpgradeable is a paid mutator transaction binding the contract method 0xb3fd159e.
//
// Solidity: function markVerifiersUpgradeable(uint48 _proposalId, uint8 _proofBitmap, bool _allowInstantUpgrade) returns()
func (_SurgeVerifier *SurgeVerifierSession) MarkVerifiersUpgradeable(_proposalId *big.Int, _proofBitmap uint8, _allowInstantUpgrade bool) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.MarkVerifiersUpgradeable(&_SurgeVerifier.TransactOpts, _proposalId, _proofBitmap, _allowInstantUpgrade)
}

// MarkVerifiersUpgradeable is a paid mutator transaction binding the contract method 0xb3fd159e.
//
// Solidity: function markVerifiersUpgradeable(uint48 _proposalId, uint8 _proofBitmap, bool _allowInstantUpgrade) returns()
func (_SurgeVerifier *SurgeVerifierTransactorSession) MarkVerifiersUpgradeable(_proposalId *big.Int, _proofBitmap uint8, _allowInstantUpgrade bool) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.MarkVerifiersUpgradeable(&_SurgeVerifier.TransactOpts, _proposalId, _proofBitmap, _allowInstantUpgrade)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SurgeVerifier *SurgeVerifierTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SurgeVerifier.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SurgeVerifier *SurgeVerifierSession) RenounceOwnership() (*types.Transaction, error) {
	return _SurgeVerifier.Contract.RenounceOwnership(&_SurgeVerifier.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SurgeVerifier *SurgeVerifierTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SurgeVerifier.Contract.RenounceOwnership(&_SurgeVerifier.TransactOpts)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x2fb8ff97.
//
// Solidity: function setVerifier(uint8 _proofBitFlag, address _verifierAddr) returns()
func (_SurgeVerifier *SurgeVerifierTransactor) SetVerifier(opts *bind.TransactOpts, _proofBitFlag uint8, _verifierAddr common.Address) (*types.Transaction, error) {
	return _SurgeVerifier.contract.Transact(opts, "setVerifier", _proofBitFlag, _verifierAddr)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x2fb8ff97.
//
// Solidity: function setVerifier(uint8 _proofBitFlag, address _verifierAddr) returns()
func (_SurgeVerifier *SurgeVerifierSession) SetVerifier(_proofBitFlag uint8, _verifierAddr common.Address) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.SetVerifier(&_SurgeVerifier.TransactOpts, _proofBitFlag, _verifierAddr)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x2fb8ff97.
//
// Solidity: function setVerifier(uint8 _proofBitFlag, address _verifierAddr) returns()
func (_SurgeVerifier *SurgeVerifierTransactorSession) SetVerifier(_proofBitFlag uint8, _verifierAddr common.Address) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.SetVerifier(&_SurgeVerifier.TransactOpts, _proofBitFlag, _verifierAddr)
}

// SetVerifierInstant is a paid mutator transaction binding the contract method 0x3bdee25a.
//
// Solidity: function setVerifierInstant(uint8 _proofBitFlag, address _verifierAddr) returns()
func (_SurgeVerifier *SurgeVerifierTransactor) SetVerifierInstant(opts *bind.TransactOpts, _proofBitFlag uint8, _verifierAddr common.Address) (*types.Transaction, error) {
	return _SurgeVerifier.contract.Transact(opts, "setVerifierInstant", _proofBitFlag, _verifierAddr)
}

// SetVerifierInstant is a paid mutator transaction binding the contract method 0x3bdee25a.
//
// Solidity: function setVerifierInstant(uint8 _proofBitFlag, address _verifierAddr) returns()
func (_SurgeVerifier *SurgeVerifierSession) SetVerifierInstant(_proofBitFlag uint8, _verifierAddr common.Address) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.SetVerifierInstant(&_SurgeVerifier.TransactOpts, _proofBitFlag, _verifierAddr)
}

// SetVerifierInstant is a paid mutator transaction binding the contract method 0x3bdee25a.
//
// Solidity: function setVerifierInstant(uint8 _proofBitFlag, address _verifierAddr) returns()
func (_SurgeVerifier *SurgeVerifierTransactorSession) SetVerifierInstant(_proofBitFlag uint8, _verifierAddr common.Address) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.SetVerifierInstant(&_SurgeVerifier.TransactOpts, _proofBitFlag, _verifierAddr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SurgeVerifier *SurgeVerifierTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SurgeVerifier.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SurgeVerifier *SurgeVerifierSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.TransferOwnership(&_SurgeVerifier.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SurgeVerifier *SurgeVerifierTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SurgeVerifier.Contract.TransferOwnership(&_SurgeVerifier.TransactOpts, newOwner)
}

// SurgeVerifierOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the SurgeVerifier contract.
type SurgeVerifierOwnershipTransferStartedIterator struct {
	Event *SurgeVerifierOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *SurgeVerifierOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeVerifierOwnershipTransferStarted)
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
		it.Event = new(SurgeVerifierOwnershipTransferStarted)
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
func (it *SurgeVerifierOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeVerifierOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeVerifierOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the SurgeVerifier contract.
type SurgeVerifierOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_SurgeVerifier *SurgeVerifierFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SurgeVerifierOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeVerifier.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SurgeVerifierOwnershipTransferStartedIterator{contract: _SurgeVerifier.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_SurgeVerifier *SurgeVerifierFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *SurgeVerifierOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeVerifier.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeVerifierOwnershipTransferStarted)
				if err := _SurgeVerifier.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_SurgeVerifier *SurgeVerifierFilterer) ParseOwnershipTransferStarted(log types.Log) (*SurgeVerifierOwnershipTransferStarted, error) {
	event := new(SurgeVerifierOwnershipTransferStarted)
	if err := _SurgeVerifier.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeVerifierOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SurgeVerifier contract.
type SurgeVerifierOwnershipTransferredIterator struct {
	Event *SurgeVerifierOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SurgeVerifierOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeVerifierOwnershipTransferred)
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
		it.Event = new(SurgeVerifierOwnershipTransferred)
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
func (it *SurgeVerifierOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeVerifierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeVerifierOwnershipTransferred represents a OwnershipTransferred event raised by the SurgeVerifier contract.
type SurgeVerifierOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SurgeVerifier *SurgeVerifierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SurgeVerifierOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeVerifier.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SurgeVerifierOwnershipTransferredIterator{contract: _SurgeVerifier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SurgeVerifier *SurgeVerifierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SurgeVerifierOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SurgeVerifier.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeVerifierOwnershipTransferred)
				if err := _SurgeVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SurgeVerifier *SurgeVerifierFilterer) ParseOwnershipTransferred(log types.Log) (*SurgeVerifierOwnershipTransferred, error) {
	event := new(SurgeVerifierOwnershipTransferred)
	if err := _SurgeVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SurgeVerifierVerifierUpdatedIterator is returned from FilterVerifierUpdated and is used to iterate over the raw logs and unpacked data for VerifierUpdated events raised by the SurgeVerifier contract.
type SurgeVerifierVerifierUpdatedIterator struct {
	Event *SurgeVerifierVerifierUpdated // Event containing the contract specifics and raw log

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
func (it *SurgeVerifierVerifierUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SurgeVerifierVerifierUpdated)
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
		it.Event = new(SurgeVerifierVerifierUpdated)
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
func (it *SurgeVerifierVerifierUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SurgeVerifierVerifierUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SurgeVerifierVerifierUpdated represents a VerifierUpdated event raised by the SurgeVerifier contract.
type SurgeVerifierVerifierUpdated struct {
	ProofBitFlag uint8
	OldVerifier  common.Address
	NewVerifier  common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterVerifierUpdated is a free log retrieval operation binding the contract event 0xed3e33a70ec4417f1f470c4efd7e210336fa1b2ac4182d8aa0b6b034897583d3.
//
// Solidity: event VerifierUpdated(uint8 indexed proofBitFlag, address oldVerifier, address newVerifier)
func (_SurgeVerifier *SurgeVerifierFilterer) FilterVerifierUpdated(opts *bind.FilterOpts, proofBitFlag []uint8) (*SurgeVerifierVerifierUpdatedIterator, error) {

	var proofBitFlagRule []interface{}
	for _, proofBitFlagItem := range proofBitFlag {
		proofBitFlagRule = append(proofBitFlagRule, proofBitFlagItem)
	}

	logs, sub, err := _SurgeVerifier.contract.FilterLogs(opts, "VerifierUpdated", proofBitFlagRule)
	if err != nil {
		return nil, err
	}
	return &SurgeVerifierVerifierUpdatedIterator{contract: _SurgeVerifier.contract, event: "VerifierUpdated", logs: logs, sub: sub}, nil
}

// WatchVerifierUpdated is a free log subscription operation binding the contract event 0xed3e33a70ec4417f1f470c4efd7e210336fa1b2ac4182d8aa0b6b034897583d3.
//
// Solidity: event VerifierUpdated(uint8 indexed proofBitFlag, address oldVerifier, address newVerifier)
func (_SurgeVerifier *SurgeVerifierFilterer) WatchVerifierUpdated(opts *bind.WatchOpts, sink chan<- *SurgeVerifierVerifierUpdated, proofBitFlag []uint8) (event.Subscription, error) {

	var proofBitFlagRule []interface{}
	for _, proofBitFlagItem := range proofBitFlag {
		proofBitFlagRule = append(proofBitFlagRule, proofBitFlagItem)
	}

	logs, sub, err := _SurgeVerifier.contract.WatchLogs(opts, "VerifierUpdated", proofBitFlagRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SurgeVerifierVerifierUpdated)
				if err := _SurgeVerifier.contract.UnpackLog(event, "VerifierUpdated", log); err != nil {
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

// ParseVerifierUpdated is a log parse operation binding the contract event 0xed3e33a70ec4417f1f470c4efd7e210336fa1b2ac4182d8aa0b6b034897583d3.
//
// Solidity: event VerifierUpdated(uint8 indexed proofBitFlag, address oldVerifier, address newVerifier)
func (_SurgeVerifier *SurgeVerifierFilterer) ParseVerifierUpdated(log types.Log) (*SurgeVerifierVerifierUpdated, error) {
	event := new(SurgeVerifierVerifierUpdated)
	if err := _SurgeVerifier.contract.UnpackLog(event, "VerifierUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
