# HTTP JSON-RPC 2.0 Example

A simple HTTP-based JSON-RPC 2.0 server implementation.

## Usage

### Start the server

```shell
go run main.go -m server
```

The server listens on port 8080 and supports `add` and `subtract` methods.

### Send requests using the built-in client

#### Add method

```bash
go run main.go -m client -d '{"jsonrpc":"2.0","method":"add","params":{"a":5,"b":3},"id":1}'
```

Expected response:

```json
{"jsonrpc":"2.0","result":8,"id":1}
```

#### Subtract method

```bash
go run main.go -m client -d '{"jsonrpc":"2.0","method":"subtract","params":{"a":10,"b":4},"id":2}'
```

Expected response:

```json
{"jsonrpc":"2.0","result":6,"id":2}
```

#### Invalid method

```bash
go run main.go -m client -d '{"jsonrpc":"2.0","method":"multiply","params":{"a":2,"b":3},"id":3}'
```

Expected response:

```json
{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":3}
```

### Alternative: Using curl

You can also test using curl:

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"add","params":{"a":5,"b":3},"id":1}'
```
