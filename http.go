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

// HTTPServer is a JSON-RPC 2.0 server that handles HTTP requests.
type HTTPServer struct {
	handlers map[string]Handler
	mux      *http.ServeMux
	server   *http.Server
}

// NewHTTPServer creates a new [HTTPServer] with an empty handlers.
func NewHTTPServer(addr, path string) *HTTPServer {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	s := &HTTPServer{
		handlers: make(map[string]Handler),
		mux:      mux,
		server:   server,
	}

	// Register the JSON-RPC handler on the specified path
	mux.HandleFunc(path, s.handleJSONRPC)

	return s
}

var _ Server = (*HTTPServer)(nil)

// Register registers a handler for a specific method.
func (s *HTTPServer) Register(method string, handler Handler) {
	s.handlers[method] = handler
}

// Run starts the HTTP server and listens for incoming requests.
func (s *HTTPServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.server.Shutdown(context.Background())
	}()

	return s.server.ListenAndServe()
}

// handleJSONRPC handles incoming JSON-RPC requests over HTTP.
func (s *HTTPServer) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errorResp := NewError(ParseError, "Parse error")
		resp := NewResponse(nil, WithError(*errorResp))
		s.writeResponse(w, resp)
		return
	}

	if handler, exists := s.handlers[req.Method]; exists {
		resp := handler(r.Context(), &req)
		s.writeResponse(w, resp)
	} else {
		err := NewError(MethodNotFound, "Method not found")
		resp := NewResponse(req.ID, WithError(*err))
		s.writeResponse(w, resp)
	}
}

// writeResponse writes a JSON-RPC response to the HTTP response writer.
func (s *HTTPServer) writeResponse(w http.ResponseWriter, resp *Response) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
