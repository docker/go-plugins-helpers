// +build !windows

package sdk

import (
	"errors"
	"net"
)

var (
	errOnlySupportedOnWindows = errors.New("named pipe creation is only supported on Windows")
)

func newWindowsListener(address, pluginName string, pipeConfig *WindowsPipeConfig) (net.Listener, string, error) {
	return nil, "", errOnlySupportedOnWindows
}

func windowsCreateDirectory(name string) error {
	return nil
}
