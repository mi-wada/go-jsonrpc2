# TCP JSON-RPC 2.0 Example

A simple TCP-based JSON-RPC 2.0 server implementation.

## Usage

### Start the server

```shell
go run main.go -m server
```

The server listens on port 8081 and supports `add` and `subtract` methods.

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
go run main.go -m client -d '{"jsonrpc":"2.0","method":"subtract","params":{"a":15,"b":6},"id":2}'
```

Expected response:

```json
{"jsonrpc":"2.0","result":9,"id":2}
```

#### Invalid method

```bash
go run main.go -m client -d '{"jsonrpc":"2.0","method":"multiply","params":{"a":2,"b":3},"id":3}'
```

Expected response:

```json
{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":3}
```

### Alternative: Using nc (netcat)

You can also test using nc:

```bash
echo '{"jsonrpc":"2.0","method":"add","params":{"a":5,"b":3},"id":1}' | nc localhost 8081
```

### Alternative: Using telnet

You can also test using telnet:

```bash
telnet localhost 8081
```

Then type the JSON requests directly and press Enter.
