package views

import (
	"context"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"

	"github.com/bharath-srinivas/giterm/config"
)

// Client represents the github client.
type Client struct {
	Username string

	*github.Client
	context.Context
}

// NewClient returns a new client with the provided token.
func NewClient(config config.Config) *Client {
	ctx := context.Background()
	client := &Client{
		Client:  githubClient(ctx, config),
		Context: ctx,
	}
	client.Username = client.getUsername()
	return client
}

// githubClient returns a new github client with the provided token and context.
func githubClient(context context.Context, config config.Config) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	oauth2Client := oauth2.NewClient(context, tokenSource)
	return github.NewClient(oauth2Client)
}

// getUsername returns the username of the current user.
func (c *Client) getUsername() string {
	user, _, err := c.Client.Users.Get(c.Context, "")
	if err != nil || user == nil || user.Login == nil {
		return ""
	}
	return *user.Login
}
