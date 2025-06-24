// Package jsonrpc2 provides structs and functions for working with JSON-RPC 2.0 protocol.
// https://www.jsonrpc.org/specification
package jsonrpc2

import (
	"context"
	"encoding/json"
	"fmt"
)

const JSONRPC = "2.0" // JSONRPC is the version of the JSON-RPC protocol.

// Request represents a JSON-RPC 2.0 request object.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`          // The version of the JSON-RPC protocol. It must be "2.0".
	Method  string          `json:"method"`           // The name of the method to be invoked.
	Params  json.RawMessage `json:"params,omitempty"` // The parameters of the method being invoked.
	ID      any             `json:"id"`               // A unique identifier for the request.
}

// UnmarshalRequest unmarshals a [Request] from JSON data.
func UnmarshalRequest(data []byte) (*Request, error) {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}
	return &req, nil
}

// NewRequest creates a new [Request].
// If you want to set the Params or ID fields, use the [WithParams] or [WithID] options.
func NewRequest(method string, opts ...NewRequestOption) (*Request, error) {
	req := &Request{
		JSONRPC: JSONRPC,
		Method:  method,
	}
	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// NewRequestOption defines a function type for setting optional fields in [Request].
type NewRequestOption func(*Request) error

// WithParams sets the Params field of a [Request].
func WithParams(params any) NewRequestOption {
	return func(r *Request) error {
		if params == nil {
			return nil
		}
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("failed to marshal params: %w", err)
		}
		r.Params = paramsJSON
		return nil
	}
}

// WithID sets the ID field of a [Request].
func WithID(id any) NewRequestOption {
	return func(r *Request) error {
		r.ID = id
		return nil
	}
}

// Response represents a JSON-RPC 2.0 response object.
type Response struct {
	JSONRPC string `json:"jsonrpc"`          // The version of the JSON-RPC protocol. It must be "2.0".
	Result  any    `json:"result,omitempty"` // The result of the method invocation. This field is omitted if there was an error.
	Error   *Error `json:"error,omitempty"`  // An error object if an error occurred.
	ID      any    `json:"id"`               // The same ID as in the request. It is used to match responses to requests.
}

// NewResponse creates a new [Response].
// If you want to set the Result or Error fields, use the [WithResult] or [WithError] options.
func NewResponse(id any, opts ...NewResponseOption) *Response {
	resp := &Response{
		JSONRPC: JSONRPC,
		ID:      id,
	}
	for _, opt := range opts {
		opt(resp)
	}
	return resp
}

// NewResponseOption defines a function type for setting optional fields in [Response].
type NewResponseOption func(*Response)

// WithResult sets the Result field of a [Response].
func WithResult(result any) NewResponseOption {
	return func(r *Response) {
		r.Result = result
	}
}

// WithError sets the Error field of a [Response].
func WithError(err Error) NewResponseOption {
	return func(r *Response) {
		r.Error = &err
	}
}

// ErrorCode represents the error codes as defined in the JSON-RPC 2.0 specification.
// For more details, see: https://www.jsonrpc.org/specification#error_object
type ErrorCode int

const (
	ParseError     ErrorCode = -32700 // Invalid JSON was received by the server.
	InvalidRequest ErrorCode = -32600 // The JSON sent is not a valid Request object.
	MethodNotFound ErrorCode = -32601 // The method does not exist / is not available.
	InvalidParams  ErrorCode = -32602 // Invalid method parameter(s).
	InternalError  ErrorCode = -32603 // Internal JSON-RPC error.
)

// Error represents a JSON-RPC 2.0 error object.
type Error struct {
	Code    ErrorCode `json:"code"`           // A number indicating the error type that occurred
	Message string    `json:"message"`        // A short description of the error
	Data    any       `json:"data,omitempty"` // Additional information about the error
}

// NewError creates a new [Error].
// If you want to set the Data field, use the [WithData] option.
func NewError(code ErrorCode, message string, opts ...NewErrorOption) *Error {
	err := &Error{
		Code:    code,
		Message: message,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

// NewErrorOption defines a function type for setting optional fields in Error.
type NewErrorOption func(*Error)

// WithData sets the Data field of an Error.
func WithData(data any) NewErrorOption {
	return func(e *Error) {
		e.Data = data
	}
}

// Error implements the [Error] interface.
func (e Error) Error() string {
	return fmt.Sprintf("JSON-RPC Error %d: %s", e.Code, e.Message)
}

// Client is an interface for making JSON-RPC 2.0 requests.
type Client interface {
	// Call sends a JSON-RPC 2.0 request and returns the response.
	Call(ctx context.Context, req *Request) (*Response, error)
	// CallBatch sends multiple JSON-RPC 2.0 requests at once and returns their responses.
	// If a ParseError occurs, returns a single [Response]. Otherwise, returns a slice of [Response].
	CallBatch(ctx context.Context, reqs []*Request) (any, error)
	// Notify sends a JSON-RPC 2.0 notification (no response expected).
	Notify(ctx context.Context, req *Request) error
}

// Handler is a function type that processes a JSON-RPC request and returns a response.
type Handler = func(ctx context.Context, req *Request) *Response

// Server is an interface for handling JSON-RPC 2.0 requests.
type Server interface {
	// Register registers a handler for a specific method.
	Register(method string, handler Handler)
	// Run starts the server and listens for incoming requests.
	Run(ctx context.Context) error
}
