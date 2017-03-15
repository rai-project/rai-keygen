package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rai-project/auth"
	"github.com/rai-project/auth/provider"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	username  string
	email     string
	firstname string
	lastname  string
	isColor   bool
	isVerbose bool
	isDebug   bool
)

var RootCmd = &cobra.Command{
	Use:   "rai-keygen",
	Short: "Generates profiles to be used with the rai client",
	Long: "Generates a profile file that needs to be placed in `~/.rai_profile` (linux/OSX) or " +
		"`%HOME%/.rai_profile` (Windows -- for me this is `C:\\Users\\abduld\\.rai_profile`). " +
		"The rai client reads these configuration files to authenticate the user. " +
		"A seed (specified by `secret`) is used to generate secure credentials",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		prof, err := provider.New(
			auth.Username(username),
			auth.Email(email),
			auth.Firstname(firstname),
			auth.Lastname(lastname))
		if err != nil {
			return err
		}

		if err := prof.Create(); err != nil {
			return err
		}

		fmt.Print(prof.Info().String())

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

	RootCmd.AddCommand(cmd.VersionCmd)
	RootCmd.AddCommand(cmd.LicenseCmd)
	RootCmd.AddCommand(cmd.EnvCmd)
	RootCmd.AddCommand(cmd.GendocCmd)
	RootCmd.AddCommand(cmd.CompletionCmd)

	RootCmd.Flags().StringVarP(&username, "username", "u", "", "The username to generate the key for.")
	RootCmd.Flags().StringVarP(&email, "email", "e", "", "The email to generate the key for.")
	RootCmd.Flags().StringVarP(&firstname, "firstname", "f", "", "The firstname to generate the key for.")
	RootCmd.Flags().StringVarP(&lastname, "lastname", "l", "", "The lastname to generate the key for.")
	RootCmd.PersistentFlags().StringVarP(&config.DefaultAppSecret, "secret", "s", "", "The application secret key.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", !color.NoColor, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")

	viper.BindPFlag("app.secret", RootCmd.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("app.color", RootCmd.PersistentFlags().Lookup("color"))

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.Init(
		config.AppName("rai"),
		config.ConfigString(configContent),
	)
}
