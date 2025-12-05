package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Manager struct {
	bundle  *i18n.Bundle
	locales []string
}

var hiddenLocales = []string{}

func NewManager() *Manager {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	return &Manager{
		bundle:  bundle,
		locales: []string{},
	}
}

func (m *Manager) LoadDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			path := filepath.Join(dir, entry.Name())
			if _, err := m.bundle.LoadMessageFile(path); err != nil {
				return fmt.Errorf("failed to load %s: %w", entry.Name(), err)
			}

			langTag := strings.TrimSuffix(entry.Name(), ".json")
			exists := slices.Contains(m.locales, langTag)
			if !exists {
				m.locales = append(m.locales, langTag)
			}
		}
	}
	return nil
}

func (m *Manager) GetLocales() []string {
	result := make([]string, 0, len(m.locales))
	for _, lang := range m.locales {
		if !slices.Contains(hiddenLocales, lang) {
			result = append(result, lang)
		}
	}

	slices.SortFunc(result, func(a, b string) int {
		if a == "es" {
			return -1
		}
		if b == "es" {
			return 1
		}
		if a == "en" {
			return -1
		}
		if b == "en" {
			return 1
		}
		return strings.Compare(a, b)
	})

	return result
}

func (m *Manager) Localizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(m.bundle, lang)
}

func (m *Manager) Translate(lang, key string, templateData ...map[string]any) string {
	localizer := m.Localizer(lang)
	config := &i18n.LocalizeConfig{
		MessageID: key,
	}

	if len(templateData) > 0 {
		config.TemplateData = templateData[0]
	}

	msg, err := localizer.Localize(config)
	if err != nil {
		return key
	}
	return msg
}
