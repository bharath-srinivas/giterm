package modules

import (
	"context"

	"github.com/gdamore/tcell"
	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"
	"golang.org/x/oauth2"

	"github.com/bharath-srinivas/giterm/config"
)

type Page struct {
	Name            string
	Parent          tview.Primitive
	ChildComponents []tview.Primitive
}

type GitApp struct {
	App         *tview.Application
	Context     context.Context
	Client      *github.Client
	GitAppPages []*Page

	pages  *tview.Pages
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

	gitApp := &GitApp{
		App:     app,
		Client:  client,
		Context: ctx,
		pages:   tview.NewPages(),
		config:  cfg,
	}
	gitApp.GitAppPages = append(gitApp.GitAppPages, gitApp.ProfilePage(), gitApp.RepoPage())
	for _, gitAppPage := range gitApp.GitAppPages {
		gitApp.pages.AddPage(gitAppPage.Name, gitAppPage.Parent, true, false)
	}
	gitApp.pages.SwitchToPage(gitApp.GitAppPages[0].Name)
	gitApp.App.SetRoot(gitApp.pages, true).SetFocus(gitApp.pages)
	return gitApp
}

func (g *GitApp) LoadInputHandler() {
	g.App.SetInputCapture(g.handleInput)
}

func (g *GitApp) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlQ:
		g.App.Stop()
		return event
	}
	return event
}
