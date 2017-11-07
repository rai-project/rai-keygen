package cmd

import (
	"github.com/rai-project/auth"
	"github.com/rai-project/auth/provider"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a user",
	Long:  "The config file must be configured with the propper mailing credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {

		prof, err := provider.New(
			auth.Username(""),
			auth.Email(email),
		)
		if err != nil {
			return err
		}

		err = prof.Delete()
		return err
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&email, "email", "e", "", "The email for the account to delete.")
	RootCmd.AddCommand(deleteCmd)
}
