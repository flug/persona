package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "persona",
	Short: T("root.short_description"),
	Long:  T("root.long_description"),
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	// Load the configuration only if the command is not "init"
	cmd := rootCmd.CalledAs()
	if cmd != "init" && cmd != "" {
		initConfig()
	}
	return nil
}

// GetRootCmd returns the root command
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// Initialize translations
	initI18N()

	c := rootCmd
	c.PersistentFlags().StringVar(&cfgFile, "config", "", T("common.config_file"))
	viper.BindPFlag("config", c.PersistentFlags().Lookup("config"))

	// Add commands
	c.AddCommand(
		addCmd,
		listCmd,
		switchCmd,
		removeCmd,
		updateCmd,
		selfUpdateCmd,
		initCmd,
	)
}

// initConfig initializes the configuration after the command execution
func initConfig() {
	// Default configuration
	viper.SetDefault("config", filepath.Join(os.Getenv("HOME"), ".persona.json"))
	viper.SetDefault("profiles", []interface{}{})
	viper.SetDefault("current", "")

	// Load the configuration
	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf(T("error.config_load", map[string]interface{}{"error": err.Error()}))
			os.Exit(1)
		}
	}

	// Create the configuration directory if it does not exist
	configDir := filepath.Dir(viper.GetString("config"))
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf(T("error.config_dir", map[string]interface{}{"error": err.Error()}))
		os.Exit(1)
	}
}
