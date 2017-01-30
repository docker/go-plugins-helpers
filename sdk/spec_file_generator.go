package sdk

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type protocol string

const (
	protoTCP       protocol = "tcp"
	protoNamedPipe protocol = "npipe"
)

// PluginSpecDir returns plugin spec dir in relation to daemon root directory.
func PluginSpecDir(daemonRoot string) string {
	return ([]string{filepath.Join(daemonRoot, "plugins")})[0]
}

func createPluginSpecDirWindows(name, address, daemonRoot string) (string, error) {
	if daemonRoot == "" {
		daemonRoot = filepath.Join(os.Getenv("programdata"), "docker")
	}
	pluginSpecDir := PluginSpecDir(daemonRoot)
	if err := windowsCreateDirectoryWithACL(pluginSpecDir); err != nil {
		return "", err
	}
	return pluginSpecDir, nil
}

func createPluginSpecDirUnix(name, address string) (string, error) {
	pluginSpecDir := PluginSpecDir("/etc/docker")
	if err := os.MkdirAll(pluginSpecDir, 0755); err != nil {
		return "", err
	}
	return pluginSpecDir, nil
}

func writeSpecFile(name, address, pluginSpecDir string, proto protocol) (string, error) {
	specFileDir := filepath.Join(pluginSpecDir, name+".spec")

	url := string(proto) + "://" + address
	if err := ioutil.WriteFile(specFileDir, []byte(url), 0644); err != nil {
		return "", err
	}

	return specFileDir, nil
}
