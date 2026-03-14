package cep

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const ibgeEstados = "https://servicodados.ibge.gov.br/api/v1/localidades/estados"

// IBGEEstado is one state from the IBGE API.
type IBGEEstado struct {
	ID    int    `json:"id"`
	Sigla string `json:"sigla"`
	Nome  string `json:"nome"`
}

// IBGEMunicipio is one municipality.
type IBGEMunicipio struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

// ufToID maps UF sigla (e.g. "ES") to IBGE state id. Filled on first use.
var ufToID map[string]int

func ensureEstados() error {
	if len(ufToID) > 0 {
		return nil
	}
	list, err := fetchEstados()
	if err != nil {
		return err
	}
	ufToID = make(map[string]int)
	for _, e := range list {
		ufToID[strings.ToUpper(e.Sigla)] = e.ID
	}
	return nil
}

func fetchEstados() ([]IBGEEstado, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(ibgeEstados)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha ao listar estados: %d", resp.StatusCode)
	}
	var list []IBGEEstado
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}

// GetEstadoID returns IBGE state id for the given UF (e.g. "ES" -> 32).
// UF is case-insensitive.
func GetEstadoID(uf string) (int, bool) {
	if err := ensureEstados(); err != nil {
		return 0, false
	}
	id, ok := ufToID[strings.ToUpper(strings.TrimSpace(uf))]
	return id, ok
}

// ListEstados returns all UF siglas (e.g. ["SP", "RJ", ...]).
func ListEstados() ([]string, error) {
	if err := ensureEstados(); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(ufToID))
	for sigla := range ufToID {
		out = append(out, sigla)
	}
	return out, nil
}

// FetchMunicipios returns municipality names for the given UF.
func FetchMunicipios(uf string) ([]string, error) {
	id, ok := GetEstadoID(uf)
	if !ok {
		return nil, fmt.Errorf("UF inválida: %s", uf)
	}
	url := fmt.Sprintf("https://servicodados.ibge.gov.br/api/v1/localidades/estados/%d/municipios", id)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UF inválida: %s", uf)
	}
	var list []IBGEMunicipio
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	names := make([]string, len(list))
	for i, m := range list {
		names[i] = m.Nome
	}
	return names, nil
}

// RandomUF returns a random UF sigla.
func RandomUF() (string, error) {
	list, err := ListEstados()
	if err != nil {
		return "", err
	}
	return list[rand.Intn(len(list))], nil
}

// RandomMunicipio returns a random municipality name for the given UF.
func RandomMunicipio(uf string) (string, error) {
	names, err := FetchMunicipios(uf)
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("nenhum município para UF: %s", uf)
	}
	return names[rand.Intn(len(names))], nil
}
