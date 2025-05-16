package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: T("list.description"),
	Long: T("list.help"),
	Run: func(cmd *cobra.Command, args []string) {
		// Get profiles
		profiles := viper.Get("profiles").([]interface{})
		current := viper.Get("current").(string)

		// Display title
		fmt.Println(T("list.title"))

		// Display each profile
		for _, p := range profiles {
			profile := p.(map[string]interface{})
			name := profile["name"].(string)
			url := profile["url"].(string)
			aliases := profile["aliases"].([]interface{})

			// Mark current profile
			if name == current {
				name = "* " + name
			}

			// Display profile information
			fmt.Printf(T("list.profile_info"), name, url)
			if len(aliases) > 0 {
				fmt.Printf(T("list.aliases"), aliases)
			}
		}

		// Display current profile
		if current != "" {
			fmt.Printf(T("list.current"), current)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
