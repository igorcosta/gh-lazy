package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cli/go-gh/v2/pkg/api"
)

type Client struct {
	client *api.RESTClient
}

func NewClient(token string) (*Client, error) {
	client, err := api.NewRESTClient(api.ClientOptions{
		AuthToken: token,
		Headers:   map[string]string{"Accept": "application/vnd.github.v3+json"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) Get(ctx context.Context, url string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	return c.client.Do(req, response)
}

func (c *Client) Post(ctx context.Context, url string, body, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	return c.client.Do(req, response)
}

func (c *Client) Patch(ctx context.Context, url string, body, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	return c.client.Do(req, response)
}
