package main

import (
	"bufio"
	"encoding/json"
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
func processRequest(req *jsonrpc2.Request) *jsonrpc2.Response {
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

func main() {
	runServer()
}
