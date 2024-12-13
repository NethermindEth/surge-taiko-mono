package producer

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
)

func TestCombinedProducerRequestProof(t *testing.T) {
	header := &types.Header{
		ParentHash:  randHash(),
		UncleHash:   randHash(),
		Coinbase:    common.BytesToAddress(randHash().Bytes()),
		Root:        randHash(),
		TxHash:      randHash(),
		ReceiptHash: randHash(),
		Difficulty:  common.Big0,
		Number:      common.Big256,
		GasLimit:    1024,
		GasUsed:     1024,
		Time:        uint64(time.Now().Unix()),
		Extra:       randHash().Bytes(),
		MixDigest:   randHash(),
		Nonce:       types.BlockNonce{},
	}

	optimisticProducer1 := &OptimisticProofProducer{}
	optimisticProducer2 := &OptimisticProofProducer{}

	producer := &CombinedProducer{
		ProofTier:      encoding.TierSgxAndZkVMID,
		RequiredProofs: 2,
		Producers:      []ProofProducer{optimisticProducer1, optimisticProducer2},
		Verifiers: []common.Address{
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			common.HexToAddress("0x0987654321098765432109876543210987654321"),
		},
		ProofStates: make(map[*big.Int]BlockProofState),
	}

	blockID := big.NewInt(1)
	meta := &metadata.TaikoDataBlockMetadataLegacy{}
	opts := &ProofRequestOptions{
		BlockID:       blockID,
		ProverAddress: common.HexToAddress("0x1234"),
	}

	res, err := producer.RequestProof(
		context.Background(),
		opts,
		blockID,
		meta,
		header,
		time.Now(),
	)

	require.Nil(t, err)
	require.Equal(t, blockID, res.BlockID)
	require.Equal(t, header, res.Header)
	require.Equal(t, producer.Tier(), res.Tier)
	require.NotEmpty(t, res.Proof)
}

func TestCombinedProducerRequestCancel(t *testing.T) {
	optimisticProducer1 := &OptimisticProofProducer{}
	optimisticProducer2 := &OptimisticProofProducer{}

	producer := &CombinedProducer{
		ProofTier:      encoding.TierSgxAndZkVMID,
		RequiredProofs: 2,
		Producers:      []ProofProducer{optimisticProducer1, optimisticProducer2},
		Verifiers: []common.Address{
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			common.HexToAddress("0x0987654321098765432109876543210987654321"),
		},
		ProofStates: make(map[*big.Int]BlockProofState),
	}

	opts := &ProofRequestOptions{
		BlockID:       big.NewInt(1),
		ProverAddress: common.HexToAddress("0x1234"),
	}

	err := producer.RequestCancel(context.Background(), opts)
	require.Nil(t, err)
}

func TestCombinedProducerTier(t *testing.T) {
	producer := &CombinedProducer{
		ProofTier: encoding.TierSgxAndZkVMID,
		ProofStates: make(map[*big.Int]BlockProofState),
	}

	require.Equal(t, encoding.TierSgxAndZkVMID, producer.Tier())
}
