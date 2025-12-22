package main

import (
	"os"

	"github.com/umisto/cities-svc/cmd/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
