package card

import (
	"math/rand"
	"strconv"
	"strings"
)

// ValidBrands are the supported card brands in v1.
var ValidBrands = map[string]bool{
	"visa": true, "master": true, "amex": true, "elo": true, "hipercard": true,
}

// BrandInfo holds length and default BIN(s) for a brand.
type BrandInfo struct {
	Length int
	BINs   []string // can be "405168" or "222100-222199"
}

// DefaultBINs returns default BIN/prefix per brand when no config or arg is provided.
var DefaultBINs = map[string]BrandInfo{
	"visa":     {Length: 16, BINs: []string{"4"}},
	"master":   {Length: 16, BINs: []string{"51", "52", "53", "54", "55", "222100-272099"}},
	"amex":     {Length: 15, BINs: []string{"34", "37"}},
	"elo":      {Length: 16, BINs: []string{"636297", "438935", "504175", "451416", "636368"}},
	"hipercard": {Length: 16, BINs: []string{"606282"}},
}

// ResolveBIN returns a 6+ digit BIN to use: from explicitBin if non-empty, else from configBINs, else from DefaultBINs.
// configBINs can be from config.Cards[brand] (e.g. ["405168", "411111"] or ["222100-222199"]).
func ResolveBIN(brand string, explicitBin string, configBINs []string) string {
	if explicitBin != "" {
		return normalizeBIN(explicitBin)
	}
	if len(configBINs) > 0 {
		s := configBINs[rand.Intn(len(configBINs))]
		return pickOneBIN(s)
	}
	info, ok := DefaultBINs[brand]
	if !ok {
		return ""
	}
	s := info.BINs[rand.Intn(len(info.BINs))]
	return pickOneBIN(s)
}

// normalizeBIN strips non-digits and ensures at least 6 digits for PAN generation.
func normalizeBIN(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// pickOneBIN interprets s as either a single BIN (e.g. "405168") or a range (e.g. "222100-222199")
// and returns one BIN (full digits) to use.
func pickOneBIN(s string) string {
	idx := strings.Index(s, "-")
	if idx >= 0 {
		low := normalizeBIN(s[:idx])
		high := normalizeBIN(s[idx+1:])
		if low != "" && high != "" {
			lowN, _ := strconv.Atoi(low)
			highN, _ := strconv.Atoi(high)
			if lowN > highN {
				lowN, highN = highN, lowN
			}
			n := lowN
			if highN > lowN {
				n = lowN + rand.Intn(highN-lowN+1)
			}
			return strconv.Itoa(n)
		}
	}
	return normalizeBIN(s)
}

// LengthForBrand returns the PAN length for the brand.
func LengthForBrand(brand string) int {
	if info, ok := DefaultBINs[brand]; ok {
		return info.Length
	}
	return 16
}
