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

// BondManagerMetaData contains all meta data concerning the BondManager contract.
var BondManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"cancelWithdrawal\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositTo\",\"inputs\":[{\"name\":\"_recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getBond\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"bond_\",\"type\":\"tuple\",\"internalType\":\"structIBondManager.Bond\",\"components\":[{\"name\":\"balance\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"withdrawalRequestedAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestWithdrawal\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BondDeposited\",\"inputs\":[{\"name\":\"depositor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BondWithdrawn\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LivenessBondSettled\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"livenessBond\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"credited\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"slashed\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalCancelled\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalRequested\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawableAt\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false}]",
}

// BondManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use BondManagerMetaData.ABI instead.
var BondManagerABI = BondManagerMetaData.ABI

// BondManager is an auto generated Go binding around an Ethereum contract.
type BondManager struct {
	BondManagerCaller     // Read-only binding to the contract
	BondManagerTransactor // Write-only binding to the contract
	BondManagerFilterer   // Log filterer for contract events
}

// BondManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type BondManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BondManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BondManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BondManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BondManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BondManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BondManagerSession struct {
	Contract     *BondManager      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BondManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BondManagerCallerSession struct {
	Contract *BondManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BondManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BondManagerTransactorSession struct {
	Contract     *BondManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BondManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type BondManagerRaw struct {
	Contract *BondManager // Generic contract binding to access the raw methods on
}

// BondManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BondManagerCallerRaw struct {
	Contract *BondManagerCaller // Generic read-only contract binding to access the raw methods on
}

// BondManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BondManagerTransactorRaw struct {
	Contract *BondManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBondManager creates a new instance of BondManager, bound to a specific deployed contract.
func NewBondManager(address common.Address, backend bind.ContractBackend) (*BondManager, error) {
	contract, err := bindBondManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BondManager{BondManagerCaller: BondManagerCaller{contract: contract}, BondManagerTransactor: BondManagerTransactor{contract: contract}, BondManagerFilterer: BondManagerFilterer{contract: contract}}, nil
}

// NewBondManagerCaller creates a new read-only instance of BondManager, bound to a specific deployed contract.
func NewBondManagerCaller(address common.Address, caller bind.ContractCaller) (*BondManagerCaller, error) {
	contract, err := bindBondManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BondManagerCaller{contract: contract}, nil
}

// NewBondManagerTransactor creates a new write-only instance of BondManager, bound to a specific deployed contract.
func NewBondManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*BondManagerTransactor, error) {
	contract, err := bindBondManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BondManagerTransactor{contract: contract}, nil
}

// NewBondManagerFilterer creates a new log filterer instance of BondManager, bound to a specific deployed contract.
func NewBondManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*BondManagerFilterer, error) {
	contract, err := bindBondManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BondManagerFilterer{contract: contract}, nil
}

// bindBondManager binds a generic wrapper to an already deployed contract.
func bindBondManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BondManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BondManager *BondManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BondManager.Contract.BondManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BondManager *BondManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BondManager.Contract.BondManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BondManager *BondManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BondManager.Contract.BondManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BondManager *BondManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BondManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BondManager *BondManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BondManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BondManager *BondManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BondManager.Contract.contract.Transact(opts, method, params...)
}

// GetBond is a free data retrieval call binding the contract method 0x0d8912f3.
//
// Solidity: function getBond(address _address) view returns((uint64,uint48) bond_)
func (_BondManager *BondManagerCaller) GetBond(opts *bind.CallOpts, _address common.Address) (IBondManagerBond, error) {
	var out []interface{}
	err := _BondManager.contract.Call(opts, &out, "getBond", _address)

	if err != nil {
		return *new(IBondManagerBond), err
	}

	out0 := *abi.ConvertType(out[0], new(IBondManagerBond)).(*IBondManagerBond)

	return out0, err

}

// GetBond is a free data retrieval call binding the contract method 0x0d8912f3.
//
// Solidity: function getBond(address _address) view returns((uint64,uint48) bond_)
func (_BondManager *BondManagerSession) GetBond(_address common.Address) (IBondManagerBond, error) {
	return _BondManager.Contract.GetBond(&_BondManager.CallOpts, _address)
}

