package views

import (
	"reflect"
	"testing"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

func TestNewTextView(t *testing.T) {
	textView := NewTextView(tview.NewApplication(), config.Config{}, true)
	expected := "*views.TextWidget"
	actual := reflect.TypeOf(textView).String()
	if actual != expected {
		t.Errorf("NewTextView returned an incorrect type, got: %s, want: %s", actual, expected)
	}
}
