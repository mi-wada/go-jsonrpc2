package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
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

// HTTP Transport Layer - Server Implementation
func handleHTTP(w http.ResponseWriter, r *http.Request) {
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

	response := handleRequest(&req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func runServer() {
	http.HandleFunc("/rpc", handleHTTP)
	fmt.Println("JSON-RPC 2.0 HTTP server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
