package cnpj

import (
	"fmt"
	"math/rand"
	"strings"
)

// Generate returns a valid 14-digit CNPJ (digits only).
// Avoids trivial sequences like 00000000000000.
func Generate() string {
	const size = 12
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = byte('0' + rand.Intn(10))
	}
	// Ensure not all zeros
	allZero := true
	for _, c := range buf {
		if c != '0' {
			allZero = false
			break
		}
	}
	if allZero {
		buf[0] = '1'
	}
	d1 := digitCNPJ(buf, false)
	d2 := digitCNPJ(append(buf, d1), true)
	return string(buf) + string([]byte{d1, d2})
}

// digitCNPJ computes one CNPJ check digit.
// First DV: 12 digits × weights [5,4,3,2,9,8,7,6,5,4,3,2]
// Second DV: 13 digits × weights [6,5,4,3,2,9,8,7,6,5,4,3,2]
func digitCNPJ(prefix []byte, second bool) byte {
	var weights []int
	if second {
		weights = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	} else {
		weights = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	}
	n := len(prefix)
	if n > len(weights) {
		n = len(weights)
	}
	var sum int
	for i := 0; i < n; i++ {
		c := prefix[i]
		if c < '0' || c > '9' {
			continue
		}
		sum += int(c-'0') * weights[i]
	}
	rem := sum % 11
	if rem < 2 {
		return '0'
	}
	return byte('0' + 11 - rem)
}

// Format returns the CNPJ with mask 12.345.678/0001-95.
// Input should be 14 digits; non-digits are stripped.
func Format(s string) string {
	var digits []rune
	for _, r := range s {
		if r >= '0' && r <= '9' {
			digits = append(digits, r)
		}
	}
	if len(digits) != 14 {
		return s
	}
	return fmt.Sprintf("%s.%s.%s/%s-%s",
		string(digits[0:2]),
		string(digits[2:5]),
		string(digits[5:8]),
		string(digits[8:12]),
		string(digits[12:14]))
}

// Strip removes all non-digit characters from s.
func Strip(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
