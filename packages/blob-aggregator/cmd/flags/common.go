package flags

import (
	"github.com/urfave/cli/v2"
)

var (
	commonCategory     = "COMMON"
	aggregatorCategory = "AGGREGATOR"
	txmgrCategory      = "TXMGR"
)

var (
	QueueUsername = &cli.StringFlag{
		Name:     "queue.username",
		Usage:    "Queue connection username",
		Required: true,
		Category: commonCategory,
		EnvVars:  []string{"QUEUE_USER"},
	}
	QueuePassword = &cli.StringFlag{
		Name:     "queue.password",
		Usage:    "Queue connection password",
		Required: true,
		Category: commonCategory,
		EnvVars:  []string{"QUEUE_PASSWORD"},
	}
	QueueHost = &cli.StringFlag{
		Name:     "queue.host",
		Usage:    "Queue connection host",
		Required: true,
		Category: commonCategory,
		EnvVars:  []string{"QUEUE_HOST"},
	}
	QueuePort = &cli.Uint64Flag{
		Name:     "queue.port",
		Usage:    "Queue connection port",
		Required: true,
		Category: commonCategory,
		EnvVars:  []string{"QUEUE_PORT"},
	}
)

// All common flags.
var CommonFlags = []cli.Flag{
	QueueUsername,
	QueuePassword,
	QueueHost,
	QueuePort,
}

// MergeFlags merges the given flag slices.
func MergeFlags(groups ...[]cli.Flag) []cli.Flag {
	var merged []cli.Flag
	for _, group := range groups {
		merged = append(merged, group...)
	}

	return merged
}
