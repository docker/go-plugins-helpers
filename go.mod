module github.com/docker/go-plugins-helpers

go 1.14

replace (
	github.com/codegangsta/cli => github.com/urfave/cli v1.22.4
	github.com/docker/docker v1.13.1 => github.com/docker/engine v0.0.0-20200618181300-9dc6525e6118
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc10
)

require (
	github.com/Microsoft/go-winio v0.4.15-0.20190919025122-fc70bd9a86b5
	github.com/Microsoft/hcsshim v0.8.9 // indirect
	github.com/containerd/containerd v1.4.0 // indirect
	github.com/containerd/continuity v0.0.0-20200710164510-efbc4488d8fe // indirect
	github.com/coreos/go-systemd/v22 v22.1.0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/vbatts/tar-split v0.11.1 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	google.golang.org/grpc v1.31.1 // indirect
)
