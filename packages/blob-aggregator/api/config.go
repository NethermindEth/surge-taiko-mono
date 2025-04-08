package api

import (
	"strings"

	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/queue"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/queue/rabbitmq"
	"github.com/urfave/cli/v2"
)

type Config struct {
	CORSOrigins   []string
	HTTPPort      uint64
	OpenQueueFunc func() (queue.Queue, error)
}

// NewConfigFromCliContext creates a new config instance from command line flags.
func NewConfigFromCliContext(c *cli.Context) (*Config, error) {
	return &Config{
		CORSOrigins: strings.Split(c.String(flags.CORSOrigins.Name), ","),
		HTTPPort:    c.Uint64(flags.HTTPPort.Name),
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
	}, nil
}
