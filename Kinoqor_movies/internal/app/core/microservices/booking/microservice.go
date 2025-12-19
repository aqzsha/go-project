package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"movies/internal/app/core/retry"
	"net/http"
	"net/url"
	"path"
)

type HTTPClient struct {
	baseURL *url.URL
	retry   *retry.Retrier
}

func New(baseURL string, retry *retry.Retrier) (*HTTPClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{baseURL: u, retry: retry}, nil
}

func (c *HTTPClient) StoreSeat(ctx context.Context, input StoreSeat) error {
	body, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("StoreSeat failed to marshal input: %w", err)
	}

	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, "/seat/store")

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("StoreUser failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.retry.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("StoreSeat failed: %w", err)
	}
	defer resp.Body.Close()

	var response interface{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("StoreSeat failed to decode response: %w, status: %s", err, resp.Status)
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("StoreSeat failed to: status %s, response body:%s", resp.Status, response)
	}

	return nil
}

func (c *HTTPClient) DeleteSeat(ctx context.Context, id int64) error {
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, fmt.Sprintf("/seat/delete/%d", id))

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return fmt.Errorf("DeleteSeat failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.retry.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteSeat failed: %w", err)
	}
	defer resp.Body.Close()

	var response interface{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("DeleteSeat failed to decode response: %w, status: %s", err, resp.Status)
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("DeleteSeat failed to: status %s, response body:%s", resp.Status, response)
	}

	return nil
}
