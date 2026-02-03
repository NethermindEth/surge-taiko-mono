// Extract Azure TDX attestation data for smart contract verification
// Usage: go run cmd/extract/main.go <guest-data-file>
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/formatter"
	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/prover"
	"go.uber.org/zap"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <guest-data-file>\n", os.Args[0])
		os.Exit(1)
	}

	// Create logger
	logger, _ := zap.NewDevelopment()
	log := logger.Sugar()
	defer log.Sync()

	// Read guest data from file
	guestDataFile := os.Args[1]
	proverClient := prover.NewClient("http://dummy", log)
	guestData, err := proverClient.GetGuestDataFromFile(guestDataFile)
	if err != nil {
		log.Fatalf("Failed to read guest data: %v", err)
	}

	// Process with Azure TDX formatter
	azureTdxFormatter := formatter.NewAzureTDXFormatter(log)
	processedData, err := azureTdxFormatter.ProcessGuestData(guestData)
	if err != nil {
		log.Fatalf("Failed to process Azure TDX data: %v", err)
	}

	// Extract trusted params
	trustedParams, err := azureTdxFormatter.ExtractTrustedParams(processedData)
	if err != nil {
		log.Fatalf("Failed to extract trusted params: %v", err)
	}

	// Output as JSON
	output := map[string]interface{}{
		"processedData": processedData,
		"trustedParams": trustedParams,
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal output: %v", err)
	}

	fmt.Println(string(jsonOutput))
}
