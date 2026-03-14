package cep

import (
	"fmt"
	"math/rand"
	"strings"
)

// FetchRandom returns a random real CEP.
// If uf is empty, picks any UF. If city is empty, picks any city in the UF.
// Returns 8 digits without mask (e.g. "01310100").
func FetchRandom(uf, city string) (string, error) {
	uf = strings.TrimSpace(strings.ToUpper(uf))
	city = strings.TrimSpace(city)

	if uf == "" {
		var err error
		uf, err = RandomUF()
		if err != nil {
			return "", err
		}
	} else {
		if len(uf) != 2 {
			return "", fmt.Errorf("UF inválida: %s", uf)
		}
		if _, ok := GetEstadoID(uf); !ok {
			return "", fmt.Errorf("UF inválida: %s", uf)
		}
	}

	if city == "" {
		var err error
		city, err = RandomMunicipio(uf)
		if err != nil {
			return "", err
		}
	}

	items, err := SearchByAddress(uf, city, "Rua")
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", fmt.Errorf("nenhum CEP encontrado para uf=%s cidade=%s", uf, city)
	}
	item := items[rand.Intn(len(items))]
	return NormalizeCEP(item.CEP), nil
}
