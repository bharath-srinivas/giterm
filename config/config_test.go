package config

import (
	"reflect"
	"testing"
)

func Test_GetConfig(t *testing.T) {
	config := GetConfig()
	expectedType := "config.Config"
	returnType := reflect.TypeOf(config).String()
	if returnType != expectedType {
		t.Errorf("GetConfig returned incorrect type, got: %s, want: %s", returnType, expectedType)
	}
}
