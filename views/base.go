// Package views implements different types of views that are used by the application.
package views

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

// Base represents the base view.
type Base struct {
	app *tview.Application
	*Client
}

// NewBase returns an instance of the base with new client.
func NewBase(app *tview.Application, config config.Config) *Base {
	return &Base{
		app:    app,
		Client: NewClient(config),
	}
}
