package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mi-wada/go-jsonrpc2"
)

func runServer() {
	// Create and configure the HTTP server
	server := jsonrpc2.NewHTTPServer(":8080", "/rpc")
	server.Register("add", addHandler)
	server.Register("subtract", subtractHandler)

	fmt.Println("JSON-RPC 2.0 HTTP server starting on :8080")

	// Run the server
	ctx := context.Background()
	if err := server.Run(ctx); err != nil {
		log.Printf("Server error: %v", err)
	}
}

func runClient(data string) {
	client := jsonrpc2.NewHTTPClient("http://localhost:8080/rpc", nil)

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

type addParams struct {
	A int `json:"a"`
	B int `json:"b"`
}

func addHandler(ctx context.Context, req *jsonrpc2.Request) *jsonrpc2.Response {
	if req.Params == nil {
		jsonErr := jsonrpc2.NewError(jsonrpc2.InvalidParams, "Invalid params")
		return jsonrpc2.NewResponse(req.ID, jsonrpc2.WithError(*jsonErr))
	}

	var params addParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		jsonErr := jsonrpc2.NewError(jsonrpc2.InvalidParams, "Invalid params")
		return jsonrpc2.NewResponse(req.ID, jsonrpc2.WithError(*jsonErr))
	}
	result := params.A + params.B
	return jsonrpc2.NewResponse(req.ID, jsonrpc2.WithResult(result))
}

type subtractParams struct {
	A int `json:"a"`
	B int `json:"b"`
}

func subtractHandler(ctx context.Context, req *jsonrpc2.Request) *jsonrpc2.Response {
	if req.Params == nil {
		jsonErr := jsonrpc2.NewError(jsonrpc2.InvalidParams, "Invalid params")
		return jsonrpc2.NewResponse(req.ID, jsonrpc2.WithError(*jsonErr))
	}
	var params subtractParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		jsonErr := jsonrpc2.NewError(jsonrpc2.InvalidParams, "Invalid params")
		return jsonrpc2.NewResponse(req.ID, jsonrpc2.WithError(*jsonErr))
	}
	result := params.A - params.B
	return jsonrpc2.NewResponse(req.ID, jsonrpc2.WithResult(result))
}
