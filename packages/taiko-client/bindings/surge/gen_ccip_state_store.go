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

// AzureTDXAkPub is an auto generated low-level Go binding around an user-defined struct.
type AzureTDXAkPub struct {
	ExponentRaw *big.Int
	ModulusRaw  []byte
}

// AzureTDXAttestation is an auto generated low-level Go binding around an user-defined struct.
type AzureTDXAttestation struct {
	TpmQuote AzureTDXTPMQuote
}

// AzureTDXAttestationDocument is an auto generated low-level Go binding around an user-defined struct.
type AzureTDXAttestationDocument struct {
	Attestation  AzureTDXAttestation
	InstanceInfo AzureTDXInstanceInfo
	UserData     []byte
}

// AzureTDXInstanceInfo is an auto generated low-level Go binding around an user-defined struct.
type AzureTDXInstanceInfo struct {
	AttestationReport []byte
	RuntimeData       AzureTDXRuntimeData
}

// AzureTDXPCR is an auto generated low-level Go binding around an user-defined struct.
type AzureTDXPCR struct {
	Index  *big.Int
	Digest [32]byte
}

// AzureTDXRuntimeData is an auto generated low-level Go binding around an user-defined struct.
type AzureTDXRuntimeData struct {
	Raw      []byte
	HclAkPub AzureTDXAkPub
}

// AzureTDXTPMQuote is an auto generated low-level Go binding around an user-defined struct.
type AzureTDXTPMQuote struct {
	Quote        []byte
	RsaSignature []byte
	Pcrs         [24][32]byte
}

// AzureTDXVerifierTrustedParams is an auto generated low-level Go binding around an user-defined struct.
type AzureTDXVerifierTrustedParams struct {
	TeeTcbSvn [16]byte
	PcrBitmap *big.Int
	MrSeam    []byte
	MrTd      []byte
	Pcrs      [][32]byte
}

// AzureTDXVerifyParams is an auto generated low-level Go binding around an user-defined struct.
type AzureTDXVerifyParams struct {
	AttestationDocument AzureTDXAttestationDocument
	Pcrs                []AzureTDXPCR
	Nonce               []byte
}

// ICCIPStateStoreSyncedState is an auto generated low-level Go binding around an user-defined struct.
type ICCIPStateStoreSyncedState struct {
	SyncedAt  *big.Int
	StateRoot [32]byte
	BlockHash [32]byte
}

// CCIPStateStoreMetaData contains all meta data concerning the CCIPStateStore contract.
var CCIPStateStoreMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_automataDcapAttestation\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MIN_SYNC_DELAY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addInstances\",\"inputs\":[{\"name\":\"_instances\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"automataDcapAttestation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deleteInstances\",\"inputs\":[{\"name\":\"_instances\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getSyncedState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structICCIPStateStore.SyncedState\",\"components\":[{\"name\":\"syncedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"impl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inNonReentrant\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"init\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"instances\",\"inputs\":[{\"name\":\"instanceAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"isRegistered\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isInstanceRegistered\",\"inputs\":[{\"name\":\"_instance\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonceUsed\",\"inputs\":[{\"name\":\"nonceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"isUsed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerInstance\",\"inputs\":[{\"name\":\"_trustedParamsIdx\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_attestation\",\"type\":\"tuple\",\"internalType\":\"structAzureTDX.VerifyParams\",\"components\":[{\"name\":\"attestationDocument\",\"type\":\"tuple\",\"internalType\":\"structAzureTDX.AttestationDocument\",\"components\":[{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAzureTDX.Attestation\",\"components\":[{\"name\":\"tpmQuote\",\"type\":\"tuple\",\"internalType\":\"structAzureTDX.TPMQuote\",\"components\":[{\"name\":\"quote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rsaSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pcrs\",\"type\":\"bytes32[24]\",\"internalType\":\"bytes32[24]\"}]}]},{\"name\":\"instanceInfo\",\"type\":\"tuple\",\"internalType\":\"structAzureTDX.InstanceInfo\",\"components\":[{\"name\":\"attestationReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"runtimeData\",\"type\":\"tuple\",\"internalType\":\"structAzureTDX.RuntimeData\",\"components\":[{\"name\":\"raw\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hclAkPub\",\"type\":\"tuple\",\"internalType\":\"structAzureTDX.AkPub\",\"components\":[{\"name\":\"exponentRaw\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"modulusRaw\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]},{\"name\":\"userData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"pcrs\",\"type\":\"tuple[]\",\"internalType\":\"structAzureTDX.PCR[]\",\"components\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"nonce\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setTrustedParams\",\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_params\",\"type\":\"tuple\",\"internalType\":\"structAzureTDXVerifier.TrustedParams\",\"components\":[{\"name\":\"teeTcbSvn\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"},{\"name\":\"pcrBitmap\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"mrSeam\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mrTd\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pcrs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"syncState\",\"inputs\":[{\"name\":\"_proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"trustedParams\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"teeTcbSvn\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"},{\"name\":\"pcrBitmap\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"mrSeam\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mrTd\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AdminChanged\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BeaconUpgraded\",\"inputs\":[{\"name\":\"beacon\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InstanceAdded\",\"inputs\":[{\"name\":\"instance\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InstanceDeleted\",\"inputs\":[{\"name\":\"instance\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StateSynced\",\"inputs\":[{\"name\":\"blockHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"syncedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TrustedParamsUpdated\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"params\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAzureTDXVerifier.TrustedParams\",\"components\":[{\"name\":\"teeTcbSvn\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"},{\"name\":\"pcrBitmap\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"mrSeam\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mrTd\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pcrs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ACCESS_DENIED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AttestationReportHashMismatch\",\"inputs\":[{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicatePCR\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ExtraDataMismatch\",\"inputs\":[{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"FUNC_NOT_IMPLEMENTED\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"INVALID_PAUSE_STATUS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAttestationType\",\"inputs\":[{\"name\":\"actual\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"expected\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"type\":\"error\",\"name\":\"InvalidExtraDataLength\",\"inputs\":[{\"name\":\"actual\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"expected\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"type\":\"error\",\"name\":\"InvalidHashAlgorithm\",\"inputs\":[{\"name\":\"actual\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"expected\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"type\":\"error\",\"name\":\"InvalidMagicValue\",\"inputs\":[{\"name\":\"actual\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expected\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidPCRBitmap\",\"inputs\":[{\"name\":\"actual\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expected\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidPCRDigestLength\",\"inputs\":[{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidPCRIndex\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidPCRSelectionCount\",\"inputs\":[{\"name\":\"actual\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expected\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidRuntimeData\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PCRDigestMismatch\",\"inputs\":[{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"PCRMismatch\",\"inputs\":[{\"name\":\"entryIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"QuoteTooShort\",\"inputs\":[{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"REENTRANT_CALL\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_InvalidProofLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_InvalidSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_SyncTooFrequent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_TdxAlreadyAttested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_TdxInvalidAttestation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_TdxInvalidInstance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_TdxInvalidMrSeam\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_TdxInvalidMrTd\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_TdxInvalidPcr\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_TdxInvalidTcbSvn\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_TdxInvalidTrustedParams\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SurgeCCIP_TdxInvalidVersionType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UNEXPECTED_OFFSET\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_ADDRESS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZERO_VALUE\",\"inputs\":[]}]",
}

// CCIPStateStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use CCIPStateStoreMetaData.ABI instead.
var CCIPStateStoreABI = CCIPStateStoreMetaData.ABI

// CCIPStateStore is an auto generated Go binding around an Ethereum contract.
type CCIPStateStore struct {
	CCIPStateStoreCaller     // Read-only binding to the contract
	CCIPStateStoreTransactor // Write-only binding to the contract
	CCIPStateStoreFilterer   // Log filterer for contract events
}

// CCIPStateStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type CCIPStateStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CCIPStateStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CCIPStateStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CCIPStateStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CCIPStateStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CCIPStateStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CCIPStateStoreSession struct {
	Contract     *CCIPStateStore   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CCIPStateStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CCIPStateStoreCallerSession struct {
	Contract *CCIPStateStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// CCIPStateStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CCIPStateStoreTransactorSession struct {
	Contract     *CCIPStateStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// CCIPStateStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type CCIPStateStoreRaw struct {
	Contract *CCIPStateStore // Generic contract binding to access the raw methods on
}

// CCIPStateStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CCIPStateStoreCallerRaw struct {
	Contract *CCIPStateStoreCaller // Generic read-only contract binding to access the raw methods on
}

// CCIPStateStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CCIPStateStoreTransactorRaw struct {
	Contract *CCIPStateStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCCIPStateStore creates a new instance of CCIPStateStore, bound to a specific deployed contract.
func NewCCIPStateStore(address common.Address, backend bind.ContractBackend) (*CCIPStateStore, error) {
	contract, err := bindCCIPStateStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStore{CCIPStateStoreCaller: CCIPStateStoreCaller{contract: contract}, CCIPStateStoreTransactor: CCIPStateStoreTransactor{contract: contract}, CCIPStateStoreFilterer: CCIPStateStoreFilterer{contract: contract}}, nil
}

// NewCCIPStateStoreCaller creates a new read-only instance of CCIPStateStore, bound to a specific deployed contract.
func NewCCIPStateStoreCaller(address common.Address, caller bind.ContractCaller) (*CCIPStateStoreCaller, error) {
	contract, err := bindCCIPStateStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreCaller{contract: contract}, nil
}

// NewCCIPStateStoreTransactor creates a new write-only instance of CCIPStateStore, bound to a specific deployed contract.
func NewCCIPStateStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*CCIPStateStoreTransactor, error) {
	contract, err := bindCCIPStateStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreTransactor{contract: contract}, nil
}

// NewCCIPStateStoreFilterer creates a new log filterer instance of CCIPStateStore, bound to a specific deployed contract.
func NewCCIPStateStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*CCIPStateStoreFilterer, error) {
	contract, err := bindCCIPStateStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreFilterer{contract: contract}, nil
}

// bindCCIPStateStore binds a generic wrapper to an already deployed contract.
func bindCCIPStateStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CCIPStateStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CCIPStateStore *CCIPStateStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPStateStore.Contract.CCIPStateStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CCIPStateStore *CCIPStateStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.CCIPStateStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CCIPStateStore *CCIPStateStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.CCIPStateStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CCIPStateStore *CCIPStateStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPStateStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CCIPStateStore *CCIPStateStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CCIPStateStore *CCIPStateStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.contract.Transact(opts, method, params...)
}

// MINSYNCDELAY is a free data retrieval call binding the contract method 0x3c68a96d.
//
// Solidity: function MIN_SYNC_DELAY() view returns(uint256)
func (_CCIPStateStore *CCIPStateStoreCaller) MINSYNCDELAY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "MIN_SYNC_DELAY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINSYNCDELAY is a free data retrieval call binding the contract method 0x3c68a96d.
//
// Solidity: function MIN_SYNC_DELAY() view returns(uint256)
func (_CCIPStateStore *CCIPStateStoreSession) MINSYNCDELAY() (*big.Int, error) {
	return _CCIPStateStore.Contract.MINSYNCDELAY(&_CCIPStateStore.CallOpts)
}

// MINSYNCDELAY is a free data retrieval call binding the contract method 0x3c68a96d.
//
// Solidity: function MIN_SYNC_DELAY() view returns(uint256)
func (_CCIPStateStore *CCIPStateStoreCallerSession) MINSYNCDELAY() (*big.Int, error) {
	return _CCIPStateStore.Contract.MINSYNCDELAY(&_CCIPStateStore.CallOpts)
}

