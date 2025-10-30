package main

import (
	"os"

	"github.com/taikoxyz/taiko-mono/packages/prover-register/cmd"
	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/logger"
)

func main() {
	log := logger.NewJSONLogger()
	if err := cmd.Execute(); err != nil {
		log.Error("execution failed", "error", err)
		os.Exit(1)
	}
}
