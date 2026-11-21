package main

import (
	"os"

	"github.com/v1kz25/openapi-contract-lab/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
