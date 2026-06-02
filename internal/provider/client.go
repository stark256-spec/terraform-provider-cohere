package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.cohere.com"

type CohereClient struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func newClient(apiKey, baseURL string) *CohereClient {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &CohereClient{apiKey: apiKey, baseURL: baseURL, http: &http.Client{Timeout: 30 * time.Second}}
}

func (c *CohereClient) do(ctx context.Context, method, path string, body any) ([]byte, int, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request: %w", err)
		}
		buf = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, buf)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}
	return data, resp.StatusCode, nil
}

type Workspace struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (c *CohereClient) CreateWorkspace(ctx context.Context, name string) (*Workspace, error) {
	data, _, err := c.do(ctx, http.MethodPost, "/v2/teams", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	var w Workspace
	return &w, json.Unmarshal(data, &w)
}

func (c *CohereClient) GetWorkspace(ctx context.Context, id string) (*Workspace, error) {
	data, _, err := c.do(ctx, http.MethodGet, "/v2/teams/"+id, nil)
	if err != nil {
		return nil, err
	}
	var w Workspace
	return &w, json.Unmarshal(data, &w)
}

func (c *CohereClient) DeleteWorkspace(ctx context.Context, id string) error {
	_, _, err := c.do(ctx, "/v2/teams/"+id, http.MethodDelete, nil)
	return err
}

type APIKey struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Status     string  `json:"status"`
	PartialKey string  `json:"partial_key"`
	CreatedAt  string  `json:"created_at"`
	Key        *string `json:"key,omitempty"`
}

func (c *CohereClient) CreateAPIKey(ctx context.Context, name string) (*APIKey, error) {
	data, _, err := c.do(ctx, http.MethodPost, "/v2/api-keys", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	var k APIKey
	return &k, json.Unmarshal(data, &k)
}

func (c *CohereClient) GetAPIKey(ctx context.Context, id string) (*APIKey, error) {
	data, _, err := c.do(ctx, http.MethodGet, "/v2/api-keys/"+id, nil)
	if err != nil {
		return nil, err
	}
	var k APIKey
	return &k, json.Unmarshal(data, &k)
}

func (c *CohereClient) DeleteAPIKey(ctx context.Context, id string) error {
	_, _, err := c.do(ctx, http.MethodDelete, "/v2/api-keys/"+id, nil)
	return err
}
