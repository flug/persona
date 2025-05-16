package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a profile",
	Long:  `Removes a profile and its symbolic links. After removal, the application may prompt to switch profiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get profile name from flag
		profileName := cmd.Flag("profile").Value.String()
		if profileName == "" {
			fmt.Println("The --profile parameter is required")
			cmd.Help()
			os.Exit(1)
		}

		// Check if profile exists
		profiles := viper.Get("profiles").([]interface{})
		var profilePath string
		var profileIndex int
		found := false

		for i, p := range profiles {
			profile := p.(map[string]interface{})
			if profile["name"] == profileName {
				found = true
				profilePath = profile["path"].(string)
				profileIndex = i
				break
			}
		}

		if !found {
			fmt.Printf(T("profile.not_found", map[string]interface{}{"name": profileName}))
			os.Exit(1)
		}

		// Ask for confirmation
		fmt.Printf(T("remove.confirmation", map[string]interface{}{"name": profileName}))
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			fmt.Println(T("error.reading_response"))
			os.Exit(1)
		}

		if strings.ToLower(response) != "y" {
			fmt.Println(T("operation.cancelled"))
			os.Exit(0)
		}

		// Remove symbolic links
		if err := removeProfileLinks(profilePath); err != nil {
			fmt.Printf(T("error.removing_links", map[string]interface{}{"error": err.Error()}))
			os.Exit(1)
		}

		// Remove profile directory
		if err := os.RemoveAll(profilePath); err != nil {
			fmt.Printf(T("error.removing_directory", map[string]interface{}{"error": err.Error()}))
			os.Exit(1)
		}

		// Remove profile from configuration
		profiles = append(profiles[:profileIndex], profiles[profileIndex+1:]...)
		viper.Set("profiles", profiles)

		// Update current profile if needed
		if viper.GetString("current") == profileName {
			viper.Set("current", "")
		}

		// Save configuration
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf(T("error.saving_configuration", map[string]interface{}{"error": err.Error()}))
			os.Exit(1)
		}

		fmt.Printf(T("profile.removed_successfully", map[string]interface{}{"name": profileName}))

		// Prompt to switch profiles if others remain
		if len(profiles) > 0 {
			fmt.Println(T("remove.switch_prompt"))
			for i, p := range profiles {
				profile := p.(map[string]interface{})
				fmt.Printf("%d) %s\n", i+1, profile["name"])
			}
			fmt.Println(T("remove.switch_prompt_2"))
			fmt.Println(T("remove.switch_prompt_3"))

			var choice int
			fmt.Print(T("remove.choice"))
			if _, err := fmt.Scanln(&choice); err != nil {
				fmt.Println(T("error.reading_choice"))
				os.Exit(1)
			}

			if choice > 0 && choice <= len(profiles) {
				selectedProfile := profiles[choice-1].(map[string]interface{})["name"].(string)
				fmt.Printf(T("remove.switching_to_profile", map[string]interface{}{"name": selectedProfile}))
				if err := executeSwitch(selectedProfile); err != nil {
					fmt.Printf(T("error.switching_profile", map[string]interface{}{"error": err.Error()}))
					os.Exit(1)
				}
			}
		}
	},
}

func init() {
	removeCmd.Flags().StringP("profile", "p", "", "name of the profile to remove")
	removeCmd.MarkFlagRequired("profile")
	rootCmd.AddCommand(removeCmd)
}

// removeProfileLinks removes all symbolic links pointing to the profile
func removeProfileLinks(profilePath string) error {
	return filepath.Walk(profilePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore the root profile directory
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

		// Check if the symbolic link exists
		if _, err := os.Readlink(targetPath); err == nil {
			// Remove the symbolic link
			if err := os.Remove(targetPath); err != nil {
				return fmt.Errorf(T("error.removing_link", map[string]interface{}{"error": err.Error()}))
			}
			fmt.Printf(T("removed.symbolic_link", map[string]interface{}{"path": targetPath}))
		}

		return nil
	})
}

// executeSwitch executes the switch command with the specified profile
func executeSwitch(profileName string) error {
	cmd := &cobra.Command{}
	cmd.SetArgs([]string{"switch", "--profile", profileName})
	switchCmd.Run(cmd, nil)
	return nil
}
