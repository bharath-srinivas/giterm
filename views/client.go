package views

import (
	"context"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"

	"github.com/bharath-srinivas/giterm/config"
)

type Client struct {
	Username string

	*github.Client
	context.Context
}

func NewClient(config config.Config) *Client {
	ctx := context.Background()
	client := &Client{
		Client:  githubClient(config, ctx),
		Context: ctx,
	}
	client.Username = client.getUsername()
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
	user, _, err := c.Client.Users.Get(c.Context, "")
	if err != nil || user == nil || user.Login == nil {
		return ""
	}
	return *user.Login
}
