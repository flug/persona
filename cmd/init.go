package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: T("init.description"),
	Long:  T("init.help"),
	Run: func(cmd *cobra.Command, args []string) {
		configDir := filepath.Join(os.Getenv("HOME"), ".persona")
		configPath := filepath.Join(configDir, ".persona.json")

		// Verify if configuration file exists
		if _, err := os.Stat(configPath); !os.IsNotExist(err) {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Println(T("init.exists_force", map[string]interface{}{"path": configPath}))
				return
			}
		}

		// Create configuration directory
		if err := os.MkdirAll(configDir, 0755); err != nil {
			panic(fmt.Errorf(T("init.error_dir", err)))
		}

		// Retrieve available translations
		translations, err := getAvailableTranslations()
		if err != nil {
			panic(err)
		}

		// Prompt user for preferred language
		fmt.Println(T("init.select_language"))
		for i, lang := range translations {
			fmt.Println(i+1, ".", lang)
		}

		var choice int
		fmt.Println(T("init.prompt_choice", len(translations)))
		fmt.Scan(&choice)

		if choice < 1 || choice > len(translations) {
			fmt.Println(T("init.invalid_choice"))
			return
		}

		// Create configuration file with selected language
		viper.SetConfigFile(configPath)
		viper.SetConfigType("json")
		viper.Set("profiles", []interface{}{})
		viper.Set("current", "")
		viper.Set("language", translations[choice-1])
		if err := viper.WriteConfig(); err != nil {
			panic(fmt.Errorf(T("init.error_file", err)))
		}

		fmt.Println(T("init.success", map[string]interface{}{"path": configPath}))
	},
}

func getAvailableTranslations() ([]string, error) {
	// Retrieve i18n directory path
	i18nDir := filepath.Join(os.Getenv("HOME"), ".persona", "i18n")
	if err := os.MkdirAll(i18nDir, 0755); err != nil {
		return nil, err
	}

	// List JSON files in directory
	files, err := os.ReadDir(i18nDir)
	if err != nil {
		return nil, err
	}

	// Extract filenames (removing .json extension)
	translations := make([]string, 0)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			translations = append(translations, strings.TrimSuffix(file.Name(), ".json"))
		}
	}

	// Add default translations if directory is empty
	if len(translations) == 0 {
		translations = append(translations, "en", "fr")
	}

	return translations, nil
}

var force bool

func init() {
	initCmd.Flags().BoolVarP(&force, "force", "f", false, T("init.force_flag"))
	viper.BindPFlag("force", initCmd.Flags().Lookup("force"))
	rootCmd.AddCommand(initCmd)
}
