package main

import (
	"io"

	"github.com/cli/go-gh/v2/pkg/api"
)

type GitHubClient interface {
	Get(url string, response interface{}) error
	Post(url string, body io.Reader, response interface{}) error
	Patch(url string, body io.Reader, response interface{}) error
}

type githubClientImpl struct {
	client *api.RESTClient
}

func NewGitHubClient(token string) (GitHubClient, error) {
	client, err := api.NewRESTClient(api.ClientOptions{AuthToken: token})
	if err != nil {
		return nil, err
	}
	return &githubClientImpl{client: client}, nil
}

func (c *githubClientImpl) Get(url string, response interface{}) error {
	return c.client.Get(url, response)
}

func (c *githubClientImpl) Post(url string, body io.Reader, response interface{}) error {
	return c.client.Post(url, body, response)
}

func (c *githubClientImpl) Patch(url string, body io.Reader, response interface{}) error {
	return c.client.Patch(url, body, response)
}
