package cep

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const viaCEPBase = "https://viacep.com.br/ws"

// ViaCEPItem is one result from the address search API.
type ViaCEPItem struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Erro        bool   `json:"erro"`
}

// SearchByAddress calls ViaCEP search by UF, city and street (logradouro).
// Returns up to 50 CEPs. UF and city must be non-empty; logradouro min 3 chars.
func SearchByAddress(uf, city, logradouro string) ([]ViaCEPItem, error) {
	uf = strings.TrimSpace(strings.ToUpper(uf))
	city = strings.TrimSpace(city)
	logradouro = strings.TrimSpace(logradouro)
	if len(uf) != 2 {
		return nil, fmt.Errorf("UF inválida: %s", uf)
	}
	if len(city) < 3 {
		return nil, fmt.Errorf("cidade deve ter ao menos 3 caracteres")
	}
	if len(logradouro) < 3 {
		logradouro = "Rua"
	}
	path := fmt.Sprintf("%s/%s/%s/%s/json/", viaCEPBase,
		url.PathEscape(uf),
		url.PathEscape(city),
		url.PathEscape(logradouro))
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UF inválida: %s", uf)
	}
	var items []ViaCEPItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	// Filter out erro and empty CEP
	valid := items[:0]
	for _, it := range items {
		if it.Erro || it.CEP == "" {
			continue
		}
		valid = append(valid, it)
	}
	return valid, nil
}

// NormalizeCEP returns 8 digits without hyphen (e.g. "01310100").
func NormalizeCEP(cep string) string {
	var b strings.Builder
	for _, r := range cep {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