// AutomataDcapAttestation is a free data retrieval call binding the contract method 0x36383dc7.
//
// Solidity: function automataDcapAttestation() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCaller) AutomataDcapAttestation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "automataDcapAttestation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AutomataDcapAttestation is a free data retrieval call binding the contract method 0x36383dc7.
//
// Solidity: function automataDcapAttestation() view returns(address)
func (_CCIPStateStore *CCIPStateStoreSession) AutomataDcapAttestation() (common.Address, error) {
	return _CCIPStateStore.Contract.AutomataDcapAttestation(&_CCIPStateStore.CallOpts)
}

// AutomataDcapAttestation is a free data retrieval call binding the contract method 0x36383dc7.
//
// Solidity: function automataDcapAttestation() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCallerSession) AutomataDcapAttestation() (common.Address, error) {
	return _CCIPStateStore.Contract.AutomataDcapAttestation(&_CCIPStateStore.CallOpts)
}

// GetSyncedState is a free data retrieval call binding the contract method 0x7367cf81.
//
// Solidity: function getSyncedState() view returns((uint256,bytes32,bytes32))
func (_CCIPStateStore *CCIPStateStoreCaller) GetSyncedState(opts *bind.CallOpts) (ICCIPStateStoreSyncedState, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "getSyncedState")

	if err != nil {
		return *new(ICCIPStateStoreSyncedState), err
	}

	out0 := *abi.ConvertType(out[0], new(ICCIPStateStoreSyncedState)).(*ICCIPStateStoreSyncedState)

	return out0, err

}

// GetSyncedState is a free data retrieval call binding the contract method 0x7367cf81.
//
// Solidity: function getSyncedState() view returns((uint256,bytes32,bytes32))
func (_CCIPStateStore *CCIPStateStoreSession) GetSyncedState() (ICCIPStateStoreSyncedState, error) {
	return _CCIPStateStore.Contract.GetSyncedState(&_CCIPStateStore.CallOpts)
}

