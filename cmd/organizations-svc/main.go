package main

import (
	"fmt"
	"os"

	"github.com/netbill/organizations-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}

func buildProcessID(service string) string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("%s-%s-%d", service, hostname, os.Getpid())
}
