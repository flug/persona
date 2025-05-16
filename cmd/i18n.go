package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

// Initialize translations
func initI18N() {
	// Get language from configuration
	lang := viper.GetString("language")
	if lang == "" {
		// Default to English if no language is set
		lang = "en"
	}

	// Configure Viper for translations
	viper.SetConfigName(lang)
	viper.AddConfigPath("i18n")
	viper.SetConfigType("json")

	// Load translations
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(T("error.reading_translations", map[string]interface{}{"error": err.Error()}))
	}
}

// T translates a translation key
func T(key string, args ...interface{}) string {
	// Get translation
	translation := viper.GetString(key)
	if translation == "" {
		// If translation doesn't exist, use the key as text
		return key
	}

	// Format the string with arguments
	return fmt.Sprintf(translation, args...)
}
