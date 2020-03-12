package app

import (
	"log"
	"os"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

type GitApp struct {
	app      *tview.Application
	appPages []*page
	config   config.Config
	pages    *tview.Pages
}

func New(app *tview.Application) *GitApp {
	config.Init()
	cfg := config.GetConfig()

	gitApp := &GitApp{
		app:    app,
		pages:  tview.NewPages(),
		config: cfg,
	}
	gitApp.LoadPages()
	gitApp.app.SetInputCapture(gitApp.inputHandler)
	gitApp.app.SetRoot(gitApp.pages, true).SetFocus(gitApp.pages)
	return gitApp
}

func (g *GitApp) Start() {
	if err := g.app.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
