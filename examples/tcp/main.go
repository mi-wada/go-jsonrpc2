package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/mi-wada/go-jsonrpc2"
)

type AddParams struct {
	A int `json:"a"`
	B int `json:"b"`
}

type SubtractParams struct {
	A int `json:"a"`
	B int `json:"b"`
}

// processRequest handles JSON-RPC request processing (common logic)
func processRequest(req jsonrpc2.Request) *jsonrpc2.Response {
	if req.JSONRPC != "2.0" {
		err := jsonrpc2.NewError(jsonrpc2.InvalidRequest, "Invalid Request")
		return jsonrpc2.NewResponse(req.ID, jsonrpc2.WithError(*err))
	}

	var result any
	var err *jsonrpc2.Error

	switch req.Method {
	case "add":
		var params AddParams
		if req.Params != nil {
			if jsonErr := json.Unmarshal(req.Params, &params); jsonErr != nil {
				err = jsonrpc2.NewError(jsonrpc2.InvalidParams, "Invalid params")
			} else {
				result = params.A + params.B
			}
		} else {
			err = jsonrpc2.NewError(jsonrpc2.InvalidParams, "Invalid params")
		}
	case "subtract":
		var params SubtractParams
		if req.Params != nil {
			if jsonErr := json.Unmarshal(req.Params, &params); jsonErr != nil {
				err = jsonrpc2.NewError(jsonrpc2.InvalidParams, "Invalid params")
			} else {
				result = params.A - params.B
			}
		} else {
			err = jsonrpc2.NewError(jsonrpc2.InvalidParams, "Invalid params")
		}
	default:
		err = jsonrpc2.NewError(jsonrpc2.MethodNotFound, "Method not found")
	}

	if err != nil {
		return jsonrpc2.NewResponse(req.ID, jsonrpc2.WithError(*err))
	} else {
		return jsonrpc2.NewResponse(req.ID, jsonrpc2.WithResult(result))
	}
}

// TCP Transport Layer - Server Implementation
func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	encoder := json.NewEncoder(conn)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req jsonrpc2.Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendError(encoder, nil, int(jsonrpc2.ParseError), "Parse error", nil)
			continue
		}

		response := processRequest(req)
		if encErr := encoder.Encode(response); encErr != nil {
			log.Printf("Error encoding response: %v", encErr)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
}

func sendError(encoder *json.Encoder, id any, code int, message string, data any) {
	var err *jsonrpc2.Error
	if data != nil {
		err = jsonrpc2.NewError(jsonrpc2.ErrorCode(code), message, jsonrpc2.WithData(data))
	} else {
		err = jsonrpc2.NewError(jsonrpc2.ErrorCode(code), message)
	}

	response := jsonrpc2.NewResponse(id, jsonrpc2.WithError(*err))
	encoder.Encode(response)
}

func runServer() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Error starting TCP server:", err)
	}
	defer listener.Close()

	fmt.Println("JSON-RPC 2.0 TCP server starting on :8081")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

// TCP Transport Layer - Client Implementation
func callRPC(conn net.Conn, method string, params any, id any) (*jsonrpc2.Response, error) {
	req, err := jsonrpc2.NewRequest(method, jsonrpc2.WithParams(params), jsonrpc2.WithID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(req); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return nil, fmt.Errorf("no response from server")
	}

	var resp jsonrpc2.Response
	if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func runClient() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	fmt.Println("JSON-RPC 2.0 TCP Client Example")
	fmt.Println("================================")

	// Test add method
	fmt.Println("\n1. Testing add method (7 + 2):")
	addParams := AddParams{A: 7, B: 2}
	resp, err := callRPC(conn, "add", addParams, 1)
	if err != nil {
		log.Printf("Error calling add: %v", err)
	} else {
		if resp.Error != nil {
			fmt.Printf("RPC Error: %+v\n", resp.Error)
		} else {
			var result int
			resultBytes, _ := json.Marshal(resp.Result)
			json.Unmarshal(resultBytes, &result)
			fmt.Printf("Result: %d\n", result)
		}
	}

	// Test subtract method
	fmt.Println("\n2. Testing subtract method (15 - 6):")
	subtractParams := SubtractParams{A: 15, B: 6}
	resp, err = callRPC(conn, "subtract", subtractParams, 2)
	if err != nil {
		log.Printf("Error calling subtract: %v", err)
	} else {
		if resp.Error != nil {
			fmt.Printf("RPC Error: %+v\n", resp.Error)
		} else {
			var result int
			resultBytes, _ := json.Marshal(resp.Result)
			json.Unmarshal(resultBytes, &result)
			fmt.Printf("Result: %d\n", result)
		}
	}

	// Test invalid method
	fmt.Println("\n3. Testing invalid method:")
	resp, err = callRPC(conn, "multiply", nil, 3)
	if err != nil {
		log.Printf("Error calling invalid method: %v", err)
	} else {
		if resp.Error != nil {
			fmt.Printf("RPC Error: %+v\n", resp.Error)
		} else {
			fmt.Printf("Unexpected success: %+v\n", resp)
		}
	}

	// Test invalid params
	fmt.Println("\n4. Testing invalid params:")
	resp, err = callRPC(conn, "add", "invalid", 4)
	if err != nil {
		log.Printf("Error calling with invalid params: %v", err)
	} else {
		if resp.Error != nil {
			fmt.Printf("RPC Error: %+v\n", resp.Error)
		} else {
			fmt.Printf("Unexpected success: %+v\n", resp)
		}
	}
}

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "server", "Mode: server or client")
	flag.StringVar(&mode, "m", "server", "Mode: server or client (shorthand)")
	flag.Parse()

	switch mode {
	case "server":
		runServer()
	case "client":
		runClient()
	default:
		fmt.Printf("Usage: %s -mode server|client\n", "go run examples/tcp/main.go")
		flag.PrintDefaults()
	}
}
