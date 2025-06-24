# TCP JSON-RPC 2.0 Example

A simple TCP-based JSON-RPC 2.0 server implementation.

## Usage

Start the server:

```shell
go run main.go
```

The server listens on port 8081 and supports `add` and `subtract` methods.

## Example Requests

Send these JSON requests via TCP connection to `localhost:8081`:

### Add method

```bash
echo '{"jsonrpc":"2.0","method":"add","params":{"a":7,"b":2},"id":1}' | nc localhost 8081
```

Expected response:

```json
{"jsonrpc":"2.0","result":9,"id":1}
```

### Subtract method

```bash
echo '{"jsonrpc":"2.0","method":"subtract","params":{"a":15,"b":6},"id":2}' | nc localhost 8081
```

Expected response:

```json
{"jsonrpc":"2.0","result":9,"id":2}
```

### Invalid method

```bash
echo '{"jsonrpc":"2.0","method":"multiply","params":{"a":2,"b":3},"id":3}' | nc localhost 8081
```

Expected response:

```json
{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":3}
```

### Testing with telnet

You can also test using telnet:

```bash
telnet localhost 8081
```

Then type the JSON requests directly and press Enter.
