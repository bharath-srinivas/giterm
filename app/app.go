// Package app implements the terminal based giterm application.
package app

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/pages"
)

// GitApp represents the main application.
type GitApp struct {
	app      *tview.Application
	appPages pages.Pages
	pages    *tview.Pages
}

// New returns a new instance of GitApp.
func New(app *tview.Application, config config.Config) *GitApp {
	gitApp := &GitApp{
		app:   app,
		pages: tview.NewPages(),
	}
	gitApp.appPages = pages.MakePages(app, config)
	for _, page := range gitApp.appPages {
		gitApp.pages.AddPage(page.Name, page.Widgets.Parent, true, false)
	}
	gitApp.pages.SwitchToPage(gitApp.appPages[0].Name)
	gitApp.app.SetInputCapture(gitApp.inputHandler)
	gitApp.app.SetRoot(gitApp.pages, true).SetFocus(gitApp.pages)
	return gitApp
}

// Run starts the GitApp application in an event loop.
func (g *GitApp) Run() error {
	return g.app.Run()
}
