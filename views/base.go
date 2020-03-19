package views

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

type Base struct {
	app *tview.Application
	*Client
}

func NewBase(app *tview.Application, config config.Config) *Base {
	return &Base{
		app:    app,
		Client: NewClient(config),
	}
}
