package app

import (
	"reflect"
	"testing"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

func Test_New(t *testing.T) {
	gitApp := New(tview.NewApplication(), config.Config{})
	expectedType := "*app.GitApp"
	returnType := reflect.TypeOf(gitApp).String()
	if returnType != expectedType {
		t.Errorf("Write returned incorrect type, got: %s, want: %s", returnType, expectedType)
	}
}
