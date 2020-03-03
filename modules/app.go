package modules

import (
	"context"

	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"
	"golang.org/x/oauth2"

	"github.com/bharath-srinivas/giterm/config"
)

type GitApp struct {
	app     *tview.Application
	Context context.Context
	Client  *github.Client

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
		app:     app,
		config:  cfg,
	}
}
