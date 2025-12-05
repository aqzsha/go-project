package microservices

import (
	"context"
	"net/http"
	"net/url"
	"path"
)

type Client interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
	FullPath(url string) string
}

type BaseClient struct {
	Client  *http.Client
	baseURL *url.URL
	apiKey  string
}

func NewBaseClient(client *http.Client, baseURL, apiKey string) (*BaseClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &BaseClient{
		baseURL: u,
		apiKey:  apiKey,
		Client:  client,
	}, nil
}

func (c *BaseClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Set("x-api-key", c.apiKey)

	return c.Client.Do(req)
}

func (c *BaseClient) FullPath(url string) string {
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, url)

	return u.String()
}
