// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/rai-project/auth"
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
