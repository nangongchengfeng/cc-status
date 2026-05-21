package main

import (
	"fmt"
	"os"

	"cc-status/client/internal/cli"
)

func main() {
	app := cli.NewApp(cli.NewBootstrapRunner(cli.BootstrapOptions{
		EnvLookup: os.Getenv,
	}))

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
