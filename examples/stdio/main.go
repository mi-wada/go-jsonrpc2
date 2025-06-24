package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/mi-wada/go-jsonrpc2"
)

type addParams struct {
	A int `json:"a"`
	B int `json:"b"`
}

type subtractParams struct {
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

func main() {
	log.SetOutput(os.Stderr) // Set log output to stderr

	// Create and configure the stdio server
	server := jsonrpc2.NewStdioServer()
	server.Register("add", addHandler)
	server.Register("subtract", subtractHandler)

	// Run the server
	ctx := context.Background()
	if err := server.Run(ctx); err != nil {
		log.Printf("Server error: %v", err)
	}
}
