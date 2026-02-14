package main

import (
	"os"

	"github.com/netbill/organizations-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
