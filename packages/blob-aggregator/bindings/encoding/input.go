package encoding

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/bindings"
)

// ABI arguments marshaling components.
var (
	CallComponents = []abi.ArgumentMarshaling{
		{
			Name: "target",
			Type: "address",
		},
		{
			Name: "value",
			Type: "uint256",
		},
		{
			Name: "data",
			Type: "bytes",
		},
	}
)

var (
	CallsComponentsArrayType, _ = abi.NewType("tuple[]", "MinimalBatcher.Call", CallComponents)
	CallsComponentsArrayArgs    = abi.Arguments{
		{Name: "MinimalBatcher.Call[]", Type: CallsComponentsArrayType},
	}
)

// Contract ABIs.
var (
	MinimalBatcherABI *abi.ABI

	customErrorMaps []map[string]abi.Error
)

func init() {
	var err error

	if MinimalBatcherABI, err = bindings.MinimalBatcherMetaData.GetAbi(); err != nil {
		log.Crit("Get MinimalBatcher ABI error", "error", err)
	}

	customErrorMaps = []map[string]abi.Error{
		MinimalBatcherABI.Errors,
	}
}

// EncodeExecuteBatchInput performs the solidity `abi.encode` for the given MinimalBatcher.executeBatch input.
func EncodeExecuteBatchInput(calls []bindings.Call) ([]byte, error) {
	b, err := CallsComponentsArrayArgs.Pack(calls)
	if err != nil {
		return nil, fmt.Errorf("failed to abi.encode MinimalBatcher.executeBatch input, %w", err)
	}
	return b, nil
}
