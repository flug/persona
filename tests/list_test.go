package tests

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/flug/persona/cmd"
)

func initI18N(t *testing.T) {
	viper.SetConfigType("json")
	viper.SetConfigName("en")
	viper.AddConfigPath("../i18n")

	// Vérifier que les fichiers de traduction existent
	for _, lang := range []string{"en.json", "fr.json", "de.json"} {
		if !assert.FileExists(t, filepath.Join("../i18n", lang)) {
			t.Fatal("Translation file missing:", lang)
		}
	}
}

func initConfig(t *testing.T) {
	initI18N(t)

	// Vérifier que les traductions sont correctement chargées
	translations := map[string]string{
		"list.title": "Available profiles",
		"list.current": "Current profile",
	}
	for key, expected := range translations {
		if value := viper.GetString(key); value != expected {
			t.Fatal("Translation mismatch for", key)
		}
	}
}

func GetRootCmd() *cobra.Command {
	return cmd.GetRootCmd()
}

func TestListCommand(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".persona", ".persona.json")

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	initConfig(t)

	viper.SetConfigFile(configPath)
	viper.SetDefault("profiles", []interface{}{
		map[string]interface{}{
			"name": "profile1",
			"url":  "https://github.com/user/profile1.git",
			"path": filepath.Join(tempDir, ".persona", "profiles", "profile1"),
		},
		map[string]interface{}{
			"name": "profile2",
			"url":  "https://github.com/user/profile2.git",
			"path": filepath.Join(tempDir, ".persona", "profiles", "profile2"),
		},
	})
	viper.SetDefault("current", "profile1")
	if err := viper.WriteConfig(); err != nil {
		t.Fatal(err)
	}

	rootCmd := cmd.GetRootCmd()
	listCmd := rootCmd.Commands()[0]
	listCmd.SetArgs([]string{})
	
	var buf bytes.Buffer
	oldWriter := listCmd.OutOrStdout()
	listCmd.SetOut(&buf)
	defer func() { listCmd.SetOut(oldWriter) }()

	err := listCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, viper.GetString("list.title"))
	assert.Contains(t, output, "* profile1")
	assert.Contains(t, output, "  profile2")
	assert.Contains(t, output, fmt.Sprintf(viper.GetString("list.current"), "profile1"))
}

func TestListTranslations(t *testing.T) {
	initConfig(t)

	// Vérifier les traductions
	assert.NotEmpty(t, viper.GetString("list.description"))
	assert.NotEmpty(t, viper.GetString("list.help"))
	assert.NotEmpty(t, viper.GetString("list.current_profile"))
	assert.NotEmpty(t, viper.GetString("list.no_profiles"))

	// Vérifier les traductions des messages de sortie
	assert.NotEmpty(t, viper.GetString("list.profile_active"))
	assert.NotEmpty(t, viper.GetString("list.profile_inactive"))
}

func TestListCommandOutput(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".persona", ".persona.json")

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	initConfig(t)

	viper.SetConfigFile(configPath)
	viper.SetDefault("profiles", []interface{}{
		map[string]interface{}{
			"name": "profile1",
			"url":  "https://github.com/user/profile1.git",
			"path": filepath.Join(tempDir, ".persona", "profiles", "profile1"),
		},
		map[string]interface{}{
			"name": "profile2",
			"url":  "https://github.com/user/profile2.git",
			"path": filepath.Join(tempDir, ".persona", "profiles", "profile2"),
		},
	})
	viper.SetDefault("current", "profile1")
	if err := viper.WriteConfig(); err != nil {
		t.Fatal(err)
	}

	rootCmd := GetRootCmd()
	listCmd := rootCmd.Commands()[0]
	listCmd.SetArgs([]string{})
	
	var buf bytes.Buffer
	oldWriter := listCmd.OutOrStdout()
	listCmd.SetOut(&buf)
	defer func() { listCmd.SetOut(oldWriter) }()

	err := listCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Available profiles:")
	assert.Contains(t, output, "* profile1")
	assert.Contains(t, output, "  profile2")
	assert.Contains(t, output, "\nCurrent profile: profile1")
}
