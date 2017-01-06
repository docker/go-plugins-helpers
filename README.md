# go-plugins-helpers

A collection of helper packages to extend Docker Engine in Go

 Plugin type   | Documentation | Description
 --------------|---------------|--------------------------------------------------
 Authorization | [Link](https://docs.docker.com/engine/extend/authorization/)   | Extend API authorization mechanism
 Network       | [Link](https://docs.docker.com/engine/extend/plugins_network/) | Extend network management
 Volume        | [Link](https://docs.docker.com/engine/extend/plugins_volume/)  | Extend persistent storage
 IPAM          | [Link](https://github.com/docker/libnetwork/blob/master/docs/ipam.md) | Extend IP address management
 Graph (experimental) | [Link](https://github.com/docker/docker/blob/master/experimental/plugins_graphdriver.md) | Extend image and container fs storage

See the [understand Docker plugins documentation section](https://docs.docker.com/engine/extend/plugins/).

# Serving multiple drivers on a single socket
```go
import (
    "github.com/docker/go-plugins-helpers/sdk"
    "github.com/docker/go-plugins-helpers/ipam"
    "github.com/docker/go-plugins-helpers/network"
)
dIPAM := MyIPAMDriver{}
dNetwork := MyNetworkDriver{}
h := sdk.NewHandler()
ipam.RegisterDriver(dIPAM, h)
network.RegisterDriver(dNetwork,h)
h.ServeUnix("/var/run/docker/plugins/myplugin.sock", 1001)
```
