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

var version = "v0.1.0"

func init() {
	pflag.StringP("token", "t", "", "your github personal access token or oauth token")
	pflag.StringP("feeds-url", "f", "", "your github private feeds URL")
	pflag.BoolP("help", "h", false, "prints this message")
	pflag.BoolP("version", "v", false, "prints the version")
}

func visitFlags(flag *pflag.Flag) {
	if flag.Name == "help" {
		pflag.Usage()
		os.Exit(0)
	}
	if flag.Name == "version" {
		fmt.Println(version)
		os.Exit(0)
	}
	if flag.Name == "token" {
		_ = viper.BindPFlag(flag.Name, pflag.Lookup(flag.Name))
		if err := config.Write(); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("success!")
		os.Exit(0)
	}
	if flag.Name == "feeds-url" {
		_ = viper.BindPFlag("feeds_url", pflag.Lookup(flag.Name))
		if err := config.Write(); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("success!")
		os.Exit(0)
	}
}

func main() {
	pflag.Parse()
	pflag.Visit(visitFlags)

	cfg, err := config.GetConfig()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	gitApp := app.New(tview.NewApplication(), cfg)
	if err := gitApp.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
