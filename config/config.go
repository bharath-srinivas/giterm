// Package config implements the configurations required by the application to function properly.
package config

import (
	"errors"
	"os"
	"os/user"
	"path"

	"github.com/spf13/viper"
)

// Config represents the configuration fields of the giterm app.
type Config struct {
	Token    string
	FeedsUrl string
}

// Write creates a new config file if one does not exist or modifies the existing config file. It returns an error if it
//fails to do both.
func Write() error {
	_ = readConfig()
	if err := viper.WriteConfig(); err != nil {
		return viper.SafeWriteConfig()
	}
	return nil
}

// GetConfig returns an instance of config.
func GetConfig() (Config, error) {
	var config Config
	if err := readConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return config, errors.New("config file not found. please set your personal access token")
		}
		return config, err
	}
	config.Token = viper.GetString("token")
	config.FeedsUrl = viper.GetString("feeds_url")
	return config, nil
}

// setConfigPath sets the config path if it exists already or creates a new one if it doesn't exist and set the newly created config path.
func setConfigPath() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	configDir := ".giterm"
	configPath := path.Join(usr.HomeDir, configDir)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.Mkdir(configPath, 0777); err != nil {
			return err
		}
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	return nil
}

// readConfig reads the config file from the set config path and returns an error if it fails to find the config file.
func readConfig() error {
	if err := setConfigPath(); err != nil {
		return err
	}
	return viper.ReadInConfig()
}
