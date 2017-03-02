package cmd

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/rai-project/auth"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/spf13/cobra"
)

var (
	username  string
	appsecret string
)

var RootCmd = &cobra.Command{
	Use: "rai-keygen",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := auth.Hash(username)
		if err != nil {
			return err
		}
		enc := toml.NewEncoder(os.Stdout)
		err = enc.Encode(User{
			Username:  username,
			AccessKey: accessKey,
			SecretKey: secretKey,
		})
		if err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.AddCommand(cmd.VersionCmd)
	RootCmd.AddCommand(cmd.LicenseCmd)
	RootCmd.AddCommand(cmd.EnvCmd)

	RootCmd.Flags().StringVarP(&username, "username", "u", "",
		"The username to generate the command for.")
	RootCmd.Flags().StringVarP(&appsecret, "appsecret", "s", "",
		"The application secret key.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.Init(
		config.AppName("rai"),
	)
}
