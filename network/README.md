# Docker network extension API

Go handler to create external network extensions for Docker.

## Usage

This library is designed to be integrated in your program.

1. Implement the `network.Driver` interface.
2. Initialize a `network.Handler` with your implementation.
3. Call either `ServeTCP`, `ServeUnix` or `ServeWindows` from the `network.Handler`.
4. On Windows, docker daemon data dir must be provided for ServeTCP and ServeWindows functions.
On Unix, this parameter is ignored.

## Quickstart (using TCP sockets)

Here is a minimalist example to start the network driver on your host machine.
```go
package main

import (
  "fmt"
  "os"

  "github.com/docker/go-plugins-helpers/network"
)

type MyNetworkDriver struct {
  network.Driver
}

func main() {
  d := MyNetworkDriver{}
  h := network.NewHandler(&d)
  err := h.ServeTCP("test_network", "localhost:8080", "", nil)
  if err != nil {
    fmt.Printf("error occurred: %s\n", err.Error())
    os.Exit(1)
  }
}

```

You can test this out by the following:
```bash
$ curl http://localhost:8080/Plugin.Activate
{"Implements": ["NetworkDriver"]}
```

## Further Examples

### Example using Unix sockets:

```go
  import "github.com/docker/go-plugins-helpers/network"

  d := MyNetworkDriver{}
  h := network.NewHandler(d)
  h.ServeUnix("test_network", 0)
```

### Example using Windows named pipes:

```go
import "github.com/docker/go-plugins-helpers/network"
import "github.com/docker/go-plugins-helpers/sdk"

d := MyNetworkDriver{}
h := network.NewHandler(d)

config := sdk.WindowsPipeConfig{
  // open, read, write permissions for everyone 
  // (uses Windows Security Descriptor Definition Language)
  SecurityDescriptor: AllowServiceSystemAdmin,
  InBufferSize:       4096,
  OutBufferSize:      4096,
}

h.ServeWindows("//./pipe/testpipe", "test_network", WindowsDefaultDaemonRootDir(), &config)
```

## Full example plugins

- [docker-ovs-plugin](https://github.com/gopher-net/docker-ovs-plugin) - An Open vSwitch Networking Plugin
