package config

import (
	"os"
	"os/user"
	"path"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func Test_Write(t *testing.T) {
	configName = "test_config"
	viper.Set("test_token", "aTestOauthToken123")
	expected := ""
	if err := Write(); err != nil {
		t.Errorf("got: %s, want: %s", err.Error(), expected)
	}
}

func Test_GetConfig(t *testing.T) {
	config, _ := GetConfig()
	expectedType := "config.Config"
	returnType := reflect.TypeOf(config).String()
	if returnType != expectedType {
		cleanup()
		t.Errorf("GetConfig returned incorrect type, got: %s, want: %s", returnType, expectedType)
	}
	cleanup()
}

func cleanup() {
	if err := readConfig(); err != nil {
		return
	}

	usr, err := user.Current()
	if err != nil {
		return
	}

	configPath := path.Join(usr.HomeDir, configDir)
	if _, err := os.Stat(configPath); err == nil {
		if err = os.Chdir(configPath); err != nil {
			return
		}
		_ = os.Remove(configName + "." + configType)
	}
}
