# STDIO JSON-RPC 2.0 Example

A simple standard input/output based JSON-RPC 2.0 server implementation.

## Usage

Start the server:

```shell
go run main.go
```

The server reads JSON-RPC requests from stdin and writes responses to stdout. Logs are written to stderr.

## Example Requests

Send these JSON requests via stdin (one per line):

### Add method

```json
{"jsonrpc":"2.0","method":"add","params":{"a":5,"b":3},"id":1}
```

Expected response:

```json
{"jsonrpc":"2.0","result":8,"id":1}
```

### Subtract method

```json
{"jsonrpc":"2.0","method":"subtract","params":{"a":10,"b":4},"id":2}
```

Expected response:

```json
{"jsonrpc":"2.0","result":6,"id":2}
```

### Invalid method

```json
{"jsonrpc":"2.0","method":"multiply","params":{"a":2,"b":3},"id":3}
```

Expected response:

```json
{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":3}
```

## Testing with echo

You can test using echo and pipes:

```shell
echo '{"jsonrpc":"2.0","method":"add","params":{"a":5,"b":3},"id":1}' | go run main.go
```
