package flags

import (
	"github.com/urfave/cli/v2"
)

var (
	L1AggregatorPrivKey = &cli.StringFlag{
		Name:     "l1.proposerPrivKey",
		Usage:    "Private key of the L1 aggregator, who will send the L2 proposals with a shared blob",
		Required: true,
		Category: aggregatorCategory,
		EnvVars:  []string{"L1_AGGREGATOR_PRIV_KEY"},
	}
)
