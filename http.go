package jsonrpc2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	defaultHTTPMethod    = http.MethodPost
	defaultSuccessStatus = http.StatusOK
)

// HTTPClient is a JSON-RPC 2.0 client that communicates over HTTP.
type HTTPClient struct {
	endpoint string
	client   *http.Client
}

// NewHTTPClient creates a new [HTTPClient].
func NewHTTPClient(endpoint string, client *http.Client) *HTTPClient {
	if client == nil {
		client = &http.Client{}
	}
	return &HTTPClient{
		endpoint: endpoint,
		client:   client,
	}
}

var _ Client = (*HTTPClient)(nil)

// Call sends a JSON-RPC request over HTTP and returns the response.
func (c *HTTPClient) Call(ctx context.Context, req *Request) (*Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != defaultSuccessStatus {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var rpcResp Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&rpcResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &rpcResp, nil
}

// CallBatch sends a batch of JSON-RPC requests over HTTP and returns the responses.
func (c *HTTPClient) CallBatch(ctx context.Context, reqs []*Request) (any, error) {
	body, err := json.Marshal(reqs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal batch request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, defaultHTTPMethod, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != defaultSuccessStatus {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var rpcResp any
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&rpcResp); err != nil {
		return nil, fmt.Errorf("failed to decode response as single or batch: %w", err)
	}
	return rpcResp, nil
}

// Notify sends a JSON-RPC notification over HTTP.
func (c *HTTPClient) Notify(ctx context.Context, req *Request) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, defaultHTTPMethod, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
