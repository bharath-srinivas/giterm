package app

import (
	"reflect"
	"testing"

	"github.com/rivo/tview"
)

func Test_New(t *testing.T) {
	gitApp := New(tview.NewApplication())
	expectedType := "*app.GitApp"
	returnType := reflect.TypeOf(gitApp).String()
	if returnType != expectedType {
		t.Errorf("New returned incorrect type, got: %s, want: %s", returnType, expectedType)
	}
}
