# go-jsonrpc2

Package go-jsonrpc2 provides server and client implementations for the [JSON-RPC 2.0 protocol](https://www.jsonrpc.org/specification).

## Install

```shell
go get github.com/mi-wada/go-jsonrpc2@latest
```

## Usage

See the [examples](https://github.com/mi-wada/go-jsonrpc2/tree/main/examples) for usage.

## LICENSE

MIT

## ToDo

- [ ] Add tests
- [ ] It would be useful if HTTP could also be used as a Handler. For example, Server.HTTPHandler() could return an http.Handler.
- [ ] maybe it's good to add NextID() func to client. Generate random or sequential ID.
- [ ] Consider logging strategy at server
