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

const blockGasLimit = 60000000
const blockGasTarget = blockGasLimit / 2
const txGasLimit = 21000

type SyntheticBlockGenerator struct {
	blockTime  time.Duration
	lastBlock  *types.Header
	accounts   []*ecdsa.PrivateKey // Store multiple private keys
	initialKey *ecdsa.PrivateKey   // Store initial key
	nonce      uint64
}

func NewSyntheticBlockGenerator(blockTime time.Duration, numAccounts int, initialKey *ecdsa.PrivateKey, initialNonce uint64) *SyntheticBlockGenerator {
	accounts := make([]*ecdsa.PrivateKey, numAccounts)

	return &SyntheticBlockGenerator{
		blockTime:  blockTime,
		accounts:   accounts,
		initialKey: initialKey,
		nonce:      initialNonce,
	}
}

// Helper function to create a self-transfer transaction
func createSelfTransferTx(nonce uint64, privateKey *ecdsa.PrivateKey) []byte {
	account := crypto.PubkeyToAddress(privateKey.PublicKey)
	tx := types.NewTransaction(
		nonce,         // nonce
		account,       // to (same as sender)
		big.NewInt(0), // value (0 ETH)
		txGasLimit,    // gas limit (standard transfer)
		big.NewInt(1), // gas price (1 wei)
		nil,           // data
	)

	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privateKey)
	txBytes, _ := signedTx.MarshalBinary()
	return txBytes
}

// Helper function to create a transfer to next account
func createTransferToNextTx(nonce uint64, fromKey *ecdsa.PrivateKey, toKey *ecdsa.PrivateKey, value *big.Int) []byte {
	toAddr := crypto.PubkeyToAddress(toKey.PublicKey)
	tx := types.NewTransaction(
		nonce,         // nonce (will be 127)
		toAddr,        // to (next account)
		value,         // transfer amount
		txGasLimit,    // gas limit
		big.NewInt(1), // gas price (1 wei)
		nil,           // data
	)
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), fromKey)
	txBytes, _ := signedTx.MarshalBinary()
	return txBytes
}

func (g *SyntheticBlockGenerator) generateTransactions() [][]byte {
	// Generate accounts
	for i := 0; i < len(g.accounts); i++ {
		key, _ := crypto.GenerateKey()
		g.accounts[i] = key
	}

	var transactions [][]byte
	transferAmount := big.NewInt(1e17) // 0.1 ETH

	availableGas := blockGasTarget - txGasLimit*2

	// initial funding transfer
	lastAccount := g.accounts[0]
	transactions = append(transactions, createTransferToNextTx(g.nonce, g.initialKey, lastAccount, transferAmount))
	g.nonce++

	lastNonce := uint64(0)
	i := 0
	for i, lastAccount = range g.accounts {
		// Generate 126 self-transfers (nonce 0-126)
		for ; lastNonce < 127; lastNonce++ {
			if availableGas-txGasLimit < 0 {
				break
			}
			transactions = append(transactions, createSelfTransferTx(lastNonce, lastAccount))
		}

		transferAmount.Sub(transferAmount, big.NewInt(int64(lastNonce+1)*21000))

		// Transfer to next account with nonce 127
		if availableGas-txGasLimit < 0 {
			break
		}
		nextIdx := (i + 1) % len(g.accounts)
		transactions = append(transactions, createTransferToNextTx(lastNonce, lastAccount, g.accounts[nextIdx], transferAmount))
		lastNonce = 0
	}

	// Transfer remaining back to initial account
	transactions = append(transactions, createTransferToNextTx(lastNonce, g.accounts[(i)%len(g.accounts)], lastAccount, transferAmount))

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
		GasLimit:      blockGasLimit,
		GasUsed:       txGasLimit * uint64(len(transactions)), // Update gas used (21000 per transaction)
		Timestamp:     timestamp,
		ExtraData:     []byte{},
		BaseFeePerGas: big.NewInt(1), // 1 wei
		BlockHash:     common.Hash{},
		Transactions:  transactions,
	}
}

// Usage example:
// initialKey, _ := crypto.HexToECDSA("your_private_key_hex")
// generator := NewSyntheticBlockGenerator(time.Second, 8, initialKey)
