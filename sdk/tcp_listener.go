package sdk

import (
	"crypto/tls"
	"net"
	"runtime"

	"github.com/docker/go-connections/sockets"
)

// NewTCPListener constructs a net.Listener to use for serving requests at the given TCP address.
// It also writes the spec file in the right directory for docker to read.
//
// Due to constrains for running Docker in Docker on Windows, data-root directory
// of docker daemon must be provided. To get default directory, use
// WindowsDefaultDaemonRootDir() function. On Unix, this parameter is ignored.
func NewTCPListener(address, pluginName, daemonDir string, tlsConfig *tls.Config) (net.Listener, string, error) {
	listener, err := sockets.NewTCPSocket(address, tlsConfig)
	if err != nil {
		return nil, "", err
	}

	addr := listener.Addr().String()

	var specDir string
	if runtime.GOOS == "windows" {
		specDir, err = createPluginSpecDirWindows(pluginName, addr, daemonDir)
	} else {
		specDir, err = createPluginSpecDirUnix(pluginName, addr)
	}
	if err != nil {
		return nil, "", err
	}

	specFile, err := writeSpecFile(pluginName, addr, specDir, protoTCP)
	if err != nil {
		return nil, "", err
	}
	return listener, specFile, nil
}
