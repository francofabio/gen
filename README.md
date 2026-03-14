# gen

Utilitário de linha de comando para geração rápida de dados de teste comuns no desenvolvimento de sistemas no Brasil. Um binário único, sem dependência de runtime no host.

## Recursos

- **CPF** — Gera CPF válido (com ou sem máscara)
- **CNPJ** — Gera CNPJ válido (com ou sem máscara)
- **CEP** — Retorna um CEP real (qualquer, por UF ou por cidade)
- **Cartão** — Gera número de cartão válido (Luhn) para bandeiras visa, master, amex, elo, hipercard

## Instalação

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/franco/gen/main/scripts/install.sh | sh
```

O script detecta OS e arquitetura, baixa o binário do [GitHub Releases](https://github.com/franco/gen/releases) e instala em `~/.local/bin` ou `~/bin`. Não usa `sudo`.

Se o diretório não estiver no `PATH`, adicione no seu shell:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/franco/gen/main/scripts/install.ps1 | iex
```

Instala em `%USERPROFILE%\bin`. Verifique se esse diretório está no `PATH`.

### Build a partir do código

Requisito: [Go 1.21+](https://go.dev/dl/).

```bash
git clone https://github.com/franco/gen.git
cd gen
go build -o gen ./cmd/gen
```

Para gerar binários para todas as plataformas:

```bash
./scripts/build.sh          # versão "dev"
./scripts/build.sh v1.0.0  # versão específica
```

Artefatos em `dist/`: `gen_darwin_arm64.tar.gz`, `gen_linux_amd64.tar.gz`, `gen_windows_amd64.zip`, etc.

## Uso

Forma geral:

```bash
gen <recurso> [argumentos] [flags]
```

A saída é apenas o valor gerado (uma linha), pronta para pipe.

### Exemplos

```bash
# CPF (sem máscara ou com -f/--format)
gen cpf
gen cpf -f

# CNPJ
gen cnpj
gen cnpj --format

# CEP (qualquer, por UF ou por cidade)
gen cep
gen cep es
gen cep es vitoria
gen cep sp campinas

# Cartão (bandeira e opcionalmente BIN)
gen card visa
gen card master
gen card visa 405168
gen card elo 636297
```

### Flags globais

| Flag | Descrição |
|------|-----------|
| `-h`, `--help` | Exibe a ajuda |
| `-v`, `--version` | Exibe a versão |
| `-c`, `--clipboard` | Copia o resultado também para o clipboard |

Exemplos com clipboard:

```bash
gen cpf -c
gen cnpj -c -f
gen cep es vitoria -c
gen card visa -c
```

No Linux, a flag `-c` usa `wl-copy`, `xclip` ou `xsel` (um deles precisa estar instalado). No macOS usa `pbcopy`; no Windows usa `clip`.

### Ajuda

```bash
gen help
gen --help
gen card --help
gen cep --help
```

## Configuração (opcional)

Arquivo: `~/.gen/config.json` (no Windows: `%USERPROFILE%\.gen\config.json`).

Use para definir BINs preferidos para o comando `card`. Se não informar BIN na linha de comando, o gen usa os valores da config; caso contrário, os padrões da bandeira.

Exemplo:

```json
{
  "cards": {
    "visa": ["405168", "411111"],
    "master": ["555555", "222100-222199"],
    "elo": ["636297"]
  }
}
```

BIN pode ser um valor fixo ou um intervalo (`"222100-222199"`). O arquivo é opcional; se não existir ou estiver inválido, a mensagem de erro será clara.

## Requisitos

- **CPF, CNPJ, cartão**: nenhum (totalmente offline).
- **CEP**: acesso à internet (ViaCEP e IBGE).
- **Clipboard (`-c`) no Linux**: `wl-copy` (Wayland), `xclip` ou `xsel` instalado.

## Plataformas

Binários para:

- macOS (arm64, amd64)
- Linux (amd64, arm64)
- Windows (amd64)

Build com Go puro (`CGO_ENABLED=0`), sem runtime externo.

## Licença

MIT. Ver [LICENSE](LICENSE).
