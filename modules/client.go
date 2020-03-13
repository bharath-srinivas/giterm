package modules

import (
	"context"

	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"
	"golang.org/x/oauth2"

	"github.com/bharath-srinivas/giterm/config"
)

type Client struct {
	app      *tview.Application
	client   *github.Client
	ctx      context.Context
	username string
}

func NewClient(app *tview.Application, config config.Config) *Client {
	ctx := context.Background()
	client := &Client{
		app:    app,
		client: githubClient(config, ctx),
		ctx:    ctx,
	}
	client.username = client.getUsername()
	return client
}

func githubClient(config config.Config, context context.Context) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	oauth2Client := oauth2.NewClient(context, tokenSource)
	return github.NewClient(oauth2Client)
}

func (c *Client) getUsername() string {
	user, _, err := c.client.Users.Get(c.ctx, "")
	if err != nil || user == nil || user.Login == nil {
		return ""
	}
	return *user.Login
}
