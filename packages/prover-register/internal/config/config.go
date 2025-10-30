package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	VerifierAddress  string
	VerifierType     string
	ProverAddress    string
	RPCURL           string
	PrivateKey       string
	TrustCollateral  bool
	RegisterInstance bool
}

func LoadEnv(filepath string) error {
	if filepath == "" {
		filepath = ".env"
	}

	if _, err := os.Stat(filepath); err == nil {
		return godotenv.Load(filepath)
	}

	return nil
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ValidateConfig(cfg *Config) error {
	if cfg.VerifierAddress == "" {
		return fmt.Errorf("verifier address is required")
	}

	if cfg.VerifierType == "" {
		return fmt.Errorf("verifier type is required")
	}

	validTypes := map[string]bool{
		"sgx":       true,
		"sp1":       true,
		"risc0":     true,
		"tdx":       true,
		"azure-tdx": true,
	}

	if !validTypes[cfg.VerifierType] {
		return fmt.Errorf("invalid verifier type: %s", cfg.VerifierType)
	}

	if cfg.ProverAddress == "" {
		return fmt.Errorf("prover address is required")
	}

	if cfg.RPCURL == "" {
		return fmt.Errorf("RPC URL is required")
	}

	if cfg.PrivateKey == "" {
		return fmt.Errorf("private key is required")
	}

	if !cfg.TrustCollateral && !cfg.RegisterInstance {
		return fmt.Errorf("at least one of trust or register must be specified")
	}

	return nil
}
