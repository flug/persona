package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: T("update.description"),
	Long:  T("update.long_description"),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := cmd.Flag("profile").Value.String()
		profiles := viper.Get("profiles").([]interface{})

		// If a specific profile is requested, only process that profile
		if profileName != "" {
			if err := updateProfile(profileName, profiles); err != nil {
				fmt.Printf(T("error.updating_profile", map[string]interface{}{"name": profileName, "error": err.Error()}))
				os.Exit(1)
			}
			return
		}

		// Otherwise, update all profiles
		for _, p := range profiles {
			profile := p.(map[string]interface{})
			name := profile["name"].(string)
			if err := updateProfile(name, profiles); err != nil {
				fmt.Printf(T("error.updating_profile", map[string]interface{}{"name": name, "error": err.Error()}))
				os.Exit(1)
			}
		}
	},
}

func init() {
	updateCmd.Flags().StringP("profile", "p", "", "name of the profile to update (optional)")
	rootCmd.AddCommand(updateCmd)
}

// updateProfile updates a specific profile
func updateProfile(profileName string, profiles []interface{}) error {
	// Find the profile
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
		return fmt.Errorf(T("error.profile_not_found", map[string]interface{}{"name": profileName}))
	}

	// Ouvrir le repository
	repo, err := git.PlainOpen(profilePath)
	if err != nil {
		return fmt.Errorf(T("error.opening_repository", map[string]interface{}{"error": err.Error()}))
	}

	// Récupérer l'URL du remote origin
	remote, err := repo.Remote("origin")
	if err != nil {
		return fmt.Errorf(T("error.fetching_remote", map[string]interface{}{"error": err.Error()}))
	}

	// Fetch les derniers changements
	fmt.Printf(T("update.fetching", profileName))
	if err := remote.Fetch(&git.FetchOptions{
		Progress: os.Stdout,
	}); err != nil {
		return fmt.Errorf(T("error.fetching", map[string]interface{}{"error": err.Error()}))
	}

	// Obtenir le HEAD
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf(T("error.fetching_worktree", map[string]interface{}{"error": err.Error()}))
	}

	// Pull les changements
	if err := worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	}); err != nil {
		// Si la branche est déjà à jour, c'est OK
		if err.Error() == "already up-to-date" {
			fmt.Printf(T("update.already_up_to_date", profileName))
			return nil
		}
		return fmt.Errorf(T("error.pull", map[string]interface{}{"error": err.Error()}))
	}

	fmt.Printf(T("update.success", profileName))
	return nil
}
