package card

import (
	"math/rand"
	"strings"

	"github.com/francofabio/gen/internal/i18n"
)

// LuhnCheckDigit returns the check digit (0-9) for the given digits string (only digits).
func LuhnCheckDigit(digits string) byte {
	var sum int
	parity := len(digits) % 2
	for i, r := range digits {
		if r < '0' || r > '9' {
			continue
		}
		n := int(r - '0')
		if i%2 == parity {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
	}
	check := (10 - (sum % 10)) % 10
	return byte('0' + check)
}

// ValidLuhn reports whether the digit string (PAN) passes Luhn.
func ValidLuhn(digits string) bool {
	var sum int
	parity := len(digits) % 2
	for i, r := range digits {
		if r < '0' || r > '9' {
			return false
		}
		n := int(r - '0')
		if i%2 == parity {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
	}
	return sum%10 == 0
}

// Generate produces a valid PAN for the given brand and optional BIN.
// binPrefix is from ResolveBIN (already resolved from config/defaults if empty).
// configBINs is config.Cards[brand] for when no explicit BIN was passed.
func Generate(brand string, explicitBIN string, configBINs []string) (string, error) {
	brand = strings.ToLower(brand)
	if !ValidBrands[brand] {
		return "", ErrInvalidBrand{brand}
	}
	length := LengthForBrand(brand)
	bin := ResolveBIN(brand, explicitBIN, configBINs)
	if bin == "" {
		bin = DefaultBINs[brand].BINs[0]
		if len(bin) < 6 {
			// e.g. "4" for Visa: pad to at least 6 for BIN
			for len(bin) < 6 {
				bin += string(byte('0' + rand.Intn(10)))
			}
		}
	}
	if len(bin) > length-1 {
		bin = bin[:length-1]
	}
	// Fill to length-1 with random digits, then Luhn
	for len(bin) < length-1 {
		bin += string(byte('0' + rand.Intn(10)))
	}
	pan := bin + string(LuhnCheckDigit(bin))
	return pan, nil
}

// ErrInvalidBrand is returned when the brand is not supported.
type ErrInvalidBrand struct {
	Brand string
}

func (e ErrInvalidBrand) Error() string {
	return i18n.T("card_invalid_brand", e.Brand)
}
