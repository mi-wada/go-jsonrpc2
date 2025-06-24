package jsonrpc2

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
)

// TCPClient is a JSON-RPC 2.0 client that communicates over TCP.
type TCPClient struct {
	conn net.Conn
}

// NewTCPClient creates a new [TCPClient].
func NewTCPClient(conn net.Conn) *TCPClient {
	return &TCPClient{
		conn: conn,
	}
}

var _ Client = (*TCPClient)(nil)

// Call sends a JSON-RPC 2.0 request over TCP and returns the response.
func (c *TCPClient) Call(ctx context.Context, req *Request) (*Response, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if err := c.conn.SetDeadline(deadline); err != nil {
			return nil, fmt.Errorf("failed to set connection deadline: %w", err)
		}
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if _, err := c.conn.Write(append(reqData, '\n')); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	scanner := bufio.NewScanner(c.conn)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
		return nil, fmt.Errorf("connection closed without response")
	}

	var rpcResp Response
	if err := json.Unmarshal(scanner.Bytes(), &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &rpcResp, nil
}

// CallBatch sends a batch of JSON-RPC requests over TCP and returns the responses.
func (c *TCPClient) CallBatch(ctx context.Context, reqs []*Request) (any, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if err := c.conn.SetDeadline(deadline); err != nil {
			return nil, fmt.Errorf("failed to set connection deadline: %w", err)
		}
	}

	reqData, err := json.Marshal(reqs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal batch request: %w", err)
	}

	if _, err := c.conn.Write(append(reqData, '\n')); err != nil {
		return nil, fmt.Errorf("failed to send batch request: %w", err)
	}

	scanner := bufio.NewScanner(c.conn)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read batch response: %w", err)
		}
		return nil, fmt.Errorf("connection closed without response")
	}

	var rpcResp any
	if err := json.Unmarshal(scanner.Bytes(), &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to decode batch response: %w", err)
	}

	return rpcResp, nil
}

// Notify sends a JSON-RPC notification over TCP.
func (c *TCPClient) Notify(ctx context.Context, req *Request) error {
	if deadline, ok := ctx.Deadline(); ok {
		if err := c.conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("failed to set connection deadline: %w", err)
		}
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	if _, err := c.conn.Write(append(reqData, '\n')); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	return nil
}

// TCPServer is a JSON-RPC 2.0 server that handles TCP connections.
type TCPServer struct {
	handlers map[string]Handler
	addr     string
}

// NewTCPServer creates a new [TCPServer] with an empty handlers.
func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		handlers: make(map[string]Handler),
		addr:     addr,
	}
}

var _ Server = (*TCPServer)(nil)

// Register registers a handler for a specific method.
func (s *TCPServer) Register(method string, handler Handler) {
	s.handlers[method] = handler
}

// Run starts the TCP server and listens for incoming connections.
func (s *TCPServer) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}
	defer listener.Close()

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				continue
			}
		}
		go s.handleConnection(ctx, conn)
	}
}

// handleConnection handles a single TCP connection.
func (s *TCPServer) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	encoder := json.NewEncoder(conn)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			errorResp := NewError(ParseError, "Parse error")
			resp := NewResponse(nil, WithError(*errorResp))
			encoder.Encode(resp)
			continue
		}

		if handler, exists := s.handlers[req.Method]; exists {
			resp := handler(ctx, &req)
			encoder.Encode(resp)
		} else {
			err := NewError(MethodNotFound, "Method not found")
			resp := NewResponse(req.ID, WithError(*err))
			encoder.Encode(resp)
		}
	}
}
