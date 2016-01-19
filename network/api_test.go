package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/go-plugins-helpers/sdk"
)

type TestDriver struct {
	Driver
}

func (t *TestDriver) GetCapabilities() (*CapabilitiesResponse, error) {
	return &CapabilitiesResponse{Scope: LocalScope}, nil
}

func (t *TestDriver) CreateNetwork(r *CreateNetworkRequest) error {
	return nil
}

func (t *TestDriver) DeleteNetwork(r *DeleteNetworkRequest) error {
	return nil
}

func (t *TestDriver) CreateEndpoint(r *CreateEndpointRequest) (*CreateEndpointResponse, error) {
	return &CreateEndpointResponse{}, nil
}

func (t *TestDriver) DeleteEndpoint(r *DeleteEndpointRequest) error {
	return nil
}

func (t *TestDriver) Join(r *JoinRequest) (*JoinResponse, error) {
	return &JoinResponse{}, nil
}

func (t *TestDriver) Leave(r *LeaveRequest) error {
	return nil
}

type ErrDriver struct {
	Driver
}

func (e *ErrDriver) GetCapabilities() (*CapabilitiesResponse, error) {
	return nil, fmt.Errorf("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) CreateNetwork(r *CreateNetworkRequest) error {
	return fmt.Errorf("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) DeleteNetwork(r *DeleteNetworkRequest) error {
	return errors.New("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) CreateEndpoint(r *CreateEndpointRequest) (*CreateEndpointResponse, error) {
	return nil, errors.New("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) DeleteEndpoint(r *DeleteEndpointRequest) error {
	return errors.New("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) Join(r *JoinRequest) (*JoinResponse, error) {
	return nil, errors.New("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) Leave(r *LeaveRequest) error {
	return errors.New("I CAN HAZ ERRORZ")
}

func TestMain(m *testing.M) {
	d := &TestDriver{}
	h1 := NewHandler(d)
	go h1.ServeTCP("test", ":328234")

	e := &ErrDriver{}
	h2 := NewHandler(e)
	go h2.ServeTCP("err", ":328567")

	m.Run()
}

func TestActivate(t *testing.T) {
	response, err := http.Get("http://localhost:328234/Plugin.Activate")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if string(body) != manifest+"\n" {
		t.Fatalf("Expected %s, got %s\n", manifest+"\n", string(body))
	}
}

func TestCapabilitiesExchange(t *testing.T) {
	response, err := http.Get("http://localhost:328234/NetworkDriver.GetCapabilities")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if string(body) != defaultScope+"\n" {
		t.Fatalf("Expected %s, got %s\n", defaultScope+"\n", string(body))
	}
}

func TestCreateNetworkSuccess(t *testing.T) {
	request := `{"NetworkID":"d76cfa51738e8a12c5eca71ee69e9d65010a4b48eaad74adab439be7e61b9aaf","Options":{"com.docker.network.generic":{}},"IPv4Data":[{"AddressSpace":"","Gateway":"172.18.0.1/16","Pool":"172.18.0.0/16"}],"IPv6Data":[]}`

	response, err := http.Post("http://localhost:328234/NetworkDriver.CreateNetwork",
		sdk.DefaultContentTypeV1_1,
		strings.NewReader(request),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d\n", response.StatusCode)
	}
	if string(body) != "{}\n" {
		t.Fatalf("Expected %s, got %s\n", "{}\n", string(body))
	}
}

func TestCreateNetworkError(t *testing.T) {
	request := `{"NetworkID":"d76cfa51738e8a12c5eca71ee69e9d65010a4b48eaad74adab439be7e61b9aaf","Options":{"com.docker.network.generic":    {}},"IPv4Data":[{"AddressSpace":"","Gateway":"172.18.0.1/16","Pool":"172.18.0.0/16"}],"IPv6Data":[]}`
	response, err := http.Post("http://localhost:328567/NetworkDriver.CreateNetwork",
		sdk.DefaultContentTypeV1_1,
		strings.NewReader(request))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected 500, got %d\n", response.StatusCode)
	}
	if string(body) != "{\"Err\":\"I CAN HAZ ERRORZ\"}\n" {
		t.Fatalf("Expected %s, got %s\n", "{\"Err\":\"I CAN HAZ ERRORZ\"}\n", string(body))
	}
}
