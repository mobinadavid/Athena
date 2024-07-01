package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"sync"
)

//go:embed locales/*.json
var localeFS embed.FS

var (
	Locales  = []string{"fa", "en", "fr", "es", "ar", "tr"}
	once     sync.Once
	instance *i18n.Bundle
)

// Init initializes the global i18n bundle instance.
// It's safe to be called multiple times but only the first call will have an effect.
func Init() error {
	once.Do(func() {
		instance = i18n.NewBundle(language.English)
		instance.RegisterUnmarshalFunc("json", json.Unmarshal)

		for _, locale := range Locales {
			loadMessageFile(fmt.Sprintf("locales/%s.json", locale))
		}
	})

	return nil
}

// loadMessageFile is a helper function to load a message file into the bundle.
func loadMessageFile(filename string) {
	content, err := localeFS.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read locale file '%s': %v", filename, err)
	}

	_, err = instance.ParseMessageFileBytes(content, filename)
	if err != nil {
		log.Fatalf("Could not parse locale file '%s': %v", filename, err)
	}
}

// Localize uses the singleton instance to localize a message ID for a given language.
// It accepts the language code, message ID, and optionally any necessary template data.
// The templateData parameter is variadic, allowing calls without template data.
func Localize(lang, messageID string, templateData ...map[string]interface{}) string {
	localizer := i18n.NewLocalizer(instance, lang)

	var config *i18n.LocalizeConfig
	if len(templateData) > 0 && templateData[0] != nil {
		// If templateData is provided, use it
		config = &i18n.LocalizeConfig{
			MessageID:    messageID,
			TemplateData: templateData[0],
		}
	} else {
		// If no templateData is provided, just use the message ID
		config = &i18n.LocalizeConfig{
			MessageID: messageID,
		}
	}

	message, err := localizer.Localize(config)
	if err != nil {
		log.Printf("Could not localize message ID '%s' for language '%s': %v", messageID, lang, err)
		return messageID // Fallback to message ID in case of an error
	}

	return message
}
