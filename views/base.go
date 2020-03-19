package views

import (
	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

type Base struct {
	app *tview.Application

	*Client
	*github.ListOptions
	*github.ListCursorOptions
	*github.Response
}

func NewBase(app *tview.Application, config config.Config) *Base {
	return &Base{
		app:               app,
		Client:            NewClient(config),
		ListOptions:       &github.ListOptions{},
		ListCursorOptions: &github.ListCursorOptions{},
		Response:          &github.Response{},
	}
}
