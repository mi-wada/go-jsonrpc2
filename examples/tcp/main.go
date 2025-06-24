package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mi-wada/go-jsonrpc2"
)

func main() {
	// Create and configure the TCP server
	server := jsonrpc2.NewTCPServer(":8081")
	server.Register("add", addHandler)
	server.Register("subtract", subtractHandler)

	// Run the server
	ctx := context.Background()
	if err := server.Run(ctx); err != nil {
		log.Printf("Server error: %v", err)
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
