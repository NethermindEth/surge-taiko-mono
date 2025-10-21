package cmd

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/contract"
	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/formatter"
	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/logger"
	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/prover"
)

var (
	cfgFile          string
	verifierAddress  string
	verifierType     string
	proverAddress    string
	rpcURL           string
	privateKey       string
	trustCollateral  bool
	registerInstance bool
	envFile          string
	dryRun           bool
	dryRunAsOwner    bool
	logJSON          bool
	logDebug         bool
)

var rootCmd = &cobra.Command{
	Use:   "register",
	Short: "Register prover instances in verifier contracts",
	Long:  `A CLI tool for registering and trusting prover instances in their verifier contracts.`,
	RunE:  runRegister,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&verifierAddress, "verifier", "", "Verifier contract address (required)")
	rootCmd.Flags().StringVar(&verifierType, "type", "", "Verifier type: sgx, sp1, risc0, tdx, or azure-tdx (required)")
	rootCmd.Flags().StringVar(&proverAddress, "prover", "", "Prover address (e.g., http://ip:port) (required)")
	rootCmd.Flags().StringVar(&rpcURL, "rpc", "", "RPC URL for blockchain connection (required)")
	rootCmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for transaction signing")
	rootCmd.Flags().BoolVar(&trustCollateral, "trust", false, "Trust the collateral")
	rootCmd.Flags().BoolVar(&registerInstance, "register", false, "Register the instance")
	rootCmd.Flags().StringVar(&envFile, "env", ".env", "Environment file path")
	rootCmd.Flags().BoolVar(&dryRun, "dry", false, "Dry run mode - simulate transactions without sending")
	rootCmd.Flags().BoolVar(&dryRunAsOwner, "dry-as-owner", false, "Dry run mode simulating as contract owner")
	rootCmd.Flags().BoolVar(&logJSON, "log.json", false, "Output logs in JSON format")
	rootCmd.Flags().BoolVar(&logDebug, "log.debug", false, "Enable debug logging")

	rootCmd.MarkFlagRequired("verifier")
	rootCmd.MarkFlagRequired("type")
	rootCmd.MarkFlagRequired("prover")
	rootCmd.MarkFlagRequired("rpc")
}

func initConfig() {
	if envFile != "" {
		viper.SetConfigFile(envFile)
		viper.SetConfigType("env")
		viper.ReadInConfig()
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if privateKey == "" {
		privateKey = viper.GetString("PRIVATE_KEY")
	}
}

func runRegister(cmd *cobra.Command, args []string) error {
	log, err := logger.NewLogger(logJSON, logDebug)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer log.Sync()
	ctx := context.Background()

	// If --dry-as-owner is set, automatically enable dry run
	if dryRunAsOwner {
		dryRun = true
	}

	log.Infow("starting registration process",
		"verifier", verifierAddress,
		"type", verifierType,
		"prover", proverAddress,
		"dryRun", dryRun,
		"dryRunAsOwner", dryRunAsOwner)

	// Validate inputs
	if !trustCollateral && !registerInstance {
		return fmt.Errorf("at least one of --trust or --register must be specified")
	}

	if privateKey == "" && !dryRun {
		return fmt.Errorf("private key is required (use --private-key or set PRIVATE_KEY in .env)")
	}

	// Parse private key
	var privKey *ecdsa.PrivateKey
	if privateKey != "" {
		var err error
		privKey, err = crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
		if err != nil {
			return fmt.Errorf("invalid private key: %w", err)
		}
	}

	// Connect to blockchain
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	// Create prover client
	proverClient := prover.NewClient(proverAddress, log)

	// Fetch guest data
	log.Info("fetching guest data from prover")
	guestData, err := proverClient.GetGuestData(ctx)
	if err != nil {
		return fmt.Errorf("failed to get guest data: %w", err)
	}

	log.Debugw("received guest data", "data", guestData)

	// Process guest data based on verifier type
	var processedData interface{}
	switch strings.ToLower(verifierType) {
	case "tdx":
		log.Info("processing TDX attestation")
		tdxFormatter := formatter.NewTDXFormatter(log)
		processedData, err = tdxFormatter.ProcessGuestData(guestData)
		if err != nil {
			return fmt.Errorf("failed to process TDX data: %w", err)
		}
	case "azure-tdx":
		log.Info("processing Azure TDX attestation")
		azureTdxFormatter := formatter.NewAzureTDXFormatter(log)
		processedData, err = azureTdxFormatter.ProcessGuestData(guestData)
		if err != nil {
			return fmt.Errorf("failed to process Azure TDX data: %w", err)
		}
	case "sgx":
		log.Info("processing SGX attestation")
		sgxFormatter := formatter.NewSGXFormatter(log)
		processedData, err = sgxFormatter.ProcessGuestData(guestData)
		if err != nil {
			return fmt.Errorf("failed to process SGX data: %w", err)
		}
	case "sp1":
		log.Info("processing SP1 data")
		sp1Formatter := formatter.NewSP1Formatter(log)
		processedData, err = sp1Formatter.ProcessGuestData(guestData)
		if err != nil {
			return fmt.Errorf("failed to process SP1 data: %w", err)
		}
	case "risc0":
		log.Info("processing RISC0 data")
		risc0Formatter := formatter.NewRISC0Formatter(log)
		processedData, err = risc0Formatter.ProcessGuestData(guestData)
		if err != nil {
			return fmt.Errorf("failed to process RISC0 data: %w", err)
		}
	default:
		return fmt.Errorf("unsupported verifier type: %s", verifierType)
	}

	// Create contract client
	verifierAddr := common.HexToAddress(verifierAddress)
	contractClient, err := contract.NewClient(client, verifierAddr, privKey, log, dryRun, dryRunAsOwner)
	if err != nil {
		return fmt.Errorf("failed to create contract client: %w", err)
	}

	// Execute operations
	if trustCollateral {
		log.Info("trusting collateral")
		if err := contractClient.TrustCollateral(ctx, verifierType, processedData); err != nil {
			return fmt.Errorf("failed to trust collateral: %w", err)
		}
		log.Info("collateral trusted successfully")
	}

	if registerInstance {
		log.Info("registering instance")
		if err := contractClient.RegisterInstance(ctx, verifierType, processedData); err != nil {
			return fmt.Errorf("failed to register instance: %w", err)
		}
		log.Info("instance registered successfully")
	}

	// Output final result
	result := map[string]interface{}{
		"success":       true,
		"verifier":      verifierAddress,
		"type":          verifierType,
		"dryRun":        dryRun,
		"dryRunAsOwner": dryRunAsOwner,
		"operations": map[string]bool{
			"trust":    trustCollateral,
			"register": registerInstance,
		},
	}

	resultJSON, _ := json.Marshal(result)
	fmt.Fprintln(os.Stderr, string(resultJSON))

	return nil
}
