package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/flug/persona/cmd"
)

func TestAddCommand(t *testing.T) {
	tempDir := t.TempDir()

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	initConfig(t)

	configDir := filepath.Join(tempDir, ".persona")
	configPath := filepath.Join(configDir, ".persona.json")

	viper.SetConfigFile(configPath)
	viper.SetDefault("profiles", []interface{}{})
	viper.SetDefault("current", "")
	viper.SetDefault("language", "en")

	if err := viper.WriteConfig(); err != nil {
		t.Fatal(err)
	}

	rootCmd := cmd.GetRootCmd()
	addCmd := rootCmd.Commands()[1] // Assuming add is the second command

	t.Run("AddNewProfile", func(t *testing.T) {
		addCmd.SetArgs([]string{"--profile", "test", "--url", "https://github.com/test/repo.git"})
		
		var buf bytes.Buffer
		oldWriter := addCmd.OutOrStdout()
		addCmd.SetOut(&buf)
		defer func() { addCmd.SetOut(oldWriter) }()

		err := addCmd.Execute()
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, viper.GetString("add.success"))

		profiles := viper.Get("profiles").([]interface{})
		assert.Len(t, profiles, 1)
		profile := profiles[0].(map[string]interface{})
		assert.Equal(t, "test", profile["name"])
		assert.Equal(t, "https://github.com/test/repo.git", profile["url"])

		// Vérifier que le dossier du profil a été créé
		profilePath := filepath.Join(configDir, "profiles", "test")
		assert.DirExists(t, profilePath)
	})

	t.Run("AddExistingProfile", func(t *testing.T) {
		addCmd.SetArgs([]string{"--profile", "test", "--url", "https://github.com/test/repo.git"})
		if err := addCmd.Execute(); err == nil {
			t.Fatal("Expected error for existing profile")
		}
	})

	t.Run("AddProfileWithInvalidURL", func(t *testing.T) {
		addCmd.SetArgs([]string{"--profile", "invalid", "--url", "invalid-url"})
		if err := addCmd.Execute(); err == nil {
			t.Fatal("Expected error for invalid URL")
		}
	})

	t.Run("AddProfileWithMissingURL", func(t *testing.T) {
		addCmd.SetArgs([]string{"--profile", "missing-url"})
		if err := addCmd.Execute(); err == nil {
			t.Fatal("Expected error for missing URL")
		}
	})
}
