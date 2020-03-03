package config

import "os"

type Config struct {
	Token string
}

var config Config

func Init() {
	config.Token = os.Getenv("TOKEN")
}

func GetConfig() Config {
	return config
}
