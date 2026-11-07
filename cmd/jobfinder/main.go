package main

import (
	"os"

	"github.com/akkul143/ai-analysis/internal/cli"
)

func main() {
	os.Exit(cli.Execute(os.Args[1:], os.Stdout, os.Stderr))
}
