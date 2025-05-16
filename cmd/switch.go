package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch profile",
	Long:  `Switch the active profile and apply the corresponding configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		profileName := cmd.Flag("profile").Value.String()
		if profileName == "" {
			fmt.Println(T("switch.profile_required"))
			cmd.Help()
			os.Exit(1)
		}

		// Verify if profile exists
		profiles := viper.Get("profiles").([]interface{})
		var profilePath string
		found := false

		for _, p := range profiles {
			profile := p.(map[string]interface{})
			if profile["name"] == profileName {
				found = true
				profilePath = profile["path"].(string)
				break
			}
		}

		if !found {
			fmt.Printf(T("switch.profile_not_found", map[string]interface{}{"name": profileName}))
			os.Exit(1)
		}

		// Update active profile in configuration
		viper.Set("current", profileName)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf(T("error.saving_config", map[string]interface{}{"error": err.Error()}))
			os.Exit(1)
		}

		// Traverse profile files and create symbolic links
		if err := applyProfileConfig(profilePath); err != nil {
			fmt.Printf(T("error.applying_profile", map[string]interface{}{"error": err.Error()}))
			os.Exit(1)
		}

		fmt.Printf(T("switch.success", map[string]interface{}{"name": profileName}))
	},
}

func init() {
	switchCmd.Flags().StringP("profile", "p", "", "nom du profil à activer")
	switchCmd.MarkFlagRequired("profile")
	rootCmd.AddCommand(switchCmd)
}

// applyProfileConfig applies the profile configuration by creating symbolic links
func applyProfileConfig(profilePath string) error {
	return filepath.Walk(profilePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore the profile root directory
		if path == profilePath {
			return nil
		}

		// Get the relative path with respect to the profile directory
		relPath, err := filepath.Rel(profilePath, path)
		if err != nil {
			return err
		}

		// Build the target path
		targetPath := filepath.Join(os.Getenv("HOME"), relPath)

		// Check if the target file already exists
		if _, err := os.Stat(targetPath); err == nil {
			// The file already exists
			if info.IsDir() {
				// For directories, check if it's a symbolic link
				if _, err := os.Readlink(targetPath); err == nil {
					// It's already a symbolic link, do nothing
					return nil
				}
			} else {
				// For files, check if it's a symbolic link
				if _, err := os.Readlink(targetPath); err == nil {
					// It's already a symbolic link, do nothing
					return nil
				}
			}

			// Ask for confirmation from the user
			fmt.Printf(T("switch.confirmation", targetPath))
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				return err
			}

			if strings.ToLower(response) != "o" {
				// The user does not want to replace
				return nil
			}

			// Remove the old file/directory
			if err := os.RemoveAll(targetPath); err != nil {
				return fmt.Errorf(T("error.removing_old_file", map[string]interface{}{"error": err.Error()}))
			}
		}

		// Create the parent directory if it does not exist
		parentDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf(T("error.creating_parent_directory", map[string]interface{}{"error": err.Error()}))
		}

		// Create the symbolic link
		if err := os.Symlink(path, targetPath); err != nil {
			return fmt.Errorf(T("error.creating_symbolic_link", map[string]interface{}{"error": err.Error()}))
		}

		fmt.Printf(T("switch.link_created", targetPath, path))
		return nil
	})
}
