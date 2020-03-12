package modules

import (
	"context"

	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"
	"golang.org/x/oauth2"

	"github.com/bharath-srinivas/giterm/config"
)

type Client struct {
	app    *tview.Application
	client *github.Client
	ctx    context.Context
}

func NewClient(app *tview.Application, config config.Config) *Client {
	ctx := context.Background()
	return &Client{
		app:    app,
		client: githubClient(config, ctx),
		ctx:    ctx,
	}
}

func githubClient(config config.Config, context context.Context) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	oauth2Client := oauth2.NewClient(context, tokenSource)
	return github.NewClient(oauth2Client)
}
