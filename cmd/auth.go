package cmd

import (
	"fmt"

	"github.com/jrogala/jellyfin-cli/config"
	"github.com/jrogala/jellyfin-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("password", "p", "", "Password")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Jellyfin",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		result, err := ops.Authenticate(config.URL(), username, password)
		if err != nil {
			return err
		}

		if err := config.SaveSession(&config.Session{Token: result.Token, UserID: result.UserID}); err != nil {
			return err
		}
		fmt.Println("Authenticated successfully!")
		return nil
	},
}
