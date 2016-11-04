package sdk

import (
	"crypto/tls"
	"net"

	"github.com/docker/go-connections/sockets"
)

func newTCPListener(address string, pluginName string, tlsConfig *tls.Config) func() (net.Listener, string, string, error) {
	return func() (net.Listener, string, string, error) {
		listener, err := sockets.NewTCPSocket(address, tlsConfig)
		if err != nil {
			return nil, "", "", err
		}
		spec, err := writeSpec(pluginName, listener.Addr().String(), protoTCP)
		if err != nil {
			return nil, "", "", err
		}
		return listener, address, spec, nil
	}
}
