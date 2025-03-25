package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/api"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/cmd/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	log.SetOutput(os.Stdout)
	// attempt to load a .env file to overwrite CLI flags, but allow it to not
	// exist.

	envFile := os.Getenv("BLOB_AGGREGATOR_ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_ = godotenv.Load(envFile)

	app.Name = "Taiko Blob Aggregator"
	app.Usage = "The taiko blob aggreagtor software command line interface"
	app.Copyright = ""
	app.Description = "Blob aggregator implementation in Golang for Taiko protocol"
	app.Authors = []*cli.Author{{Name: "", Email: ""}}
	app.EnableBashCompletion = true

	// All supported sub commands.
	app.Commands = []*cli.Command{
		{
			Name:        "api",
			Flags:       flags.APIFlags,
			Usage:       "Starts the blob aggregator http API software",
			Description: "Taiko blob aggregator http API software",
			Action:      utils.SubcommandAction(new(api.API)),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
