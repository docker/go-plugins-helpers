# Docker secrets extension API

Go handler to get secrets from external secret stores in Docker.

## Usage

This library is designed to be integrated in your program.

1. Implement the `secrets.Driver` interface.
2. Initialize a `secrets.Handler` with your implementation.
3. Call either `ServeTCP` or `ServeUnix` from the `secrets.Handler`.

### Example using TCP sockets:

```go
  import "github.com/docker/go-plugins-helpers/secrets"

  d := MySecretsDriver{}
  h := secrets.NewHandler(d)
  h.ServeTCP("test_secrets", ":8080")
```

### Example using Unix sockets:

```go
  import "github.com/docker/go-plugins-helpers/secrets"

  d := MySecretsDriver{}
  h := secrets.NewHandler(d)
  h.ServeUnix("test_secrets", 0)
```
