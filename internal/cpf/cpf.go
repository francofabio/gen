package cpf

import (
	"fmt"
	"math/rand"
	"strings"
)

// Generate returns a valid 11-digit CPF (digits only).
func Generate() string {
	const size = 9
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = byte('0' + rand.Intn(10))
	}
	d1 := digitCPF(buf, 10)
	d2 := digitCPF(append(buf, d1), 11)
	return string(buf) + string([]byte{d1, d2})
}

// digitCPF computes one CPF check digit (weights 2..10 or 2..11).
func digitCPF(prefix []byte, maxWeight int) byte {
	var sum int
	for i, c := range prefix {
		if c < '0' || c > '9' {
			continue
		}
		w := maxWeight - i
		sum += int(c-'0') * w
	}
	rem := sum % 11
	if rem < 2 {
		return '0'
	}
	return byte('0' + 11 - rem)
}

// Format returns the CPF with mask 123.456.789-09.
// Input should be 11 digits (no mask); non-digits are stripped.
func Format(s string) string {
	var digits []rune
	for _, r := range s {
		if r >= '0' && r <= '9' {
			digits = append(digits, r)
		}
	}
	if len(digits) != 11 {
		return s
	}
	return fmt.Sprintf("%s.%s.%s-%s",
		string(digits[0:3]),
		string(digits[3:6]),
		string(digits[6:9]),
		string(digits[9:11]))
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
