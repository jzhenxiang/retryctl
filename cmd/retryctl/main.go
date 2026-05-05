package main

import (
	"os"

	"github.com/retryctl/retryctl/cmd/retryctl/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