// GetBond is a free data retrieval call binding the contract method 0x0d8912f3.
//
// Solidity: function getBond(address _address) view returns((uint64,uint48) bond_)
func (_BondManager *BondManagerCallerSession) GetBond(_address common.Address) (IBondManagerBond, error) {
	return _BondManager.Contract.GetBond(&_BondManager.CallOpts, _address)
}

// CancelWithdrawal is a paid mutator transaction binding the contract method 0x22611280.
//
// Solidity: function cancelWithdrawal() returns()
func (_BondManager *BondManagerTransactor) CancelWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BondManager.contract.Transact(opts, "cancelWithdrawal")
}

// CancelWithdrawal is a paid mutator transaction binding the contract method 0x22611280.
//
// Solidity: function cancelWithdrawal() returns()
func (_BondManager *BondManagerSession) CancelWithdrawal() (*types.Transaction, error) {
	return _BondManager.Contract.CancelWithdrawal(&_BondManager.TransactOpts)
}

// CancelWithdrawal is a paid mutator transaction binding the contract method 0x22611280.
//
// Solidity: function cancelWithdrawal() returns()
func (_BondManager *BondManagerTransactorSession) CancelWithdrawal() (*types.Transaction, error) {
	return _BondManager.Contract.CancelWithdrawal(&_BondManager.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x13765838.
//
// Solidity: function deposit(uint64 _amount) payable returns()
func (_BondManager *BondManagerTransactor) Deposit(opts *bind.TransactOpts, _amount uint64) (*types.Transaction, error) {
	return _BondManager.contract.Transact(opts, "deposit", _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x13765838.
//
// Solidity: function deposit(uint64 _amount) payable returns()
func (_BondManager *BondManagerSession) Deposit(_amount uint64) (*types.Transaction, error) {
	return _BondManager.Contract.Deposit(&_BondManager.TransactOpts, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x13765838.
//
// Solidity: function deposit(uint64 _amount) payable returns()
func (_BondManager *BondManagerTransactorSession) Deposit(_amount uint64) (*types.Transaction, error) {
	return _BondManager.Contract.Deposit(&_BondManager.TransactOpts, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xefba83c9.
//
// Solidity: function depositTo(address _recipient, uint64 _amount) payable returns()
func (_BondManager *BondManagerTransactor) DepositTo(opts *bind.TransactOpts, _recipient common.Address, _amount uint64) (*types.Transaction, error) {
	return _BondManager.contract.Transact(opts, "depositTo", _recipient, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xefba83c9.
//
// Solidity: function depositTo(address _recipient, uint64 _amount) payable returns()
func (_BondManager *BondManagerSession) DepositTo(_recipient common.Address, _amount uint64) (*types.Transaction, error) {
	return _BondManager.Contract.DepositTo(&_BondManager.TransactOpts, _recipient, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xefba83c9.
//
// Solidity: function depositTo(address _recipient, uint64 _amount) payable returns()
func (_BondManager *BondManagerTransactorSession) DepositTo(_recipient common.Address, _amount uint64) (*types.Transaction, error) {
	return _BondManager.Contract.DepositTo(&_BondManager.TransactOpts, _recipient, _amount)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0xdbaf2145.
//
// Solidity: function requestWithdrawal() returns()
func (_BondManager *BondManagerTransactor) RequestWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BondManager.contract.Transact(opts, "requestWithdrawal")
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0xdbaf2145.
//
// Solidity: function requestWithdrawal() returns()
func (_BondManager *BondManagerSession) RequestWithdrawal() (*types.Transaction, error) {
	return _BondManager.Contract.RequestWithdrawal(&_BondManager.TransactOpts)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0xdbaf2145.
//
// Solidity: function requestWithdrawal() returns()
func (_BondManager *BondManagerTransactorSession) RequestWithdrawal() (*types.Transaction, error) {
	return _BondManager.Contract.RequestWithdrawal(&_BondManager.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd6dad060.
//
// Solidity: function withdraw(address _to, uint64 _amount) returns()
func (_BondManager *BondManagerTransactor) Withdraw(opts *bind.TransactOpts, _to common.Address, _amount uint64) (*types.Transaction, error) {
	return _BondManager.contract.Transact(opts, "withdraw", _to, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd6dad060.
//
// Solidity: function withdraw(address _to, uint64 _amount) returns()
func (_BondManager *BondManagerSession) Withdraw(_to common.Address, _amount uint64) (*types.Transaction, error) {
	return _BondManager.Contract.Withdraw(&_BondManager.TransactOpts, _to, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd6dad060.
//
// Solidity: function withdraw(address _to, uint64 _amount) returns()
func (_BondManager *BondManagerTransactorSession) Withdraw(_to common.Address, _amount uint64) (*types.Transaction, error) {
	return _BondManager.Contract.Withdraw(&_BondManager.TransactOpts, _to, _amount)
}

// BondManagerBondDepositedIterator is returned from FilterBondDeposited and is used to iterate over the raw logs and unpacked data for BondDeposited events raised by the BondManager contract.
type BondManagerBondDepositedIterator struct {
	Event *BondManagerBondDeposited // Event containing the contract specifics and raw log

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
func (it *BondManagerBondDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondManagerBondDeposited)
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
		it.Event = new(BondManagerBondDeposited)
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
func (it *BondManagerBondDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondManagerBondDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondManagerBondDeposited represents a BondDeposited event raised by the BondManager contract.
type BondManagerBondDeposited struct {
	Depositor common.Address
	Recipient common.Address
	Amount    uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBondDeposited is a free log retrieval operation binding the contract event 0xe5e95641fa87bdfef3ce0d39f0c9a37c200f3bf59f53623b3de21e03ed33e3d2.
//
// Solidity: event BondDeposited(address indexed depositor, address indexed recipient, uint64 amount)
func (_BondManager *BondManagerFilterer) FilterBondDeposited(opts *bind.FilterOpts, depositor []common.Address, recipient []common.Address) (*BondManagerBondDepositedIterator, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BondManager.contract.FilterLogs(opts, "BondDeposited", depositorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BondManagerBondDepositedIterator{contract: _BondManager.contract, event: "BondDeposited", logs: logs, sub: sub}, nil
}

// WatchBondDeposited is a free log subscription operation binding the contract event 0xe5e95641fa87bdfef3ce0d39f0c9a37c200f3bf59f53623b3de21e03ed33e3d2.
//
// Solidity: event BondDeposited(address indexed depositor, address indexed recipient, uint64 amount)
func (_BondManager *BondManagerFilterer) WatchBondDeposited(opts *bind.WatchOpts, sink chan<- *BondManagerBondDeposited, depositor []common.Address, recipient []common.Address) (event.Subscription, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BondManager.contract.WatchLogs(opts, "BondDeposited", depositorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondManagerBondDeposited)
				if err := _BondManager.contract.UnpackLog(event, "BondDeposited", log); err != nil {
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
func (_BondManager *BondManagerFilterer) ParseBondDeposited(log types.Log) (*BondManagerBondDeposited, error) {
	event := new(BondManagerBondDeposited)
	if err := _BondManager.contract.UnpackLog(event, "BondDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondManagerBondWithdrawnIterator is returned from FilterBondWithdrawn and is used to iterate over the raw logs and unpacked data for BondWithdrawn events raised by the BondManager contract.
type BondManagerBondWithdrawnIterator struct {
	Event *BondManagerBondWithdrawn // Event containing the contract specifics and raw log

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
func (it *BondManagerBondWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondManagerBondWithdrawn)
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
		it.Event = new(BondManagerBondWithdrawn)
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
func (it *BondManagerBondWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondManagerBondWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondManagerBondWithdrawn represents a BondWithdrawn event raised by the BondManager contract.
type BondManagerBondWithdrawn struct {
	Account common.Address
	Amount  uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBondWithdrawn is a free log retrieval operation binding the contract event 0x3362c96009316515fccd3dd29c7036c305ad9e892d83dd5681845ac9edb0c9a8.
//
// Solidity: event BondWithdrawn(address indexed account, uint64 amount)
func (_BondManager *BondManagerFilterer) FilterBondWithdrawn(opts *bind.FilterOpts, account []common.Address) (*BondManagerBondWithdrawnIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _BondManager.contract.FilterLogs(opts, "BondWithdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return &BondManagerBondWithdrawnIterator{contract: _BondManager.contract, event: "BondWithdrawn", logs: logs, sub: sub}, nil
}

// WatchBondWithdrawn is a free log subscription operation binding the contract event 0x3362c96009316515fccd3dd29c7036c305ad9e892d83dd5681845ac9edb0c9a8.
//
// Solidity: event BondWithdrawn(address indexed account, uint64 amount)
func (_BondManager *BondManagerFilterer) WatchBondWithdrawn(opts *bind.WatchOpts, sink chan<- *BondManagerBondWithdrawn, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _BondManager.contract.WatchLogs(opts, "BondWithdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondManagerBondWithdrawn)
				if err := _BondManager.contract.UnpackLog(event, "BondWithdrawn", log); err != nil {
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
func (_BondManager *BondManagerFilterer) ParseBondWithdrawn(log types.Log) (*BondManagerBondWithdrawn, error) {
	event := new(BondManagerBondWithdrawn)
	if err := _BondManager.contract.UnpackLog(event, "BondWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondManagerLivenessBondSettledIterator is returned from FilterLivenessBondSettled and is used to iterate over the raw logs and unpacked data for LivenessBondSettled events raised by the BondManager contract.
type BondManagerLivenessBondSettledIterator struct {
	Event *BondManagerLivenessBondSettled // Event containing the contract specifics and raw log

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
func (it *BondManagerLivenessBondSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondManagerLivenessBondSettled)
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
		it.Event = new(BondManagerLivenessBondSettled)
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
func (it *BondManagerLivenessBondSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondManagerLivenessBondSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondManagerLivenessBondSettled represents a LivenessBondSettled event raised by the BondManager contract.
type BondManagerLivenessBondSettled struct {
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
func (_BondManager *BondManagerFilterer) FilterLivenessBondSettled(opts *bind.FilterOpts, payer []common.Address, payee []common.Address) (*BondManagerLivenessBondSettledIterator, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	logs, sub, err := _BondManager.contract.FilterLogs(opts, "LivenessBondSettled", payerRule, payeeRule)
	if err != nil {
		return nil, err
	}
	return &BondManagerLivenessBondSettledIterator{contract: _BondManager.contract, event: "LivenessBondSettled", logs: logs, sub: sub}, nil
}

// WatchLivenessBondSettled is a free log subscription operation binding the contract event 0xaa22f5157944b5fa6846460e159d57ea9c3878e71fda274af372fa2ccf285aa0.
//
// Solidity: event LivenessBondSettled(address indexed payer, address indexed payee, uint64 livenessBond, uint64 credited, uint64 slashed)
func (_BondManager *BondManagerFilterer) WatchLivenessBondSettled(opts *bind.WatchOpts, sink chan<- *BondManagerLivenessBondSettled, payer []common.Address, payee []common.Address) (event.Subscription, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	logs, sub, err := _BondManager.contract.WatchLogs(opts, "LivenessBondSettled", payerRule, payeeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondManagerLivenessBondSettled)
				if err := _BondManager.contract.UnpackLog(event, "LivenessBondSettled", log); err != nil {
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
func (_BondManager *BondManagerFilterer) ParseLivenessBondSettled(log types.Log) (*BondManagerLivenessBondSettled, error) {
	event := new(BondManagerLivenessBondSettled)
	if err := _BondManager.contract.UnpackLog(event, "LivenessBondSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondManagerWithdrawalCancelledIterator is returned from FilterWithdrawalCancelled and is used to iterate over the raw logs and unpacked data for WithdrawalCancelled events raised by the BondManager contract.
type BondManagerWithdrawalCancelledIterator struct {
	Event *BondManagerWithdrawalCancelled // Event containing the contract specifics and raw log

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
func (it *BondManagerWithdrawalCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondManagerWithdrawalCancelled)
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
		it.Event = new(BondManagerWithdrawalCancelled)
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
func (it *BondManagerWithdrawalCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondManagerWithdrawalCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondManagerWithdrawalCancelled represents a WithdrawalCancelled event raised by the BondManager contract.
type BondManagerWithdrawalCancelled struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalCancelled is a free log retrieval operation binding the contract event 0xc51fdb96728de385ec7859819e3997bc618362ef0dbca0ad051d856866cda3db.
//
// Solidity: event WithdrawalCancelled(address indexed account)
func (_BondManager *BondManagerFilterer) FilterWithdrawalCancelled(opts *bind.FilterOpts, account []common.Address) (*BondManagerWithdrawalCancelledIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _BondManager.contract.FilterLogs(opts, "WithdrawalCancelled", accountRule)
	if err != nil {
		return nil, err
	}
	return &BondManagerWithdrawalCancelledIterator{contract: _BondManager.contract, event: "WithdrawalCancelled", logs: logs, sub: sub}, nil
}

// WatchWithdrawalCancelled is a free log subscription operation binding the contract event 0xc51fdb96728de385ec7859819e3997bc618362ef0dbca0ad051d856866cda3db.
//
// Solidity: event WithdrawalCancelled(address indexed account)
func (_BondManager *BondManagerFilterer) WatchWithdrawalCancelled(opts *bind.WatchOpts, sink chan<- *BondManagerWithdrawalCancelled, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _BondManager.contract.WatchLogs(opts, "WithdrawalCancelled", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondManagerWithdrawalCancelled)
				if err := _BondManager.contract.UnpackLog(event, "WithdrawalCancelled", log); err != nil {
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
func (_BondManager *BondManagerFilterer) ParseWithdrawalCancelled(log types.Log) (*BondManagerWithdrawalCancelled, error) {
	event := new(BondManagerWithdrawalCancelled)
	if err := _BondManager.contract.UnpackLog(event, "WithdrawalCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondManagerWithdrawalRequestedIterator is returned from FilterWithdrawalRequested and is used to iterate over the raw logs and unpacked data for WithdrawalRequested events raised by the BondManager contract.
type BondManagerWithdrawalRequestedIterator struct {
	Event *BondManagerWithdrawalRequested // Event containing the contract specifics and raw log

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
func (it *BondManagerWithdrawalRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondManagerWithdrawalRequested)
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
		it.Event = new(BondManagerWithdrawalRequested)
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
func (it *BondManagerWithdrawalRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondManagerWithdrawalRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondManagerWithdrawalRequested represents a WithdrawalRequested event raised by the BondManager contract.
type BondManagerWithdrawalRequested struct {
	Account        common.Address
	WithdrawableAt *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalRequested is a free log retrieval operation binding the contract event 0x3bbe41cfdd142e0f9b2224dac18c6efd2a6966e35a9ec23ab57ce63a60b33604.
//
// Solidity: event WithdrawalRequested(address indexed account, uint48 withdrawableAt)
func (_BondManager *BondManagerFilterer) FilterWithdrawalRequested(opts *bind.FilterOpts, account []common.Address) (*BondManagerWithdrawalRequestedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _BondManager.contract.FilterLogs(opts, "WithdrawalRequested", accountRule)
	if err != nil {
		return nil, err
	}
	return &BondManagerWithdrawalRequestedIterator{contract: _BondManager.contract, event: "WithdrawalRequested", logs: logs, sub: sub}, nil
}

// WatchWithdrawalRequested is a free log subscription operation binding the contract event 0x3bbe41cfdd142e0f9b2224dac18c6efd2a6966e35a9ec23ab57ce63a60b33604.
//
// Solidity: event WithdrawalRequested(address indexed account, uint48 withdrawableAt)
func (_BondManager *BondManagerFilterer) WatchWithdrawalRequested(opts *bind.WatchOpts, sink chan<- *BondManagerWithdrawalRequested, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _BondManager.contract.WatchLogs(opts, "WithdrawalRequested", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondManagerWithdrawalRequested)
				if err := _BondManager.contract.UnpackLog(event, "WithdrawalRequested", log); err != nil {
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
func (_BondManager *BondManagerFilterer) ParseWithdrawalRequested(log types.Log) (*BondManagerWithdrawalRequested, error) {
	event := new(BondManagerWithdrawalRequested)
	if err := _BondManager.contract.UnpackLog(event, "WithdrawalRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
