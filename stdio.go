package jsonrpc2

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
)

// StdioServer is a JSON-RPC 2.0 server that reads requests from standard input and writes responses to standard output.
type StdioServer struct {
	handlers map[string]Handler
}

// NewStdioServer creates a new [StdioServer] with an empty handlers.
func NewStdioServer() *StdioServer {
	return &StdioServer{
		handlers: make(map[string]Handler),
	}
}

var _ Server = (*StdioServer)(nil)

// Register registers a handler for a specific method.
func (s *StdioServer) Register(method string, handler Handler) {
	s.handlers[method] = handler
}

// Run starts the server, reading requests from standard input and writing responses to standard output.
func (s *StdioServer) Run(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	log.Println("JSON-RPC 2.0 stdio server started")

	for scanner.Scan() {
		line := scanner.Text()
		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("Error unmarshalling request: %v", err)
			continue
		}

		if handler, exists := s.handlers[req.Method]; exists {
			resp := handler(ctx, &req)
			if err := encoder.Encode(resp); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
		} else {
			err := NewError(MethodNotFound, "Method not found")
			resp := NewResponse(req.ID, WithError(*err))
			if err := encoder.Encode(resp); err != nil {
				log.Printf("Error encoding error response: %v", err)
			}
		}
	}

	return scanner.Err()
}
