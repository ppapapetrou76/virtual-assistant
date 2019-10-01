package github

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v27/github"
	"golang.org/x/oauth2"
)

const (
	// RepoEnvVar represents the environment variable GITHUB_REPOSITORY
	RepoEnvVar = "GITHUB_REPOSITORY"
	// ShaEnvVar represents the environment variable GITHUB_SHA
	ShaEnvVar = "GITHUB_SHA"
	// EventNameEnvVar represents the environment variable GITHUB_EVENT_NAME
	EventNameEnvVar = "GITHUB_EVENT_NAME"
	// TokenEnvVar represents the environment variable GITHUB_TOKEN
	TokenEnvVar = "GITHUB_TOKEN"
	// EventPathEnvVar represents the environment variable GITHUB_EVENT_PATH
	EventPathEnvVar = "GITHUB_EVENT_PATH"
	// InputConfigPathEnvVar represents the environment variable INPUT_CONFIG_PATH
	InputConfigPathEnvVar = "INPUT_CONFIG_PATH"
)

// ClientWrapper wraps the github client
type ClientWrapper struct {
	*github.Client
}

// Client returns a github client wrapper with the given http client
func Client(client *http.Client) ClientWrapper {
	return ClientWrapper{Client: github.NewClient(client)}
}

// DefaultClient returns the default client wrapper with an Oath2 ready http client
func DefaultClient() ClientWrapper {
	ghToken := os.Getenv(TokenEnvVar)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return Client(tc)
}
