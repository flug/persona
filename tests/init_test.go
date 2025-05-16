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

func TestInitCommand(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".persona", ".persona.json")

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	initConfig(t)

	rootCmd := cmd.GetRootCmd()
	initCmd := rootCmd.Commands()[0] // Assuming init is the first command

	t.Run("InitWithoutForce", func(t *testing.T) {
		var buf bytes.Buffer
		oldWriter := initCmd.OutOrStdout()
		initCmd.SetOut(&buf)
		defer func() { initCmd.SetOut(oldWriter) }()

		initCmd.SetArgs([]string{})
		err := initCmd.Execute()
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, viper.GetString("init.success"))

		assert.FileExists(t, configPath)
		assert.FileExists(t, filepath.Join(tempDir, ".persona", "i18n", "en.json"))
		assert.FileExists(t, filepath.Join(tempDir, ".persona", "profiles"))
	})

	t.Run("InitWithForce", func(t *testing.T) {
		// Créer un fichier de configuration existant
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(configPath, []byte("existing config"), 0644); err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer
		oldWriter := initCmd.OutOrStdout()
		initCmd.SetOut(&buf)
		defer func() { initCmd.SetOut(oldWriter) }()

		initCmd.SetArgs([]string{"--force"})
		err := initCmd.Execute()
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, viper.GetString("init.success"))

		assert.FileExists(t, configPath)
		config, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotContains(t, string(config), "existing config")
	})

	t.Run("NormalInitialization", func(t *testing.T) {
		initCmd.SetArgs([]string{})
		err := initCmd.Execute()
		assert.NoError(t, err)

		assert.FileExists(t, configPath)

		profiles := viper.Get("profiles").([]interface{})
		assert.Len(t, profiles, 0)
		assert.Equal(t, "", viper.Get("current"))

		assert.Equal(t, "en", viper.Get("language"))

		// Vérifier que les dossiers ont été créés
		configDir := filepath.Dir(configPath)
		assert.DirExists(t, filepath.Join(configDir, "i18n"))
		assert.DirExists(t, filepath.Join(configDir, "profiles"))

		// Vérifier que les fichiers de traduction ont été copiés
		i18nDir := filepath.Join(configDir, "i18n")
		for _, lang := range []string{"en.json", "fr.json", "de.json"} {
			assert.FileExists(t, filepath.Join(i18nDir, lang))
		}

		initCmd.SetArgs([]string{"--force"})
		err = initCmd.Execute()
		assert.NoError(t, err)

		profiles = viper.Get("profiles").([]interface{})
		assert.Len(t, profiles, 0)
	})

	t.Run("ForceInitialization", func(t *testing.T) {
		if err := os.WriteFile(configPath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		if err := initCmd.Execute(); err != nil {
			t.Fatal(err)
		}

		assert.FileExists(t, configPath)

		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []interface{}{}, viper.Get("profiles"))
		assert.Equal(t, "", viper.Get("current"))

		// Vérifier que les dossiers ont été créés
		configDir := filepath.Dir(configPath)
		assert.DirExists(t, filepath.Join(configDir, "i18n"))
		assert.DirExists(t, filepath.Join(configDir, "profiles"))

		// Vérifier que les fichiers de traduction ont été copiés
		i18nDir := filepath.Join(configDir, "i18n")
		for _, lang := range []string{"en.json", "fr.json", "de.json"} {
			assert.FileExists(t, filepath.Join(i18nDir, lang))
		}
	})
}
