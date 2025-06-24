package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

// HTTP Transport Layer - Server Implementation
func handleConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req jsonrpc2.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		parseErr := jsonrpc2.NewError(jsonrpc2.ParseError, "Parse error")
		response := jsonrpc2.NewResponse(nil, jsonrpc2.WithError(*parseErr))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := processRequest(&req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func runServer() {
	http.HandleFunc("/rpc", handleConnection)
	fmt.Println("JSON-RPC 2.0 HTTP server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	runServer()
}
