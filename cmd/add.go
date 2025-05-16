package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: T("add.description"),
	Long:  T("add.help"),
	Run: func(cmd *cobra.Command, args []string) {
		profileURL := cmd.Flag("url").Value.String()
		if profileURL == "" {
			fmt.Println(T("add.url_required"))
			cmd.Help()
			os.Exit(1)
		}

		// Extract profile name from URL if not specified
		profileName := cmd.Flag("profile").Value.String()
		if profileName == "" {
			profileName = extractProfileNameFromURL(profileURL)
		}

		// Verify if profile already exists
		profiles := viper.Get("profiles").([]interface{})
		for _, p := range profiles {
			profile := p.(map[string]interface{})
			if profile["name"] == profileName {
				fmt.Println(T("profile.exists", map[string]interface{}{"name": profileName}))
				os.Exit(1)
			}
		}

		// Clone the repository
		profilesPath := filepath.Join(os.Getenv("HOME"), ".persona", "profiles")
		profilePath := filepath.Join(profilesPath, profileName)

		if err := os.MkdirAll(profilePath, 0755); err != nil {
			fmt.Println(T("error.creating_directory", map[string]interface{}{"path": profilePath}))
			os.Exit(1)
		}

		_, err := git.PlainClone(profilePath, false, &git.CloneOptions{
			URL:      profileURL,
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Println(T("error.cloning_repository", map[string]interface{}{"url": profileURL}))
			os.Exit(1)
		}

		// Check existing files before creating symbolic links
		if !checkExistingFiles(profilePath) {
			fmt.Println(T("operation.canceled_by_user"))
			os.Exit(0)
		}

		// Create symbolic links
		files, err := os.ReadDir(profilePath)
		if err != nil {
			fmt.Println(T("error.reading_directory", map[string]interface{}{"path": profilePath}))
			os.Exit(1)
		}

		for _, file := range files {
			// Build source and destination paths
			source := filepath.Join(profilePath, file.Name())
			dest := filepath.Join(os.Getenv("HOME"), file.Name())

			// Remove existing file if it exists
			if _, err := os.Stat(dest); err == nil {
				if err := os.Remove(dest); err != nil {
					fmt.Println(T("error.removing_file", map[string]interface{}{"path": dest}))
					os.Exit(1)
				}
			}

			// Create symbolic link
			if err := os.Symlink(source, dest); err != nil {
				fmt.Println(T("error.creating_symlink", map[string]interface{}{"source": source, "dest": dest}))
				os.Exit(1)
			}
		}

		// Add profile to configuration with full path
		profiles = append(profiles, map[string]interface{}{
			"name": profileName,
			"url":  profileURL,
			"path": profilePath,
		})

		viper.Set("profiles", profiles)
		if err := viper.WriteConfig(); err != nil {
			fmt.Println(T("error.saving_config"))
			os.Exit(1)
		}

		fmt.Println(T("profile.added_successfully", map[string]interface{}{"name": profileName}))

		// Check if we should switch automatically
		if cmd.Flag("switch").Value.String() == "true" {
			viper.Set("current", profileName)
			if err := viper.WriteConfig(); err != nil {
				fmt.Println(T("error.updating_config"))
				os.Exit(1)
			}
			fmt.Println(T("switched_to_profile", map[string]interface{}{"name": profileName}))
		}
	},
}

// extractProfileNameFromURL extracts profile name from URL
func extractProfileNameFromURL(url string) string {
	// Remove .git if present
	url = strings.TrimSuffix(url, ".git")

	// Extract last segment of URL
	segments := strings.Split(url, "/")
	if len(segments) > 0 {
		return segments[len(segments)-1]
	}
	return "profile"
}

// checkExistingFiles checks for existing files and requests confirmation
func checkExistingFiles(profilePath string) bool {
	files, err := os.ReadDir(profilePath)
	if err != nil {
		fmt.Println(T("error.reading_directory", map[string]interface{}{"path": profilePath}))
		return false
	}

	for _, file := range files {
		homePath := filepath.Join(os.Getenv("HOME"), file.Name())
		if _, err := os.Stat(homePath); err == nil {
			fmt.Println(T("file.exists", map[string]interface{}{"path": homePath}))
			fmt.Println(T("replace.with_symlink"))
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				return false
			}
		}
	}

	return true
}

func init() {
	addCmd.Flags().StringP("profile", "p", "", T("add.flag.profile"))
	addCmd.Flags().StringP("url", "u", "", T("add.flag.url"))
	addCmd.MarkFlagRequired("url")
	rootCmd.AddCommand(addCmd)

	// Add automatic switch question
	addCmd.Flags().BoolP("switch", "s", false, T("add.flag.switch"))
}
