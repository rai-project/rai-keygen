package cmd

import (
	"fmt"

	"github.com/rai-project/auth"
	"github.com/rai-project/auth/provider"
	"github.com/spf13/cobra"
)

var getKeyCmd = &cobra.Command{
	Use:   "getkey",
	Short: "Retrieves the key for a user.",
	Long:  "The config file must be configured with the propper mailing credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// log := log.WithField("cmd", "getkey")

		prof, err := provider.New(
			auth.Email(email),
		)
		if err != nil {
			return err
		}

		if err := prof.GetByEmail(); err != nil {
			return err
		}

		fmt.Print(prof.Info().String())

		return nil
	},
}

func init() {
	getKeyCmd.Flags().StringVarP(&email, "email", "e", "", "The email to fetch the key for.")
	RootCmd.AddCommand(getKeyCmd)
}
