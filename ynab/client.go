package ynab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseURL = "https://api.ynab.com/v1"

// Client is an HTTP client for the YNAB API.
type Client struct {
	token      string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new YNAB API client.
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

// NewTestClient creates a client with a custom base URL, for testing.
func NewTestClient(token, baseURL string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func (c *Client) doRequest(method, path string, body any) ([]byte, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("building URL: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiResp Response[json.RawMessage]
		if err := json.Unmarshal(data, &apiResp); err == nil && apiResp.Error != nil {
			return nil, apiResp.Error
		}
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(data))
	}

	return data, nil
}

func doGet[T any](c *Client, path string) (T, error) {
	var zero T
	data, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return zero, err
	}

	var resp Response[T]
	if err := json.Unmarshal(data, &resp); err != nil {
		return zero, fmt.Errorf("decoding response: %w", err)
	}
	return resp.Data, nil
}

func doPost[T any](c *Client, path string, body any) (T, error) {
	var zero T
	data, err := c.doRequest(http.MethodPost, path, body)
	if err != nil {
		return zero, err
	}

	var resp Response[T]
	if err := json.Unmarshal(data, &resp); err != nil {
		return zero, fmt.Errorf("decoding response: %w", err)
	}
	return resp.Data, nil
}

func doPut[T any](c *Client, path string, body any) (T, error) {
	var zero T
	data, err := c.doRequest(http.MethodPut, path, body)
	if err != nil {
		return zero, err
	}

	var resp Response[T]
	if err := json.Unmarshal(data, &resp); err != nil {
		return zero, fmt.Errorf("decoding response: %w", err)
	}
	return resp.Data, nil
}
