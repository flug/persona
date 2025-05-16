package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSwitchCommand(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".persona", ".persona.json")
	profilesDir := filepath.Join(tempDir, ".persona", "profiles")

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	viper.SetConfigFile(configPath)
	viper.SetDefault("profiles", []interface{}{
		map[string]interface{}{
			"name":    "profile1",
			"url":     "https://github.com/user/profile1.git",
			"path":    filepath.Join(profilesDir, "profile1"),
			"aliases": []interface{}{"profile1", "p1"},
		},
		map[string]interface{}{
			"name":    "profile2",
			"url":     "https://github.com/user/profile2.git",
			"path":    filepath.Join(profilesDir, "profile2"),
			"aliases": []interface{}{"profile2"},
		},
	})
	viper.SetDefault("current", "profile1")
	if err := viper.WriteConfig(); err != nil {
		t.Fatal(err)
	}

	switchCmd := &cobra.Command{
		Use:   "switch",
		Short: "Switch to a profile",
		Run: func(cmd *cobra.Command, args []string) {
			profileName := args[0]
			profiles := viper.Get("profiles").([]interface{})

			found := false
			for _, p := range profiles {
				profile := p.(map[string]interface{})
				if profile["name"] == profileName {
					found = true
					viper.Set("current", profileName)
					if err := viper.WriteConfig(); err != nil {
						panic(err)
					}
					fmt.Printf("Switched to profile: %s\n", profileName)
					break
				}
			}

			if !found {
				t.Fatal("Profile not found")
			}
		},
	}

	t.Run("SwitchExistingProfile", func(t *testing.T) {
		switchCmd.SetArgs([]string{"profile2"})
		if err := switchCmd.Execute(); err != nil {
			t.Fatal(err)
		}

		current := viper.Get("current").(string)
		assert.Equal(t, "profile2", current)

		// Vérifier que le changement a été sauvegardé
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "profile2", viper.Get("current"))
	})

	t.Run("SwitchNonExistingProfile", func(t *testing.T) {
		switchCmd.SetArgs([]string{"nonexistent"})
		if err := switchCmd.Execute(); err == nil {
			t.Fatal("Expected error for non-existing profile")
		}
	})

	t.Run("SwitchUsingAlias", func(t *testing.T) {
		switchCmd.SetArgs([]string{"p1"})
		if err := switchCmd.Execute(); err != nil {
			t.Fatal(err)
		}

		current := viper.Get("current").(string)
		assert.Equal(t, "profile1", current)
	})
}
