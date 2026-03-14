package cli

import (
	"io"
	"os"
	"strings"

	"github.com/francofabio/gen/internal/card"
	"github.com/francofabio/gen/internal/cep"
	"github.com/francofabio/gen/internal/clipboard"
	"github.com/francofabio/gen/internal/config"
	"github.com/francofabio/gen/internal/cnpj"
	"github.com/francofabio/gen/internal/cpf"
	"github.com/francofabio/gen/internal/output"
)

// Run parses args and runs the appropriate subcommand.
// Returns exit code (0 for success).
func Run(args []string) int {
	args, copyToClipboard := stripClipboardFlag(args)
	if len(args) == 0 {
		printUsage(os.Stderr)
		return 1
	}
	cmd := strings.ToLower(args[0])
	rest := args[1:]

	switch cmd {
	case "cpf":
		return runCPF(rest, copyToClipboard)
	case "cnpj":
		return runCNPJ(rest, copyToClipboard)
	case "cep":
		return runCEP(rest, copyToClipboard)
	case "card":
		return runCard(rest, copyToClipboard)
	case "version":
		return runVersion(rest)
	case "help":
		return runHelp(rest)
	}
	if isHelpFlag(cmd) {
		output.Err(os.Stderr, helpFull)
		return 0
	}
	if isVersionFlag(cmd) {
		return runVersion(nil)
	}
	printUsage(os.Stderr)
	return 1
}

// stripClipboardFlag removes -c and --clipboard from args and returns (newArgs, hadFlag).
func stripClipboardFlag(args []string) ([]string, bool) {
	var out []string
	var found bool
	for _, a := range args {
		if a == "-c" || a == "--clipboard" {
			found = true
			continue
		}
		out = append(out, a)
	}
	return out, found
}

func runCPF(args []string, copyToClipboard bool) int {
	format := parseFormatFlag(args)
	v := cpf.Generate()
	if format {
		v = cpf.Format(v)
	}
	output.PrintValue(os.Stdout, v)
	if copyToClipboard {
		tryCopyToClipboard(v)
	}
	return 0
}

func runCNPJ(args []string, copyToClipboard bool) int {
	format := parseFormatFlag(args)
	v := cnpj.Generate()
	if format {
		v = cnpj.Format(v)
	}
	output.PrintValue(os.Stdout, v)
	if copyToClipboard {
		tryCopyToClipboard(v)
	}
	return 0
}

// parseFormatFlag returns true if -f or --format appears in args.
func parseFormatFlag(args []string) bool {
	for _, a := range args {
		if a == "-f" || a == "--format" {
			return true
		}
	}
	return false
}

func printUsage(w io.Writer) {
	output.Err(w, usageShort)
}

// OutWriter returns stdout for result output.
func OutWriter() io.Writer {
	return os.Stdout
}

func isHelpFlag(s string) bool {
	return s == "-h" || s == "--help"
}

func isVersionFlag(s string) bool {
	return s == "-v" || s == "--version"
}

// Version is set by main (e.g. via ldflags -X).
var Version = "dev"

func runCEP(args []string, copyToClipboard bool) int {
	if len(args) >= 1 && (args[0] == "--help" || args[0] == "-h") {
		output.Err(os.Stdout, helpCEP)
		return 0
	}
	var uf, city string
	if len(args) >= 1 {
		uf = strings.TrimSpace(args[0])
	}
	if len(args) >= 2 {
		city = strings.Join(args[1:], " ")
	}
	cepStr, err := cep.FetchRandom(uf, city)
	if err != nil {
		output.Err(os.Stderr, err.Error())
		return 1
	}
	output.PrintValue(os.Stdout, cepStr)
	if copyToClipboard {
		tryCopyToClipboard(cepStr)
	}
	return 0
}

func runCard(args []string, copyToClipboard bool) int {
	if len(args) >= 1 && (args[0] == "--help" || args[0] == "-h") {
		output.Err(os.Stdout, helpCard)
		return 0
	}
	if len(args) < 1 {
		output.Err(os.Stderr, "uso: gen card <bandeira> [bin]")
		return 1
	}
	brand := strings.ToLower(args[0])
	if !card.ValidBrands[brand] {
		output.Err(os.Stderr, "bandeira inválida: "+brand)
		return 1
	}
	var explicitBIN string
	if len(args) >= 2 {
		explicitBIN = args[1]
	}
	cfg, err := config.Load()
	if err != nil {
		output.Err(os.Stderr, "config: "+err.Error())
		return 1
	}
	configBINs := cfg.Cards[brand]
	pan, err := card.Generate(brand, explicitBIN, configBINs)
	if err != nil {
		output.Err(os.Stderr, err.Error())
		return 1
	}
	output.PrintValue(os.Stdout, pan)
	if copyToClipboard {
		tryCopyToClipboard(pan)
	}
	return 0
}

func tryCopyToClipboard(text string) {
	if err := clipboard.Write(text); err != nil {
		output.Err(os.Stderr, "aviso: não foi possível copiar para o clipboard")
	}
}

func runVersion(args []string) int {
	output.PrintValue(os.Stdout, Version)
	return 0
}

func runHelp(args []string) int {
	// Help goes to stdout so it can be piped (e.g. gen help | less).
	output.PrintValue(os.Stdout, strings.TrimSuffix(helpFull, "\n"))
	return 0
}
