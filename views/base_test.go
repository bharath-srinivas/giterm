package views

import (
	"reflect"
	"testing"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

func TestNewBase(t *testing.T) {
	base := NewBase(tview.NewApplication(), config.Config{})
	expected := "*views.Base"
	actual := reflect.TypeOf(base).String()
	if actual != expected {
		t.Errorf("NewBase returned an incorrect type, got: %s, want: %s", actual, expected)
	}
}
