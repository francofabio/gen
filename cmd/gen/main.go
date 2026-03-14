package main

import (
	"os"

	"github.com/franco/gen/internal/cli"
	"github.com/franco/gen/internal/output"
)

var Version = "dev"

func main() {
	cli.Version = Version
	code := cli.Run(os.Args[1:])
	output.Exit(code)
}
