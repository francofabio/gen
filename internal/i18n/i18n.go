package i18n

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

//go:embed locales/en.json
var enJSON []byte

//go:embed locales/pt-BR.json
var ptBRJSON []byte

var (
	mu       sync.RWMutex
	locale   string
	msgMap   map[string]string
	fallback map[string]string
	ptBRMap  map[string]string
)

func init() {
	fallback = mustParse(enJSON)
	ptBRMap = mustParse(ptBRJSON)
	msgMap = fallback
	locale = "en"
}

func mustParse(b []byte) map[string]string {
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		panic("i18n: invalid locale json: " + err.Error())
	}
	return m
}

// Init sets the current locale from a language tag (e.g. "en", "pt-BR", "pt_BR.UTF-8").
// If lang is empty, uses the LANG or LANGUAGE environment variable.
// Normalizes pt_BR -> pt-BR and uses the first segment before any dot.
// Falls back to "en" if the requested locale is not available.
func Init(lang string) {
	mu.Lock()
	defer mu.Unlock()
	if lang == "" {
		lang = os.Getenv("LANG")
		if lang == "" {
			lang = os.Getenv("LANGUAGE")
		}
	}
	// LANGUAGE can be "pt_BR:en" — use first part
	if i := strings.Index(lang, ":"); i > 0 {
		lang = strings.TrimSpace(lang[:i])
	}
	lang = normalizeLang(lang)
	switch lang {
	case "pt-BR", "pt_BR":
		msgMap = ptBRMap
		locale = "pt-BR"
	default:
		msgMap = fallback
		locale = "en"
	}
}

func normalizeLang(lang string) string {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return "en"
	}
	if i := strings.Index(lang, "."); i > 0 {
		lang = lang[:i]
	}
	return strings.ReplaceAll(lang, "_", "-")
}

// T returns the translation for key, with args applied via fmt.Sprintf.
// If the key is missing, falls back to "en" then returns the key.
func T(key string, args ...interface{}) string {
	mu.RLock()
	s, ok := msgMap[key]
	mu.RUnlock()
	if !ok {
		mu.RLock()
		s, ok = fallback[key]
		mu.RUnlock()
	}
	if !ok {
		if len(args) > 0 {
			return fmt.Sprintf(key+" %v", args...)
		}
		return key
	}
	if len(args) > 0 {
		return fmt.Sprintf(s, args...)
	}
	return s
}
