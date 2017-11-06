package cmd

import "github.com/spf13/cobra"

var getKeyCmd = &cobra.Command{
	Use:   "getkey",
	Short: "Retrieves the key for a user.",
	Long:  "The config file must be configured with the propper mailing credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {

		log := log.WithField("cmd", "getkey")

		log.Fatal("Unimplemented")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(getKeyCmd)
}
