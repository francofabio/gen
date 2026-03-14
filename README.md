# gen

Command-line utility for quickly generating common Brazilian test data. Single binary, no runtime dependency on the host.

## Features

- **CPF** — Generate valid CPF (with or without mask)
- **CNPJ** — Generate valid CNPJ (with or without mask)
- **CEP** — Return a real postal code (any, by state, or by city)
- **Card** — Generate Luhn-valid card numbers for visa, master, amex, elo, hipercard

## Installation

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/francofabio/gen/main/scripts/install.sh | sh
```

The script detects OS and architecture, downloads the binary from [GitHub Releases](https://github.com/francofabio/gen/releases), and installs to `~/.local/bin` or `~/bin`. No `sudo` required.

If that directory is not in your `PATH`, add it to your shell:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/francofabio/gen/main/scripts/install.ps1 | iex
```

Installs to `%USERPROFILE%\bin`. Ensure that directory is in your `PATH`.

### Build from source

Requires [Go 1.21+](https://go.dev/dl/).

```bash
git clone https://github.com/francofabio/gen.git
cd gen
go build -o gen ./cmd/gen
```

To build binaries for all platforms:

```bash
./scripts/build.sh          # version "dev"
./scripts/build.sh v1.0.0  # specific version
```

Artifacts go to `dist/`: `gen_darwin_arm64.tar.gz`, `gen_linux_amd64.tar.gz`, `gen_windows_amd64.zip`, etc.

## Usage

General form:

```bash
gen <resource> [arguments] [flags]
```

Output is a single line (the generated value), pipe-friendly.

### Examples

```bash
# CPF (plain or with -f/--format)
gen cpf
gen cpf -f

# CNPJ
gen cnpj
gen cnpj --format

# CEP (any, by state/UF, or by city)
gen cep
gen cep es
gen cep es vitoria
gen cep sp campinas

# Card (brand and optional BIN)
gen card visa
gen card master
gen card visa 405168
gen card elo 636297
```

### Global flags

| Flag | Description |
|------|-------------|
| `-h`, `--help` | Show help |
| `-v`, `--version` | Show version |
| `-c`, `--clipboard` | Also copy the result to the clipboard |

Examples with clipboard:

```bash
gen cpf -c
gen cnpj -c -f
gen cep es vitoria -c
gen card visa -c
```

On Linux, the `-c` flag uses `wl-copy`, `xclip`, or `xsel` (one of them must be installed). On macOS it uses `pbcopy`; on Windows, `clip`.

### Help

```bash
gen help
gen --help
gen card --help
gen cep --help
```

### Language

Help and error messages follow the `LANG` (or `LANGUAGE`) environment variable. Default is English.

```bash
gen help                    # English (default)
LANG=pt_BR gen help         # Portuguese (Brazil)
LANG=pt_BR.UTF-8 gen cpf    # Portuguese
```

Supported locales: `en`, `pt-BR` (or `pt_BR`).

## Configuration (optional)

Config file: `~/.gen/config.json` (on Windows: `%USERPROFILE%\.gen\config.json`).

Use it to define preferred BINs for the `card` command. If you don’t pass a BIN on the command line, gen uses values from config; otherwise it uses the brand defaults.

Example:

```json
{
  "cards": {
    "visa": ["405168", "411111"],
    "master": ["555555", "222100-222199"],
    "elo": ["636297"]
  }
}
```

BIN can be a fixed value or a range (`"222100-222199"`). The file is optional; if it’s missing or invalid, you’ll get a clear error message.

## Requirements

- **CPF, CNPJ, card**: none (fully offline).
- **CEP**: internet access (ViaCEP and IBGE).
- **Clipboard (`-c`) on Linux**: `wl-copy` (Wayland), `xclip`, or `xsel` installed.

## Platforms

Binaries for:

- macOS (arm64, amd64)
- Linux (amd64, arm64)
- Windows (amd64)

Built with pure Go (`CGO_ENABLED=0`), no external runtime.

## License

MIT. See [LICENSE](LICENSE).
