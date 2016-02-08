# Docker authorization extension api.

Go handler to create external authorization extensions for Docker.

## Usage

This library is designed to be integrated in your program.

1. Implement the `authorization.Plugin` interface.
2. Initialize a `authorization.Handler` with your implementation.
3. Call either `ServeTCP` or `ServeUnix` from the `authorization.Handler`.

### Example using TCP sockets:

```go
  p := MyAuthZPlugin{}
  h := authorization.NewHandler(p)
  h.ServeTCP("test_plugin", ":8080")
```

### Example using Unix sockets:

```go
  p := MyAuthZPlugin{}
  h := authorization.NewHandler(p)
  h.ServeUnix("root", "test_plugin")
```

## Full example plugins

- https://github.com/runcom/docker-novolume-plugin

## License

MIT
