package config

import (
	"log"
	"os"
	"os/user"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	Token string
}

var config Config

func New(key, value string) error {
	setConfigPath()
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return viper.SafeWriteConfig()
	}
	return nil
}

func Init() {
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

func GetConfig() Config {
	return config
}

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

func readConfig() error {
	setConfigPath()
	return viper.ReadInConfig()
}
