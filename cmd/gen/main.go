package main

import (
	"os"

	"github.com/francofabio/gen/internal/cli"
	"github.com/francofabio/gen/internal/output"
)

var Version = "dev"

func main() {
	cli.Version = Version
	code := cli.Run(os.Args[1:])
	output.Exit(code)
}
