package cmd

import "github.com/spf13/cobra"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Retrieves the key for a user.",
	Long:  "The config file must be configured with the propper mailing credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {

		log := log.WithField("cmd", "update")

		log.Fatal("Unimplemented")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
