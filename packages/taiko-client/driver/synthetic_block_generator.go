package driver

import (
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type SyntheticBlockGenerator struct {
	blockTime  time.Duration
	lastBlock  *types.Header
	accounts   []*ecdsa.PrivateKey // Store multiple private keys
	initialKey *ecdsa.PrivateKey   // Store initial key
}

func NewSyntheticBlockGenerator(blockTime time.Duration, numAccounts int, initialKey *ecdsa.PrivateKey) *SyntheticBlockGenerator {
	accounts := make([]*ecdsa.PrivateKey, numAccounts)
	accounts[0] = initialKey // Use provided initial key

	// Generate remaining accounts
	for i := 1; i < numAccounts; i++ {
		key, _ := crypto.GenerateKey()
		accounts[i] = key
	}
	return &SyntheticBlockGenerator{
		blockTime:  blockTime,
		accounts:   accounts,
		initialKey: initialKey,
	}
}

// Helper function to create a self-transfer transaction
func createSelfTransferTx(nonce uint64, privateKey *ecdsa.PrivateKey) []byte {
	account := crypto.PubkeyToAddress(privateKey.PublicKey)
	tx := types.NewTransaction(
		nonce,                  // nonce
		account,                // to (same as sender)
		big.NewInt(0),          // value (0 ETH)
		21000,                  // gas limit (standard transfer)
		big.NewInt(1000000000), // gas price (1 gwei)
		nil,                    // data
	)

	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privateKey)
	txBytes, _ := signedTx.MarshalBinary()
	return txBytes
}

// Helper function to create a transfer to next account
func createTransferToNextTx(nonce uint64, fromKey *ecdsa.PrivateKey, toKey *ecdsa.PrivateKey, value *big.Int) []byte {
	toAddr := crypto.PubkeyToAddress(toKey.PublicKey)
	tx := types.NewTransaction(
		nonce,                  // nonce (will be 127)
		toAddr,                 // to (next account)
		value,                  // transfer amount
		21000,                  // gas limit
		big.NewInt(1000000000), // gas price (1 gwei)
		nil,                    // data
	)
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), fromKey)
	txBytes, _ := signedTx.MarshalBinary()
	return txBytes
}

func (g *SyntheticBlockGenerator) generateTransactions() [][]byte {
	var transactions [][]byte
	transferAmount := big.NewInt(1e18) // 1 ETH

	for i, account := range g.accounts {
		// Generate 126 self-transfers (nonce 0-125)
		for nonce := uint64(0); nonce < 126; nonce++ {
			transactions = append(transactions, createSelfTransferTx(nonce, account))
		}

		// Transfer to next account with nonce 126
		nextIdx := (i + 1) % len(g.accounts)
		transactions = append(transactions, createTransferToNextTx(126, account, g.accounts[nextIdx], transferAmount))
	}

	return transactions
}

func (g *SyntheticBlockGenerator) GenerateBlock(parent *types.Header) *engine.ExecutableData {
	timestamp := uint64(time.Now().Unix())
	if parent != nil && timestamp < parent.Time {
		timestamp = parent.Time + uint64(g.blockTime.Seconds())
	}

	transactions := g.generateTransactions()

	return &engine.ExecutableData{
		ParentHash:    parent.Hash(),
		FeeRecipient:  common.Address{},
		StateRoot:     common.Hash{}, // Empty state root
		ReceiptsRoot:  common.Hash{}, // Empty receipts root
		LogsBloom:     types.Bloom{}.Bytes(),
		Random:        common.Hash{},
		Number:        new(big.Int).Add(parent.Number, common.Big1).Uint64(),
		GasLimit:      30000000,
		GasUsed:       21000 * uint64(len(transactions)), // Update gas used (21000 per transaction)
		Timestamp:     timestamp,
		ExtraData:     []byte{},
		BaseFeePerGas: big.NewInt(1000000000), // 1 gwei
		BlockHash:     common.Hash{},
		Transactions:  transactions,
	}
}

// Usage example:
// initialKey, _ := crypto.HexToECDSA("your_private_key_hex")
// generator := NewSyntheticBlockGenerator(time.Second, 8, initialKey)
