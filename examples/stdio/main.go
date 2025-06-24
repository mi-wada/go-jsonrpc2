package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

func runServer() {
	log.SetOutput(os.Stderr) // ログはstderrに出力

	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	log.Println("JSON-RPC 2.0 stdio server started")

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req jsonrpc2.Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			parseErr := jsonrpc2.NewError(jsonrpc2.ParseError, "Parse error")
			response := jsonrpc2.NewResponse(nil, jsonrpc2.WithError(*parseErr))
			if encErr := encoder.Encode(response); encErr != nil {
				log.Printf("Error encoding parse error response: %v", encErr)
			}
			continue
		}

		response := processRequest(req)
		if err := encoder.Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}

	log.Println("JSON-RPC 2.0 stdio server stopped")
}

func runClient() {
	fmt.Println("JSON-RPC 2.0 STDIO Client Example")
	fmt.Println("==================================")
	fmt.Println("Note: This is a demonstration client.")
	fmt.Println("In practice, stdio transport is typically used")
	fmt.Println("for inter-process communication where the server")
	fmt.Println("reads from stdin and writes to stdout.")

	// Test add method
	fmt.Println("\n1. Testing add method (12 + 8):")
	addParams := AddParams{A: 12, B: 8}
	fmt.Printf("Would send: add(%+v) with id=1\n", addParams)
	fmt.Printf("Expected result: %d\n", addParams.A+addParams.B)

	// Test subtract method
	fmt.Println("\n2. Testing subtract method (20 - 7):")
	subtractParams := SubtractParams{A: 20, B: 7}
	fmt.Printf("Would send: subtract(%+v) with id=2\n", subtractParams)
	fmt.Printf("Expected result: %d\n", subtractParams.A-subtractParams.B)

	// Test invalid method
	fmt.Println("\n3. Testing invalid method:")
	fmt.Println("Would send: multiply(null) with id=3")
	fmt.Println("Expected: Method not found error")

	// Test invalid params
	fmt.Println("\n4. Testing invalid params:")
	fmt.Println("Would send: add(\"invalid\") with id=4")
	fmt.Println("Expected: Invalid params error")
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
		fmt.Printf("Usage: %s -mode server|client\n", "go run examples/stdio/main.go")
		flag.PrintDefaults()
	}
}
