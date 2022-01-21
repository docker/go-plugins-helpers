// +build !windows

package sdk

import (
	"errors"
	"net"
)

var (
	errOnlySupportedOnWindows = errors.New("named pipe creation is only supported on Windows")
)

// NewWindowsListener constructs a net.Listener to use for serving requests at the given Windows named pipe.
// It also creates the spec file in the right directory for docker to read.
//
// Due to constrains for running Docker in Docker on Windows, the data-root directory
// of docker daemon must be provided. To get default directory, use
// WindowsDefaultDaemonRootDir() function. On Unix, this parameter is ignored.
func NewWindowsListener(address, pluginName, daemonRoot string, pipeConfig *WindowsPipeConfig) (net.Listener, string, error) {
	return nil, "", errOnlySupportedOnWindows
}

func windowsCreateDirectoryWithACL(name string) error {
	return nil
}
