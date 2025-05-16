package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Profile struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Path    string   `json:"path"`
	Aliases []string `json:"aliases"`
}

type Config struct {
	Profiles []Profile `json:"profiles"`
	Current  string    `json:"current"`
}

func init() {
	// Initialize configuration
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("~/.persona2")

	// Create the configuration directory if it does not exist
	configDir := filepath.Join(os.Getenv("HOME"), ".persona2")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf(T("error.config_dir", map[string]interface{}{"error": err.Error()}))
		os.Exit(1)
	}

	// Load configuration
	if err := viper.ReadInConfig(); err != nil {
		// If file doesn't exist, create default configuration
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.Set("profiles", []Profile{})
			viper.Set("current", "")
			if err := viper.WriteConfig(); err != nil {
				fmt.Printf(T("error.config_create", map[string]interface{}{"error": err.Error()}))
				os.Exit(1)
			}
		} else {
			fmt.Printf(T("error.config_read", map[string]interface{}{"error": err.Error()}))
			os.Exit(1)
		}
	}
}
