package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

func handleRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req jsonrpc2.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, nil, int(jsonrpc2.ParseError), "Parse error", nil)
		return
	}

	if req.JSONRPC != "2.0" {
		sendError(w, req.ID, int(jsonrpc2.InvalidRequest), "Invalid Request", nil)
		return
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

	var response *jsonrpc2.Response
	if err != nil {
		response = jsonrpc2.NewResponse(req.ID, jsonrpc2.WithError(*err))
	} else {
		response = jsonrpc2.NewResponse(req.ID, jsonrpc2.WithResult(result))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendError(w http.ResponseWriter, id any, code int, message string, data any) {
	var err *jsonrpc2.Error
	if data != nil {
		err = jsonrpc2.NewError(jsonrpc2.ErrorCode(code), message, jsonrpc2.WithData(data))
	} else {
		err = jsonrpc2.NewError(jsonrpc2.ErrorCode(code), message)
	}

	response := jsonrpc2.NewResponse(id, jsonrpc2.WithError(*err))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func runServer() {
	http.HandleFunc("/rpc", handleRPC)
	fmt.Println("JSON-RPC 2.0 HTTP server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func callRPC(url string, method string, params any, id any) (*jsonrpc2.Response, error) {
	req, err := jsonrpc2.NewRequest(method, jsonrpc2.WithParams(params), jsonrpc2.WithID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp jsonrpc2.Response
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, err
	}

	return &rpcResp, nil
}

func runClient() {
	serverURL := "http://localhost:8080/rpc"

	fmt.Println("JSON-RPC 2.0 HTTP Client Example")
	fmt.Println("=================================")

	// Test add method
	fmt.Println("\n1. Testing add method (5 + 3):")
	addParams := AddParams{A: 5, B: 3}
	resp, err := callRPC(serverURL, "add", addParams, 1)
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
	fmt.Println("\n2. Testing subtract method (10 - 4):")
	subtractParams := SubtractParams{A: 10, B: 4}
	resp, err = callRPC(serverURL, "subtract", subtractParams, 2)
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
	resp, err = callRPC(serverURL, "multiply", nil, 3)
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
	resp, err = callRPC(serverURL, "add", "invalid", 4)
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
		fmt.Printf("Usage: %s -mode server|client\n", "go run examples/http/main.go")
		flag.PrintDefaults()
	}
}
