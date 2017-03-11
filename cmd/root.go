package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/rai-project/auth"
	"github.com/rai-project/auth/auth0"
	"github.com/rai-project/auth/secret"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	username  string
	email     string
	appsecret string
	isColor   bool
	isVerbose bool
	isDebug   bool
)

var RootCmd = &cobra.Command{
	Use:          "rai-keygen",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var prof auth.Profile

		profOpts := []auth.ProfileOption{
			auth.Username(username),
			auth.Email(email),
		}
		provider := auth.Provider(strings.ToLower(auth.Config.Provider))
		switch provider {
		case auth.Auth0Provider:
			prof, err = auth0.NewProfile(profOpts...)
		case auth.SecretProvider:
			prof, err = secret.NewProfile(profOpts...)
		default:
			err = errors.Errorf("the auth provider %v specified is not supported", provider)
		}
		if err != nil {
			return err
		}

		if err := prof.Create(); err != nil {
			return err
		}

		user := prof.Info()

		enc := toml.NewEncoder(os.Stdout)
		err = enc.Encode(user)
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
		"The username to generate the key for.")
	RootCmd.Flags().StringVarP(&email, "email", "e", "",
		"The email to generate the key for.")
	RootCmd.Flags().StringVarP(&appsecret, "appsecret", "s", "",
		"The application secret key.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", color.NoColor, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")

	viper.BindPFlag("app.secret", RootCmd.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("app.color", RootCmd.PersistentFlags().Lookup("color"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.Init(
		config.AppName("rai"),
	)
}
