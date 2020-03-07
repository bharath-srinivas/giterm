package modules

import (
	"context"

	"github.com/gdamore/tcell"
	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"
	"golang.org/x/oauth2"

	"github.com/bharath-srinivas/giterm/config"
)

type GitApp struct {
	App     *tview.Application
	Context context.Context
	Client  *github.Client
	Widgets map[string]tview.Primitive

	config config.Config
}

func New(app *tview.Application) *GitApp {
	config.Init()
	cfg := config.GetConfig()

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	httpClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(httpClient)

	return &GitApp{
		Context: ctx,
		Client:  client,
		App:     app,
		Widgets: map[string]tview.Primitive{},
		config:  cfg,
	}
}

func (g *GitApp) LoadWidgets() {
	g.Widgets["profile"] = g.ProfileWidget()
	g.Widgets["repositories"] = g.RepoWidget()
}

func (g *GitApp) LoadInputHandler() {
	g.App.SetInputCapture(g.handleInput)
}

func (g *GitApp) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlQ:
		g.App.Stop()
		return event
	case tcell.KeyTab:
		if g.Widgets["profile"].GetFocusable().HasFocus() {
			g.App.SetFocus(g.Widgets["repositories"])
		} else {
			g.App.SetFocus(g.Widgets["profile"])
		}
		return event
	}
	return event
}
