package github

import (
	"context"
	"fmt"
	"io"

	"github.com/cli/go-gh/v2/pkg/api"
)

type Client struct {
	client *api.RESTClient
}

func NewClient(token string) (*Client, error) {
	client, err := api.NewRESTClient(api.ClientOptions{
		AuthToken: token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) Get(ctx context.Context, path string, response interface{}) error {
	return c.client.Get(path, response)
}

func (c *Client) Post(ctx context.Context, path string, body io.Reader, response interface{}) error {
	return c.client.Post(path, body, response)
}

func (c *Client) Patch(ctx context.Context, path string, body io.Reader, response interface{}) error {
	return c.client.Patch(path, body, response)
}

func (c *Client) GetUsername() (string, error) {
	var response struct {
		Login string `json:"login"`
	}
	err := c.Get(context.Background(), "user", &response)
	if err != nil {
		return "", fmt.Errorf("failed to get username: %w", err)
	}
	return response.Login, nil
}
