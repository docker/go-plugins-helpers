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

An abridged example of serving identically named handlers over the same socket.

```go

import (
    "github.com/docker/go-plugins-helpers/sdk"
    "github.com/docker/go-plugins-helpers/ipam"
	"github.com/docker/go-plugins-helpers/network"
)

handler := sdk.NewHandler()

driver := &NetworkDriver{}
ipamDriver := &IpamDriver{}

network.InitMux(handler, driver)
ipam.InitMux(handler, ipamDriver)

handler.ServeUnix("/run/docker/plugins/my_combined_driver.sock", 0)

```
