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

func TestRemoveCommand(t *testing.T) {
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

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a profile",
		Run: func(cmd *cobra.Command, args []string) {
			profileName := args[0]
			profiles := viper.Get("profiles").([]interface{})

			var index int
			found := false
			for i, p := range profiles {
				if p.(map[string]interface{})["name"] == profileName {
					index = i
					found = true
					break
				}
			}

			if !found {
				t.Fatal("Profile not found")
			}

			profilePath := profiles[index].(map[string]interface{})["path"].(string)
			if err := os.RemoveAll(profilePath); err != nil {
				panic(err)
			}

			profiles = append(profiles[:index], profiles[index+1:]...)
			viper.Set("profiles", profiles)

			if viper.Get("current").(string) == profileName {
				if len(profiles) > 0 {
					viper.Set("current", profiles[0].(map[string]interface{})["name"])
				} else {
					viper.Set("current", "")
				}
			}

			if err := viper.WriteConfig(); err != nil {
				panic(err)
			}

			fmt.Printf("Profile removed successfully: %s\n", profileName)
		},
	}

	t.Run("RemoveExistingProfile", func(t *testing.T) {
		removeCmd.SetArgs([]string{"profile1"})
		if err := removeCmd.Execute(); err != nil {
			t.Fatal(err)
		}

		profiles := viper.Get("profiles").([]interface{})
		assert.Len(t, profiles, 1)
		assert.Equal(t, "profile2", profiles[0].(map[string]interface{})["name"])
		assert.Equal(t, "profile2", viper.Get("current"))

		// Vérifier que le dossier du profil a été supprimé
		profilePath := filepath.Join(profilesDir, "profile1")
		assert.NoDirExists(t, profilePath)
	})

	t.Run("RemoveNonExistingProfile", func(t *testing.T) {
		removeCmd.SetArgs([]string{"nonexistent"})
		if err := removeCmd.Execute(); err == nil {
			t.Fatal("Expected error for non-existing profile")
		}
	})

	t.Run("RemoveLastProfile", func(t *testing.T) {
		removeCmd.SetArgs([]string{"profile2"})
		if err := removeCmd.Execute(); err != nil {
			t.Fatal(err)
		}

		profiles := viper.Get("profiles").([]interface{})
		assert.Len(t, profiles, 0)
		assert.Equal(t, "", viper.Get("current"))

		// Vérifier que le dossier du profil a été supprimé
		profilePath := filepath.Join(profilesDir, "profile2")
		assert.NoDirExists(t, profilePath)
	})
}
