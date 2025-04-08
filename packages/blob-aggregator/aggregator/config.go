package aggregator

import (
	"crypto/ecdsa"

	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/cmd/flags"
	pkgFlags "github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/flags"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/queue"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/queue/rabbitmq"
	"github.com/urfave/cli/v2"
)

type Config struct {
	L1RPCUrl                 string
	L1AggregatorPrivKey      *ecdsa.PrivateKey
	MinAggregatedBlobs       uint64
	MinBlobsFillupPercentage uint64
	OpenQueueFunc            func() (queue.Queue, error)
	TxMgrConfig              *txmgr.CLIConfig
}

// NewConfigFromCliContext creates a new config instance from command line flags.
func NewConfigFromCliContext(c *cli.Context) (*Config, error) {
	l1AggregatorPrivKey, err := crypto.ToECDSA(
		common.Hex2Bytes(c.String(flags.L1AggregatorPrivKey.Name)),
	)

	if err != nil {
		return nil, err
	}

	return &Config{
		L1RPCUrl:                 c.String(flags.L1RPCUrl.Name),
		L1AggregatorPrivKey:      l1AggregatorPrivKey,
		MinAggregatedBlobs:       c.Uint64(flags.MinAggregatedBlobs.Name),
		MinBlobsFillupPercentage: c.Uint64(flags.MinBlobsFillupPercentage.Name),
		OpenQueueFunc: func() (queue.Queue, error) {
			opts := queue.NewQueueOpts{
				Username: c.String(flags.QueueUsername.Name),
				Password: c.String(flags.QueuePassword.Name),
				Host:     c.String(flags.QueueHost.Name),
				Port:     c.String(flags.QueuePort.Name),
			}

			q, err := rabbitmq.NewRabbitMQ(opts)
			if err != nil {
				return nil, err
			}

			return q, nil
		},
		TxMgrConfig: pkgFlags.InitTxmgrConfigsFromCli(
			c.String(flags.L1RPCUrl.Name),
			l1AggregatorPrivKey,
			c,
		),
	}, nil
}
