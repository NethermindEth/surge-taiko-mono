package flags

import (
	"time"

	p2pFlags "github.com/ethereum-optimism/optimism/op-node/flags"
	"github.com/urfave/cli/v2"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/fork"
)

// Optional flags used by driver.
var (
	P2PSync = &cli.BoolFlag{
		Name: "p2p.sync",
		Usage: "Try P2P syncing blocks between L2 execution engines, " +
			"will be helpful to bring a new node online quickly",
		Value:    false,
		Category: driverCategory,
		EnvVars:  []string{"P2P_SYNC"},
	}
	P2PSyncTimeout = &cli.DurationFlag{
		Name: "p2p.syncTimeout",
		Usage: "P2P syncing timeout, if no sync progress is made within this time span, " +
			"driver will stop the P2P sync and insert all remaining L2 blocks one by one",
		Value:    1 * time.Hour,
		Category: driverCategory,
		EnvVars:  []string{"P2P_SYNC_TIMEOUT"},
	}
	CheckPointSyncURL = &cli.StringFlag{
		Name:     "p2p.checkPointSyncUrl",
		Usage:    "HTTP RPC endpoint of another synced L2 execution engine node",
		Category: driverCategory,
		EnvVars:  []string{"P2P_CHECK_POINT_SYNC_URL"},
	}
	// blob server endpoint
	BlobServerEndpoint = &cli.StringFlag{
		Name:     "blob.server",
		Usage:    "Blob sidecar storage server",
		Category: driverCategory,
		EnvVars:  []string{"BLOB_SERVER"},
	}
	// preconfirmation block server
	PreconfBlockServerPort = &cli.Uint64Flag{
		Name:     "preconfirmation.serverPort",
		Usage:    "HTTP port of the preconfirmation block server, 0 means disabled",
		Category: driverCategory,
		EnvVars:  []string{"PRECONFIRMATION_SERVER_PORT"},
	}
	PreconfBlockServerJWTSecret = &cli.StringFlag{
		Name:     "preconfirmation.jwtSecret",
		Usage:    "Path to a JWT secret to use for the preconfirmation block server",
		Category: driverCategory,
		EnvVars:  []string{"PRECONFIRMATION_SERVER_JWT_SECRET"},
	}
	PreconfBlockServerCORSOrigins = &cli.StringFlag{
		Name:     "preconfirmation.corsOrigins",
		Usage:    "CORS Origins settings for the preconfirmation block server",
		Category: driverCategory,
		Value:    "*",
		EnvVars:  []string{"PRECONFIRMATION_SERVER_CORS_ORIGINS"},
	}
	PreconfWhitelistAddress = &cli.StringFlag{
		Name:     "preconfirmation.whitelist",
		Usage:    "PreconfWhitelist contract L1 `address`",
		Required: false,
		Category: driverCategory,
		EnvVars:  []string{"PRECONFIRMATION_WHITELIST"},
	}
	DriverTaikoWrapperAddress = &cli.StringFlag{
		Name:     "taikoWrapper",
		Usage:    "TaikoWrapper contract `address`",
		Required: false,
		Category: driverCategory,
		EnvVars:  []string{"TAIKO_WRAPPER"},
	}
	Fork = &cli.StringFlag{
		Name:     "fork",
		Usage:    `Active protocol fork: "pacaya", "shasta", or "realtime"`,
		Value:    fork.RealTime,
		Category: driverCategory,
		EnvVars:  []string{"FORK"},
	}
	PrivacyMode = &cli.BoolFlag{
		Name:     "privacy.mode",
		Usage:    `Enable realtime blob payload privacy. Must match the proposer (Catalyst) and prover (raiko) configuration.`,
		Value:    false,
		Category: driverCategory,
		EnvVars:  []string{"SURGE_PRIVACY_MODE"},
	}
	PrivacySymmetricKey = &cli.StringFlag{
		Name:     "privacy.symmetricKey",
		Usage:    `Hex-encoded 32-byte AES-256-GCM key used to decrypt scheme-0x01 (normal proposal) blobs. Required when --privacy.mode=true.`,
		Category: driverCategory,
		EnvVars:  []string{"SURGE_PRIVACY_SYMMETRIC_KEY"},
	}
	PrivacyFIPrivateKey = &cli.StringFlag{
		Name:     "privacy.fiPrivateKey",
		Usage:    `Hex-encoded 32-byte secp256k1 system FI private key used to decrypt scheme-0x02 (forced inclusion) blobs.`,
		Category: driverCategory,
		EnvVars:  []string{"SURGE_PRIVACY_FI_PRIVKEY"},
	}
)

// DriverFlags All driver flags.
var DriverFlags = MergeFlags(CommonFlags, []cli.Flag{
	L1BeaconEndpoint,
	L2WSEndpoint,
	L2AuthEndpoint,
	JWTSecret,
	P2PSync,
	P2PSyncTimeout,
	CheckPointSyncURL,
	BlobServerEndpoint,
	PreconfBlockServerPort,
	PreconfBlockServerJWTSecret,
	PreconfBlockServerCORSOrigins,
	PreconfWhitelistAddress,
	DriverTaikoWrapperAddress,
	Fork,
	PrivacyMode,
	PrivacySymmetricKey,
	PrivacyFIPrivateKey,
}, p2pFlags.P2PFlags("PRECONFIRMATION"))
