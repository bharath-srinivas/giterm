package views

import (
	"context"

	"github.com/google/go-github/v30/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/bharath-srinivas/giterm/config"
)

// Client represents the github client.
type Client struct {
	GqlClient *githubv4.Client

	*github.Client
	context.Context
}

// NewClient returns a new client with the provided token.
func NewClient(config config.Config) *Client {
	ctx := context.Background()
	client := &Client{
		Context:   ctx,
		Client:    githubV3Client(ctx, config),
		GqlClient: githubV4Client(ctx, config),
	}
	return client
}

// githubV3Client returns a new github API client with the provided token and context.
func githubV3Client(context context.Context, config config.Config) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	oauth2Client := oauth2.NewClient(context, tokenSource)
	return github.NewClient(oauth2Client)
}

// githubV4Client returns a new github graphql client with the provided token and context.
func githubV4Client(context context.Context, config config.Config) *githubv4.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	oauth2Client := oauth2.NewClient(context, tokenSource)
	return githubv4.NewClient(oauth2Client)
}
