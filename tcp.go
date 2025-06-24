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
