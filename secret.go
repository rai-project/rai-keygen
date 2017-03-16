package main

import "github.com/rai-project/rai-keygen/cmd"

var AppSecret string

func init() {
	cmd.AppSecret = AppSecret
}
