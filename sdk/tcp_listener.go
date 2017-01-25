package sdk

import (
	"crypto/tls"
	"net"

	"github.com/docker/go-connections/sockets"
)

func newTCPListener(address string, pluginName string, tlsConfig *tls.Config) (net.Listener, string, error) {
	listener, err := sockets.NewTCPSocket(address, tlsConfig)
	if err != nil {
		return nil, "", err
	}
	spec, err := writeSpec(pluginName, listener.Addr().String(), protoTCP)
	if err != nil {
		return nil, "", err
	}
	return listener, spec, nil
}
