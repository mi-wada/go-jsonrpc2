# HTTP JSON-RPC 2.0 Example

A simple HTTP-based JSON-RPC 2.0 server and client implementation.

## Usage

Start the server:

```shell
go run main.go -m server
```

Run the client:

```shell
go run main.go -m client
```

The server listens on port 8080 and supports `add` and `subtract` methods.
