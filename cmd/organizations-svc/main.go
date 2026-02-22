package main

import (
	"os"

	"github.com/netbill/organizations-svc/internal/cli"
)

func main() {
	cli.Run(os.Args)
}
