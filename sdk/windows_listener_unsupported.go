// +build !windows

package sdk

import (
	"errors"
	"net"
)

var (
	errOnlySupportedOnWindows = errors.New("named pipe creation is only supported on Windows")
)

func newWindowsListener(address, pluginName string, pipeConfig *WindowsPipeConfig) func() (net.Listener, string, string, error) {
	return func() (net.Listener, string, string, error) {
		return nil, "", "", errOnlySupportedOnWindows
	}
}
