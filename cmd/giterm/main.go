package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rivo/tview"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bharath-srinivas/giterm/app"
	"github.com/bharath-srinivas/giterm/config"
)

func init() {
	pflag.String("token", "", "your github personal access token")
}

// sets the personal access token
func setToken() {
	token := viper.GetString("token")
	if token == "" {
		log.Fatal(fmt.Errorf("token cannot be empty"))
	} else if err := config.New("token", token); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("token set successfully.")
}

func main() {
	pflag.Parse()
	pflag.Visit(func(flag *pflag.Flag) {
		if flag.Name == "token" {
			_ = viper.BindPFlags(pflag.CommandLine)
			setToken()
			os.Exit(0)
		}
	})
	gitApp := app.New(tview.NewApplication())
	gitApp.Run()
}
