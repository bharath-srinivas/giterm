package config

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func Test_Write(t *testing.T) {
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
		t.Errorf("GetConfig returned incorrect type, got: %s, want: %s", returnType, expectedType)
	}
}
