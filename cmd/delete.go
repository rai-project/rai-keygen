package cmd

import "github.com/spf13/cobra"

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a user",
	Long:  "The config file must be configured with the propper mailing credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {

		log := log.WithField("cmd", "delete")

		log.Fatal("Unimplemented")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
