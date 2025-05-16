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

func TestUpdateCommand(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".persona", ".persona.json")
	profilesDir := filepath.Join(tempDir, ".persona", "profiles")

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	initConfig(t)

	viper.SetConfigFile(configPath)
	viper.SetDefault("profiles", []interface{}{
		map[string]interface{}{
			"name": "profile1",
			"url":  "https://github.com/user/profile1.git",
			"path": filepath.Join(profilesDir, "profile1"),
		},
	})
	viper.SetDefault("current", "profile1")
	if err := viper.WriteConfig(); err != nil {
		t.Fatal(err)
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a profile",
		Run: func(cmd *cobra.Command, args []string) {
			profileName := args[0]
			profiles := viper.Get("profiles").([]interface{})

			found := false
			for _, p := range profiles {
				profile := p.(map[string]interface{})
				if profile["name"] == profileName {
					found = true
					profilePath := profile["path"].(string)
					if err := os.RemoveAll(profilePath); err != nil {
						panic(err)
					}
					if err := os.MkdirAll(profilePath, 0755); err != nil {
						panic(err)
					}

					fmt.Printf(viper.GetString("update.success"), profileName)
					break
				}
			}

			if !found {
				t.Fatal("Profile not found")
			}
		},
	}

	t.Run("UpdateExistingProfile", func(t *testing.T) {
		updateCmd.SetArgs([]string{"profile1"})
		if err := updateCmd.Execute(); err != nil {
			t.Fatal(err)
		}

		profiles := viper.Get("profiles").([]interface{})
		assert.Equal(t, 2, len(profiles))
	})
}
