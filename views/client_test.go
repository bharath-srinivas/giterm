package views

import (
	"reflect"
	"testing"

	"github.com/bharath-srinivas/giterm/config"
)

func TestNewClient(t *testing.T) {
	client := NewClient(config.Config{})
	expected := "*views.Client"
	actual := reflect.TypeOf(client).String()
	if actual != expected {
		t.Errorf("NewClient returned an incorrect type, got: %s, want: %s", actual, expected)
	}
}