// GetSyncedState is a free data retrieval call binding the contract method 0x7367cf81.
//
// Solidity: function getSyncedState() view returns((uint256,bytes32,bytes32))
func (_CCIPStateStore *CCIPStateStoreCallerSession) GetSyncedState() (ICCIPStateStoreSyncedState, error) {
	return _CCIPStateStore.Contract.GetSyncedState(&_CCIPStateStore.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCaller) Impl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "impl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_CCIPStateStore *CCIPStateStoreSession) Impl() (common.Address, error) {
	return _CCIPStateStore.Contract.Impl(&_CCIPStateStore.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCallerSession) Impl() (common.Address, error) {
	return _CCIPStateStore.Contract.Impl(&_CCIPStateStore.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_CCIPStateStore *CCIPStateStoreCaller) InNonReentrant(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "inNonReentrant")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_CCIPStateStore *CCIPStateStoreSession) InNonReentrant() (bool, error) {
	return _CCIPStateStore.Contract.InNonReentrant(&_CCIPStateStore.CallOpts)
}

// InNonReentrant is a free data retrieval call binding the contract method 0x3075db56.
//
// Solidity: function inNonReentrant() view returns(bool)
func (_CCIPStateStore *CCIPStateStoreCallerSession) InNonReentrant() (bool, error) {
	return _CCIPStateStore.Contract.InNonReentrant(&_CCIPStateStore.CallOpts)
}

// Instances is a free data retrieval call binding the contract method 0x032bb443.
//
// Solidity: function instances(address instanceAddress) view returns(bool isRegistered)
func (_CCIPStateStore *CCIPStateStoreCaller) Instances(opts *bind.CallOpts, instanceAddress common.Address) (bool, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "instances", instanceAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Instances is a free data retrieval call binding the contract method 0x032bb443.
//
// Solidity: function instances(address instanceAddress) view returns(bool isRegistered)
func (_CCIPStateStore *CCIPStateStoreSession) Instances(instanceAddress common.Address) (bool, error) {
	return _CCIPStateStore.Contract.Instances(&_CCIPStateStore.CallOpts, instanceAddress)
}

// Instances is a free data retrieval call binding the contract method 0x032bb443.
//
// Solidity: function instances(address instanceAddress) view returns(bool isRegistered)
func (_CCIPStateStore *CCIPStateStoreCallerSession) Instances(instanceAddress common.Address) (bool, error) {
	return _CCIPStateStore.Contract.Instances(&_CCIPStateStore.CallOpts, instanceAddress)
}

// IsInstanceRegistered is a free data retrieval call binding the contract method 0x739df5de.
//
// Solidity: function isInstanceRegistered(address _instance) view returns(bool)
func (_CCIPStateStore *CCIPStateStoreCaller) IsInstanceRegistered(opts *bind.CallOpts, _instance common.Address) (bool, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "isInstanceRegistered", _instance)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsInstanceRegistered is a free data retrieval call binding the contract method 0x739df5de.
//
// Solidity: function isInstanceRegistered(address _instance) view returns(bool)
func (_CCIPStateStore *CCIPStateStoreSession) IsInstanceRegistered(_instance common.Address) (bool, error) {
	return _CCIPStateStore.Contract.IsInstanceRegistered(&_CCIPStateStore.CallOpts, _instance)
}

// IsInstanceRegistered is a free data retrieval call binding the contract method 0x739df5de.
//
// Solidity: function isInstanceRegistered(address _instance) view returns(bool)
func (_CCIPStateStore *CCIPStateStoreCallerSession) IsInstanceRegistered(_instance common.Address) (bool, error) {
	return _CCIPStateStore.Contract.IsInstanceRegistered(&_CCIPStateStore.CallOpts, _instance)
}

// NonceUsed is a free data retrieval call binding the contract method 0x61a4422b.
//
// Solidity: function nonceUsed(bytes32 nonceHash) view returns(bool isUsed)
func (_CCIPStateStore *CCIPStateStoreCaller) NonceUsed(opts *bind.CallOpts, nonceHash [32]byte) (bool, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "nonceUsed", nonceHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// NonceUsed is a free data retrieval call binding the contract method 0x61a4422b.
//
// Solidity: function nonceUsed(bytes32 nonceHash) view returns(bool isUsed)
func (_CCIPStateStore *CCIPStateStoreSession) NonceUsed(nonceHash [32]byte) (bool, error) {
	return _CCIPStateStore.Contract.NonceUsed(&_CCIPStateStore.CallOpts, nonceHash)
}

// NonceUsed is a free data retrieval call binding the contract method 0x61a4422b.
//
// Solidity: function nonceUsed(bytes32 nonceHash) view returns(bool isUsed)
func (_CCIPStateStore *CCIPStateStoreCallerSession) NonceUsed(nonceHash [32]byte) (bool, error) {
	return _CCIPStateStore.Contract.NonceUsed(&_CCIPStateStore.CallOpts, nonceHash)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CCIPStateStore *CCIPStateStoreSession) Owner() (common.Address, error) {
	return _CCIPStateStore.Contract.Owner(&_CCIPStateStore.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCallerSession) Owner() (common.Address, error) {
	return _CCIPStateStore.Contract.Owner(&_CCIPStateStore.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CCIPStateStore *CCIPStateStoreCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CCIPStateStore *CCIPStateStoreSession) Paused() (bool, error) {
	return _CCIPStateStore.Contract.Paused(&_CCIPStateStore.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CCIPStateStore *CCIPStateStoreCallerSession) Paused() (bool, error) {
	return _CCIPStateStore.Contract.Paused(&_CCIPStateStore.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CCIPStateStore *CCIPStateStoreSession) PendingOwner() (common.Address, error) {
	return _CCIPStateStore.Contract.PendingOwner(&_CCIPStateStore.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCallerSession) PendingOwner() (common.Address, error) {
	return _CCIPStateStore.Contract.PendingOwner(&_CCIPStateStore.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_CCIPStateStore *CCIPStateStoreCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_CCIPStateStore *CCIPStateStoreSession) ProxiableUUID() ([32]byte, error) {
	return _CCIPStateStore.Contract.ProxiableUUID(&_CCIPStateStore.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_CCIPStateStore *CCIPStateStoreCallerSession) ProxiableUUID() ([32]byte, error) {
	return _CCIPStateStore.Contract.ProxiableUUID(&_CCIPStateStore.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCaller) Resolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "resolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_CCIPStateStore *CCIPStateStoreSession) Resolver() (common.Address, error) {
	return _CCIPStateStore.Contract.Resolver(&_CCIPStateStore.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_CCIPStateStore *CCIPStateStoreCallerSession) Resolver() (common.Address, error) {
	return _CCIPStateStore.Contract.Resolver(&_CCIPStateStore.CallOpts)
}

// TrustedParams is a free data retrieval call binding the contract method 0x6e7afc42.
//
// Solidity: function trustedParams(uint256 index) view returns(bytes16 teeTcbSvn, uint24 pcrBitmap, bytes mrSeam, bytes mrTd)
func (_CCIPStateStore *CCIPStateStoreCaller) TrustedParams(opts *bind.CallOpts, index *big.Int) (struct {
	TeeTcbSvn [16]byte
	PcrBitmap *big.Int
	MrSeam    []byte
	MrTd      []byte
}, error) {
	var out []interface{}
	err := _CCIPStateStore.contract.Call(opts, &out, "trustedParams", index)

	outstruct := new(struct {
		TeeTcbSvn [16]byte
		PcrBitmap *big.Int
		MrSeam    []byte
		MrTd      []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TeeTcbSvn = *abi.ConvertType(out[0], new([16]byte)).(*[16]byte)
	outstruct.PcrBitmap = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.MrSeam = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.MrTd = *abi.ConvertType(out[3], new([]byte)).(*[]byte)

	return *outstruct, err

}

// TrustedParams is a free data retrieval call binding the contract method 0x6e7afc42.
//
// Solidity: function trustedParams(uint256 index) view returns(bytes16 teeTcbSvn, uint24 pcrBitmap, bytes mrSeam, bytes mrTd)
func (_CCIPStateStore *CCIPStateStoreSession) TrustedParams(index *big.Int) (struct {
	TeeTcbSvn [16]byte
	PcrBitmap *big.Int
	MrSeam    []byte
	MrTd      []byte
}, error) {
	return _CCIPStateStore.Contract.TrustedParams(&_CCIPStateStore.CallOpts, index)
}

// TrustedParams is a free data retrieval call binding the contract method 0x6e7afc42.
//
// Solidity: function trustedParams(uint256 index) view returns(bytes16 teeTcbSvn, uint24 pcrBitmap, bytes mrSeam, bytes mrTd)
func (_CCIPStateStore *CCIPStateStoreCallerSession) TrustedParams(index *big.Int) (struct {
	TeeTcbSvn [16]byte
	PcrBitmap *big.Int
	MrSeam    []byte
	MrTd      []byte
}, error) {
	return _CCIPStateStore.Contract.TrustedParams(&_CCIPStateStore.CallOpts, index)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CCIPStateStore *CCIPStateStoreSession) AcceptOwnership() (*types.Transaction, error) {
	return _CCIPStateStore.Contract.AcceptOwnership(&_CCIPStateStore.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CCIPStateStore.Contract.AcceptOwnership(&_CCIPStateStore.TransactOpts)
}

// AddInstances is a paid mutator transaction binding the contract method 0x16107290.
//
// Solidity: function addInstances(address[] _instances) returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) AddInstances(opts *bind.TransactOpts, _instances []common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "addInstances", _instances)
}

// AddInstances is a paid mutator transaction binding the contract method 0x16107290.
//
// Solidity: function addInstances(address[] _instances) returns()
func (_CCIPStateStore *CCIPStateStoreSession) AddInstances(_instances []common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.AddInstances(&_CCIPStateStore.TransactOpts, _instances)
}

// AddInstances is a paid mutator transaction binding the contract method 0x16107290.
//
// Solidity: function addInstances(address[] _instances) returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) AddInstances(_instances []common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.AddInstances(&_CCIPStateStore.TransactOpts, _instances)
}

// DeleteInstances is a paid mutator transaction binding the contract method 0x6647dd29.
//
// Solidity: function deleteInstances(address[] _instances) returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) DeleteInstances(opts *bind.TransactOpts, _instances []common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "deleteInstances", _instances)
}

// DeleteInstances is a paid mutator transaction binding the contract method 0x6647dd29.
//
// Solidity: function deleteInstances(address[] _instances) returns()
func (_CCIPStateStore *CCIPStateStoreSession) DeleteInstances(_instances []common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.DeleteInstances(&_CCIPStateStore.TransactOpts, _instances)
}

// DeleteInstances is a paid mutator transaction binding the contract method 0x6647dd29.
//
// Solidity: function deleteInstances(address[] _instances) returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) DeleteInstances(_instances []common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.DeleteInstances(&_CCIPStateStore.TransactOpts, _instances)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) Init(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "init", _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_CCIPStateStore *CCIPStateStoreSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.Init(&_CCIPStateStore.TransactOpts, _owner)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _owner) returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) Init(_owner common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.Init(&_CCIPStateStore.TransactOpts, _owner)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CCIPStateStore *CCIPStateStoreSession) Pause() (*types.Transaction, error) {
	return _CCIPStateStore.Contract.Pause(&_CCIPStateStore.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) Pause() (*types.Transaction, error) {
	return _CCIPStateStore.Contract.Pause(&_CCIPStateStore.TransactOpts)
}

// RegisterInstance is a paid mutator transaction binding the contract method 0x43541ddd.
//
// Solidity: function registerInstance(uint256 _trustedParamsIdx, ((((bytes,bytes,bytes32[24])),(bytes,(bytes,(uint24,bytes))),bytes),(uint256,bytes32)[],bytes) _attestation) returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) RegisterInstance(opts *bind.TransactOpts, _trustedParamsIdx *big.Int, _attestation AzureTDXVerifyParams) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "registerInstance", _trustedParamsIdx, _attestation)
}

// RegisterInstance is a paid mutator transaction binding the contract method 0x43541ddd.
//
// Solidity: function registerInstance(uint256 _trustedParamsIdx, ((((bytes,bytes,bytes32[24])),(bytes,(bytes,(uint24,bytes))),bytes),(uint256,bytes32)[],bytes) _attestation) returns()
func (_CCIPStateStore *CCIPStateStoreSession) RegisterInstance(_trustedParamsIdx *big.Int, _attestation AzureTDXVerifyParams) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.RegisterInstance(&_CCIPStateStore.TransactOpts, _trustedParamsIdx, _attestation)
}

// RegisterInstance is a paid mutator transaction binding the contract method 0x43541ddd.
//
// Solidity: function registerInstance(uint256 _trustedParamsIdx, ((((bytes,bytes,bytes32[24])),(bytes,(bytes,(uint24,bytes))),bytes),(uint256,bytes32)[],bytes) _attestation) returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) RegisterInstance(_trustedParamsIdx *big.Int, _attestation AzureTDXVerifyParams) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.RegisterInstance(&_CCIPStateStore.TransactOpts, _trustedParamsIdx, _attestation)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CCIPStateStore *CCIPStateStoreSession) RenounceOwnership() (*types.Transaction, error) {
	return _CCIPStateStore.Contract.RenounceOwnership(&_CCIPStateStore.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CCIPStateStore.Contract.RenounceOwnership(&_CCIPStateStore.TransactOpts)
}

// SetTrustedParams is a paid mutator transaction binding the contract method 0xd659c16b.
//
// Solidity: function setTrustedParams(uint256 _index, (bytes16,uint24,bytes,bytes,bytes32[]) _params) returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) SetTrustedParams(opts *bind.TransactOpts, _index *big.Int, _params AzureTDXVerifierTrustedParams) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "setTrustedParams", _index, _params)
}

// SetTrustedParams is a paid mutator transaction binding the contract method 0xd659c16b.
//
// Solidity: function setTrustedParams(uint256 _index, (bytes16,uint24,bytes,bytes,bytes32[]) _params) returns()
func (_CCIPStateStore *CCIPStateStoreSession) SetTrustedParams(_index *big.Int, _params AzureTDXVerifierTrustedParams) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.SetTrustedParams(&_CCIPStateStore.TransactOpts, _index, _params)
}

// SetTrustedParams is a paid mutator transaction binding the contract method 0xd659c16b.
//
// Solidity: function setTrustedParams(uint256 _index, (bytes16,uint24,bytes,bytes,bytes32[]) _params) returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) SetTrustedParams(_index *big.Int, _params AzureTDXVerifierTrustedParams) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.SetTrustedParams(&_CCIPStateStore.TransactOpts, _index, _params)
}

// SyncState is a paid mutator transaction binding the contract method 0xcb67103d.
//
// Solidity: function syncState(bytes _proof) returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) SyncState(opts *bind.TransactOpts, _proof []byte) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "syncState", _proof)
}

// SyncState is a paid mutator transaction binding the contract method 0xcb67103d.
//
// Solidity: function syncState(bytes _proof) returns()
func (_CCIPStateStore *CCIPStateStoreSession) SyncState(_proof []byte) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.SyncState(&_CCIPStateStore.TransactOpts, _proof)
}

// SyncState is a paid mutator transaction binding the contract method 0xcb67103d.
//
// Solidity: function syncState(bytes _proof) returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) SyncState(_proof []byte) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.SyncState(&_CCIPStateStore.TransactOpts, _proof)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CCIPStateStore *CCIPStateStoreSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.TransferOwnership(&_CCIPStateStore.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.TransferOwnership(&_CCIPStateStore.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CCIPStateStore *CCIPStateStoreSession) Unpause() (*types.Transaction, error) {
	return _CCIPStateStore.Contract.Unpause(&_CCIPStateStore.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) Unpause() (*types.Transaction, error) {
	return _CCIPStateStore.Contract.Unpause(&_CCIPStateStore.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_CCIPStateStore *CCIPStateStoreSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.UpgradeTo(&_CCIPStateStore.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.UpgradeTo(&_CCIPStateStore.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_CCIPStateStore *CCIPStateStoreTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _CCIPStateStore.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_CCIPStateStore *CCIPStateStoreSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.UpgradeToAndCall(&_CCIPStateStore.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_CCIPStateStore *CCIPStateStoreTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _CCIPStateStore.Contract.UpgradeToAndCall(&_CCIPStateStore.TransactOpts, newImplementation, data)
}

// CCIPStateStoreAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the CCIPStateStore contract.
type CCIPStateStoreAdminChangedIterator struct {
	Event *CCIPStateStoreAdminChanged // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreAdminChanged)
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
		it.Event = new(CCIPStateStoreAdminChanged)
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
func (it *CCIPStateStoreAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreAdminChanged represents a AdminChanged event raised by the CCIPStateStore contract.
type CCIPStateStoreAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*CCIPStateStoreAdminChangedIterator, error) {

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreAdminChangedIterator{contract: _CCIPStateStore.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreAdminChanged) (event.Subscription, error) {

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreAdminChanged)
				if err := _CCIPStateStore.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseAdminChanged(log types.Log) (*CCIPStateStoreAdminChanged, error) {
	event := new(CCIPStateStoreAdminChanged)
	if err := _CCIPStateStore.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the CCIPStateStore contract.
type CCIPStateStoreBeaconUpgradedIterator struct {
	Event *CCIPStateStoreBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreBeaconUpgraded)
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
		it.Event = new(CCIPStateStoreBeaconUpgraded)
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
func (it *CCIPStateStoreBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreBeaconUpgraded represents a BeaconUpgraded event raised by the CCIPStateStore contract.
type CCIPStateStoreBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*CCIPStateStoreBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreBeaconUpgradedIterator{contract: _CCIPStateStore.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreBeaconUpgraded)
				if err := _CCIPStateStore.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseBeaconUpgraded(log types.Log) (*CCIPStateStoreBeaconUpgraded, error) {
	event := new(CCIPStateStoreBeaconUpgraded)
	if err := _CCIPStateStore.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the CCIPStateStore contract.
type CCIPStateStoreInitializedIterator struct {
	Event *CCIPStateStoreInitialized // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreInitialized)
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
		it.Event = new(CCIPStateStoreInitialized)
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
func (it *CCIPStateStoreInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreInitialized represents a Initialized event raised by the CCIPStateStore contract.
type CCIPStateStoreInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterInitialized(opts *bind.FilterOpts) (*CCIPStateStoreInitializedIterator, error) {

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreInitializedIterator{contract: _CCIPStateStore.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreInitialized) (event.Subscription, error) {

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreInitialized)
				if err := _CCIPStateStore.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseInitialized(log types.Log) (*CCIPStateStoreInitialized, error) {
	event := new(CCIPStateStoreInitialized)
	if err := _CCIPStateStore.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreInstanceAddedIterator is returned from FilterInstanceAdded and is used to iterate over the raw logs and unpacked data for InstanceAdded events raised by the CCIPStateStore contract.
type CCIPStateStoreInstanceAddedIterator struct {
	Event *CCIPStateStoreInstanceAdded // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreInstanceAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreInstanceAdded)
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
		it.Event = new(CCIPStateStoreInstanceAdded)
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
func (it *CCIPStateStoreInstanceAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreInstanceAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreInstanceAdded represents a InstanceAdded event raised by the CCIPStateStore contract.
type CCIPStateStoreInstanceAdded struct {
	Instance common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterInstanceAdded is a free log retrieval operation binding the contract event 0xee3a98e49d5a27452a99d57c90a7f73d4b2e44de88c6ded02e69c4ed964edd5a.
//
// Solidity: event InstanceAdded(address indexed instance)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterInstanceAdded(opts *bind.FilterOpts, instance []common.Address) (*CCIPStateStoreInstanceAddedIterator, error) {

	var instanceRule []interface{}
	for _, instanceItem := range instance {
		instanceRule = append(instanceRule, instanceItem)
	}

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "InstanceAdded", instanceRule)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreInstanceAddedIterator{contract: _CCIPStateStore.contract, event: "InstanceAdded", logs: logs, sub: sub}, nil
}

// WatchInstanceAdded is a free log subscription operation binding the contract event 0xee3a98e49d5a27452a99d57c90a7f73d4b2e44de88c6ded02e69c4ed964edd5a.
//
// Solidity: event InstanceAdded(address indexed instance)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchInstanceAdded(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreInstanceAdded, instance []common.Address) (event.Subscription, error) {

	var instanceRule []interface{}
	for _, instanceItem := range instance {
		instanceRule = append(instanceRule, instanceItem)
	}

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "InstanceAdded", instanceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreInstanceAdded)
				if err := _CCIPStateStore.contract.UnpackLog(event, "InstanceAdded", log); err != nil {
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

// ParseInstanceAdded is a log parse operation binding the contract event 0xee3a98e49d5a27452a99d57c90a7f73d4b2e44de88c6ded02e69c4ed964edd5a.
//
// Solidity: event InstanceAdded(address indexed instance)
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseInstanceAdded(log types.Log) (*CCIPStateStoreInstanceAdded, error) {
	event := new(CCIPStateStoreInstanceAdded)
	if err := _CCIPStateStore.contract.UnpackLog(event, "InstanceAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreInstanceDeletedIterator is returned from FilterInstanceDeleted and is used to iterate over the raw logs and unpacked data for InstanceDeleted events raised by the CCIPStateStore contract.
type CCIPStateStoreInstanceDeletedIterator struct {
	Event *CCIPStateStoreInstanceDeleted // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreInstanceDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreInstanceDeleted)
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
		it.Event = new(CCIPStateStoreInstanceDeleted)
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
func (it *CCIPStateStoreInstanceDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreInstanceDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreInstanceDeleted represents a InstanceDeleted event raised by the CCIPStateStore contract.
type CCIPStateStoreInstanceDeleted struct {
	Instance common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterInstanceDeleted is a free log retrieval operation binding the contract event 0x68d6dc32f987c030377083eab8a382fd0b43885a3f68afb52386014d2ed994fc.
//
// Solidity: event InstanceDeleted(address indexed instance)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterInstanceDeleted(opts *bind.FilterOpts, instance []common.Address) (*CCIPStateStoreInstanceDeletedIterator, error) {

	var instanceRule []interface{}
	for _, instanceItem := range instance {
		instanceRule = append(instanceRule, instanceItem)
	}

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "InstanceDeleted", instanceRule)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreInstanceDeletedIterator{contract: _CCIPStateStore.contract, event: "InstanceDeleted", logs: logs, sub: sub}, nil
}

// WatchInstanceDeleted is a free log subscription operation binding the contract event 0x68d6dc32f987c030377083eab8a382fd0b43885a3f68afb52386014d2ed994fc.
//
// Solidity: event InstanceDeleted(address indexed instance)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchInstanceDeleted(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreInstanceDeleted, instance []common.Address) (event.Subscription, error) {

	var instanceRule []interface{}
	for _, instanceItem := range instance {
		instanceRule = append(instanceRule, instanceItem)
	}

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "InstanceDeleted", instanceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreInstanceDeleted)
				if err := _CCIPStateStore.contract.UnpackLog(event, "InstanceDeleted", log); err != nil {
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

// ParseInstanceDeleted is a log parse operation binding the contract event 0x68d6dc32f987c030377083eab8a382fd0b43885a3f68afb52386014d2ed994fc.
//
// Solidity: event InstanceDeleted(address indexed instance)
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseInstanceDeleted(log types.Log) (*CCIPStateStoreInstanceDeleted, error) {
	event := new(CCIPStateStoreInstanceDeleted)
	if err := _CCIPStateStore.contract.UnpackLog(event, "InstanceDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the CCIPStateStore contract.
type CCIPStateStoreOwnershipTransferStartedIterator struct {
	Event *CCIPStateStoreOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreOwnershipTransferStarted)
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
		it.Event = new(CCIPStateStoreOwnershipTransferStarted)
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
func (it *CCIPStateStoreOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the CCIPStateStore contract.
type CCIPStateStoreOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CCIPStateStoreOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreOwnershipTransferStartedIterator{contract: _CCIPStateStore.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreOwnershipTransferStarted)
				if err := _CCIPStateStore.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseOwnershipTransferStarted(log types.Log) (*CCIPStateStoreOwnershipTransferStarted, error) {
	event := new(CCIPStateStoreOwnershipTransferStarted)
	if err := _CCIPStateStore.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CCIPStateStore contract.
type CCIPStateStoreOwnershipTransferredIterator struct {
	Event *CCIPStateStoreOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreOwnershipTransferred)
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
		it.Event = new(CCIPStateStoreOwnershipTransferred)
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
func (it *CCIPStateStoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreOwnershipTransferred represents a OwnershipTransferred event raised by the CCIPStateStore contract.
type CCIPStateStoreOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CCIPStateStoreOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreOwnershipTransferredIterator{contract: _CCIPStateStore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreOwnershipTransferred)
				if err := _CCIPStateStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseOwnershipTransferred(log types.Log) (*CCIPStateStoreOwnershipTransferred, error) {
	event := new(CCIPStateStoreOwnershipTransferred)
	if err := _CCIPStateStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStorePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the CCIPStateStore contract.
type CCIPStateStorePausedIterator struct {
	Event *CCIPStateStorePaused // Event containing the contract specifics and raw log

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
func (it *CCIPStateStorePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStorePaused)
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
		it.Event = new(CCIPStateStorePaused)
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
func (it *CCIPStateStorePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStorePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStorePaused represents a Paused event raised by the CCIPStateStore contract.
type CCIPStateStorePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterPaused(opts *bind.FilterOpts) (*CCIPStateStorePausedIterator, error) {

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &CCIPStateStorePausedIterator{contract: _CCIPStateStore.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *CCIPStateStorePaused) (event.Subscription, error) {

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStorePaused)
				if err := _CCIPStateStore.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_CCIPStateStore *CCIPStateStoreFilterer) ParsePaused(log types.Log) (*CCIPStateStorePaused, error) {
	event := new(CCIPStateStorePaused)
	if err := _CCIPStateStore.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreStateSyncedIterator is returned from FilterStateSynced and is used to iterate over the raw logs and unpacked data for StateSynced events raised by the CCIPStateStore contract.
type CCIPStateStoreStateSyncedIterator struct {
	Event *CCIPStateStoreStateSynced // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreStateSyncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreStateSynced)
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
		it.Event = new(CCIPStateStoreStateSynced)
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
func (it *CCIPStateStoreStateSyncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreStateSyncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreStateSynced represents a StateSynced event raised by the CCIPStateStore contract.
type CCIPStateStoreStateSynced struct {
	BlockHash [32]byte
	StateRoot [32]byte
	SyncedAt  *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStateSynced is a free log retrieval operation binding the contract event 0xde1f3aaa9343dffc4bb3c5c898fb567e2ec8ebcfb38134b5a2d0f8fc49e9a5ca.
//
// Solidity: event StateSynced(bytes32 indexed blockHash, bytes32 indexed stateRoot, uint256 syncedAt)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterStateSynced(opts *bind.FilterOpts, blockHash [][32]byte, stateRoot [][32]byte) (*CCIPStateStoreStateSyncedIterator, error) {

	var blockHashRule []interface{}
	for _, blockHashItem := range blockHash {
		blockHashRule = append(blockHashRule, blockHashItem)
	}
	var stateRootRule []interface{}
	for _, stateRootItem := range stateRoot {
		stateRootRule = append(stateRootRule, stateRootItem)
	}

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "StateSynced", blockHashRule, stateRootRule)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreStateSyncedIterator{contract: _CCIPStateStore.contract, event: "StateSynced", logs: logs, sub: sub}, nil
}

// WatchStateSynced is a free log subscription operation binding the contract event 0xde1f3aaa9343dffc4bb3c5c898fb567e2ec8ebcfb38134b5a2d0f8fc49e9a5ca.
//
// Solidity: event StateSynced(bytes32 indexed blockHash, bytes32 indexed stateRoot, uint256 syncedAt)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchStateSynced(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreStateSynced, blockHash [][32]byte, stateRoot [][32]byte) (event.Subscription, error) {

	var blockHashRule []interface{}
	for _, blockHashItem := range blockHash {
		blockHashRule = append(blockHashRule, blockHashItem)
	}
	var stateRootRule []interface{}
	for _, stateRootItem := range stateRoot {
		stateRootRule = append(stateRootRule, stateRootItem)
	}

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "StateSynced", blockHashRule, stateRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreStateSynced)
				if err := _CCIPStateStore.contract.UnpackLog(event, "StateSynced", log); err != nil {
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

// ParseStateSynced is a log parse operation binding the contract event 0xde1f3aaa9343dffc4bb3c5c898fb567e2ec8ebcfb38134b5a2d0f8fc49e9a5ca.
//
// Solidity: event StateSynced(bytes32 indexed blockHash, bytes32 indexed stateRoot, uint256 syncedAt)
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseStateSynced(log types.Log) (*CCIPStateStoreStateSynced, error) {
	event := new(CCIPStateStoreStateSynced)
	if err := _CCIPStateStore.contract.UnpackLog(event, "StateSynced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreTrustedParamsUpdatedIterator is returned from FilterTrustedParamsUpdated and is used to iterate over the raw logs and unpacked data for TrustedParamsUpdated events raised by the CCIPStateStore contract.
type CCIPStateStoreTrustedParamsUpdatedIterator struct {
	Event *CCIPStateStoreTrustedParamsUpdated // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreTrustedParamsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreTrustedParamsUpdated)
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
		it.Event = new(CCIPStateStoreTrustedParamsUpdated)
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
func (it *CCIPStateStoreTrustedParamsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreTrustedParamsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreTrustedParamsUpdated represents a TrustedParamsUpdated event raised by the CCIPStateStore contract.
type CCIPStateStoreTrustedParamsUpdated struct {
	Index  *big.Int
	Params AzureTDXVerifierTrustedParams
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTrustedParamsUpdated is a free log retrieval operation binding the contract event 0xd074d0ae021db90d9b8b14060cb7cb9c28d68a4617ec11137f4c62454a443908.
//
// Solidity: event TrustedParamsUpdated(uint256 indexed index, (bytes16,uint24,bytes,bytes,bytes32[]) params)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterTrustedParamsUpdated(opts *bind.FilterOpts, index []*big.Int) (*CCIPStateStoreTrustedParamsUpdatedIterator, error) {

	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "TrustedParamsUpdated", indexRule)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreTrustedParamsUpdatedIterator{contract: _CCIPStateStore.contract, event: "TrustedParamsUpdated", logs: logs, sub: sub}, nil
}

// WatchTrustedParamsUpdated is a free log subscription operation binding the contract event 0xd074d0ae021db90d9b8b14060cb7cb9c28d68a4617ec11137f4c62454a443908.
//
// Solidity: event TrustedParamsUpdated(uint256 indexed index, (bytes16,uint24,bytes,bytes,bytes32[]) params)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchTrustedParamsUpdated(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreTrustedParamsUpdated, index []*big.Int) (event.Subscription, error) {

	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "TrustedParamsUpdated", indexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreTrustedParamsUpdated)
				if err := _CCIPStateStore.contract.UnpackLog(event, "TrustedParamsUpdated", log); err != nil {
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

// ParseTrustedParamsUpdated is a log parse operation binding the contract event 0xd074d0ae021db90d9b8b14060cb7cb9c28d68a4617ec11137f4c62454a443908.
//
// Solidity: event TrustedParamsUpdated(uint256 indexed index, (bytes16,uint24,bytes,bytes,bytes32[]) params)
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseTrustedParamsUpdated(log types.Log) (*CCIPStateStoreTrustedParamsUpdated, error) {
	event := new(CCIPStateStoreTrustedParamsUpdated)
	if err := _CCIPStateStore.contract.UnpackLog(event, "TrustedParamsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the CCIPStateStore contract.
type CCIPStateStoreUnpausedIterator struct {
	Event *CCIPStateStoreUnpaused // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreUnpaused)
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
		it.Event = new(CCIPStateStoreUnpaused)
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
func (it *CCIPStateStoreUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreUnpaused represents a Unpaused event raised by the CCIPStateStore contract.
type CCIPStateStoreUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterUnpaused(opts *bind.FilterOpts) (*CCIPStateStoreUnpausedIterator, error) {

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreUnpausedIterator{contract: _CCIPStateStore.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreUnpaused) (event.Subscription, error) {

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreUnpaused)
				if err := _CCIPStateStore.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseUnpaused(log types.Log) (*CCIPStateStoreUnpaused, error) {
	event := new(CCIPStateStoreUnpaused)
	if err := _CCIPStateStore.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CCIPStateStoreUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the CCIPStateStore contract.
type CCIPStateStoreUpgradedIterator struct {
	Event *CCIPStateStoreUpgraded // Event containing the contract specifics and raw log

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
func (it *CCIPStateStoreUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPStateStoreUpgraded)
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
		it.Event = new(CCIPStateStoreUpgraded)
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
func (it *CCIPStateStoreUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CCIPStateStoreUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CCIPStateStoreUpgraded represents a Upgraded event raised by the CCIPStateStore contract.
type CCIPStateStoreUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_CCIPStateStore *CCIPStateStoreFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*CCIPStateStoreUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _CCIPStateStore.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &CCIPStateStoreUpgradedIterator{contract: _CCIPStateStore.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_CCIPStateStore *CCIPStateStoreFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *CCIPStateStoreUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _CCIPStateStore.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CCIPStateStoreUpgraded)
				if err := _CCIPStateStore.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_CCIPStateStore *CCIPStateStoreFilterer) ParseUpgraded(log types.Log) (*CCIPStateStoreUpgraded, error) {
	event := new(CCIPStateStoreUpgraded)
	if err := _CCIPStateStore.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
