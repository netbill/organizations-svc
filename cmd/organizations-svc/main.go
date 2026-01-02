package main

import (
	"os"

	"github.com/netbill/organizations-svc/cmd/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
