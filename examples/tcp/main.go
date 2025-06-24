package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

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

// handleRequest handles JSON-RPC request processing (common logic)
func handleRequest(req *jsonrpc2.Request) *jsonrpc2.Response {
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
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		req, err := jsonrpc2.UnmarshalRequest(line)
		if err != nil {
			parseErr := jsonrpc2.NewError(jsonrpc2.ParseError, "Parse error")
			response := jsonrpc2.NewResponse(nil, jsonrpc2.WithError(*parseErr))
			if encErr := encoder.Encode(response); encErr != nil {
				log.Printf("Error encoding parse error response: %v", encErr)
			}
			continue
		}

		response := handleRequest(req)
		if encErr := encoder.Encode(response); encErr != nil {
			log.Printf("Error encoding response: %v", encErr)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
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

func runClient(data string) {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	client := jsonrpc2.NewTCPClient(conn)

	req, err := jsonrpc2.UnmarshalRequest([]byte(data))
	if err != nil {
		log.Fatal("Error parsing request data:", err)
	}

	resp, err := client.Call(context.Background(), req)
	if err != nil {
		log.Fatal("Error calling RPC:", err)
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("Error marshaling response:", err)
	}

	fmt.Println(string(respData))
}

func main() {
	var mode = flag.String("m", "server", "Mode: server or client")
	var data = flag.String("d", "", "JSON-RPC request data (for client mode)")
	flag.Parse()

	if *mode == "client" {
		if *data == "" {
			fmt.Println("Usage: go run main.go -m client -d '<json_data>'")
			os.Exit(1)
		}
		runClient(*data)
	} else {
		runServer()
	}
}
