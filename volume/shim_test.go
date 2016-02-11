package volume

import (
	"net/http"
	"testing"

	"github.com/docker/docker/volume"
	"github.com/docker/go-connections/sockets"
)

type testVolumeDriver struct{}
type testVolume struct{}

func (testVolume) Name() string           { return "" }
func (testVolume) Path() string           { return "" }
func (testVolume) Mount() (string, error) { return "", nil }
func (testVolume) Unmount() error         { return nil }
func (testVolume) DriverName() string     { return "" }

func (testVolumeDriver) Name() string                                            { return "" }
func (testVolumeDriver) Create(string, map[string]string) (volume.Volume, error) { return nil, nil }
func (testVolumeDriver) Remove(volume.Volume) error                              { return nil }
func (testVolumeDriver) List() ([]volume.Volume, error)                          { return nil, nil }
func (testVolumeDriver) Get(name string) (volume.Volume, error)                  { return nil, nil }

func TestVolumeDriver(t *testing.T) {
	h := NewHandlerFromVolumeDriver(testVolumeDriver{})
	l := sockets.NewInmemSocket("test", 0)
	go h.Serve(l)
	defer l.Close()

	client := &http.Client{Transport: &http.Transport{
		Dial: l.Dial,
	}}

	resp, err := pluginRequest(client, createPath, Request{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err != "" {
		t.Fatalf("error while creating volume: %v", err)
	}
}
