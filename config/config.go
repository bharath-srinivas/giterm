// Package config implements the configurations required by the application to function properly.
package config

import (
	"log"
	"os"
	"os/user"
	"path"

	"github.com/spf13/viper"
)

// Config represents the configuration fields of the giterm app.
type Config struct {
	Token string
}

var config Config

func init() {
	if err := readConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("config file not found. please set your personal access token.")
		} else {
			log.Println(err.Error())
		}
		os.Exit(1)
	}
	config.Token = viper.GetString("token")
}

// New creates a new config file if one does not exist or modifies the existing config file. It returns an error if it
//fails to do both.
func New(key, value string) error {
	setConfigPath()
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return viper.SafeWriteConfig()
	}
	return nil
}

// GetConfig returns the initialized config.
func GetConfig() Config {
	return config
}

// setConfigPath sets the config path if it exists already or creates a new one if it doesn't exist and set the newly created config path.
func setConfigPath() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configDir := ".giterm"
	configPath := path.Join(usr.HomeDir, configDir)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.Mkdir(configPath, 0777); err != nil {
			log.Fatal(err.Error())
		}
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
}

// readConfig reads the config file from the set config path and returns an error if it fails to find the config file.
func readConfig() error {
	setConfigPath()
	return viper.ReadInConfig()
}
