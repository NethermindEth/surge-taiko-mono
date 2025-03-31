package flags

import (
	"github.com/urfave/cli/v2"
)

var (
	L1AggregatorPrivKey = &cli.StringFlag{
		Name:     "l1.aggregatorPrivKey",
		Usage:    "Private key of the L1 aggregator, who will send the L2 proposals with a shared blob",
		Required: true,
		Category: aggregatorCategory,
		EnvVars:  []string{"L1_AGGREGATOR_PRIV_KEY"},
	}
	L1RPCUrl = &cli.StringFlag{
		Name:     "l1.rpcUrl",
		Usage:    "RPC URL of the L1 chain",
		Required: true,
		Category: aggregatorCategory,
		EnvVars:  []string{"L1_RPC_URL"},
	}
	MinAggregatedBlobs = &cli.Uint64Flag{
		Name:     "minAggregatedBlobs",
		Usage:    "Min number of blobs to aggregate block proposals across",
		Required: false,
		Value:    3,
		Category: aggregatorCategory,
		EnvVars:  []string{"MIN_AGGREGATED_BLOBS"},
	}
	MinBlobsFillupPercentage = &cli.Uint64Flag{
		Name:     "minBlobsFillupPercentage",
		Usage:    "Minimum fillup percentage of the aggregated blob space",
		Required: false,
		Value:    75,
		Category: aggregatorCategory,
		EnvVars:  []string{"MIN_BLOB_FILLUP_PERCENTAGE"},
	}
)

var AggregatorFlags = MergeFlags(CommonFlags, TxmgrFlags, []cli.Flag{
	L1RPCUrl,
	L1AggregatorPrivKey,
	MinAggregatedBlobs,
	MinBlobsFillupPercentage,
})
