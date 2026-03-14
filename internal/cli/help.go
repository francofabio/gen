package cli

const usageShort = "uso: gen <recurso> [argumentos] [flags]"

const helpFull = `gen — geração rápida de dados de teste

Uso:
  gen <recurso> [argumentos] [flags]

Recursos:
  cpf         Gera um CPF válido
  cnpj        Gera um CNPJ válido
  cep         Retorna um CEP real (opcional: UF, cidade)
  card        Gera número de cartão válido (bandeira [BIN])
  version     Exibe a versão
  help        Exibe esta ajuda

Exemplos:
  gen cpf
  gen cpf -f
  gen cnpj --format
  gen cep
  gen cep es
  gen cep es vitoria
  gen card visa
  gen card master
  gen card visa 405168

Flags globais:
  -h, --help        ajuda
  -v, --version     versão
  -c, --clipboard   copia o resultado também para o clipboard
`

const helpCEP = `gen cep — retorna um CEP real

Uso:
  gen cep [uf] [cidade]

Exemplos:
  gen cep
  gen cep es
  gen cep es vitoria
  gen cep sp campinas
`

const helpCard = `gen card — gera número de cartão válido para testes

Uso:
  gen card <bandeira> [bin]

Bandeiras: visa, master, amex, elo, hipercard

Exemplos:
  gen card visa
  gen card master
  gen card visa 405168
  gen card elo 636297
`
