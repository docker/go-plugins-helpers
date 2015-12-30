package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/go-plugins-sdk/sdk"

	. "gopkg.in/check.v1"
)

type TestDriver struct {
	Driver
}

func (t *TestDriver) CreateNetwork(r *CreateNetworkRequest) error {
	return nil
}

func (t *TestDriver) DeleteNetwork(r *DeleteNetworkRequest) error {
	return nil
}

func (t *TestDriver) CreateEndpoint(r *CreateEndpointRequest) error {
	return nil
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

func (e *ErrDriver) CreateNetwork(r *CreateNetworkRequest) error {
	return fmt.Errorf("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) DeleteNetwork(r *DeleteNetworkRequest) error {
	return errors.New("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) CreateEndpoint(r *CreateEndpointRequest) error {
	return errors.New("I CAN HAZ ERRORZ")
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

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	h1 *Handler
	h2 *Handler
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
	d := &TestDriver{}
	s.h1 = NewHandler(d)
	go s.h1.ServeTCP("test", ":8080")

	e := &ErrDriver{}
	s.h2 = NewHandler(e)
	go s.h2.ServeTCP("err", ":8888")
}

func (s *MySuite) TestActivate(c *C) {
	response, err := http.Get("http://localhost:8080/Plugin.Activate")
	if err != nil {
		c.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	c.Assert(string(body), Equals, manifest+"\n")

}

func (s *MySuite) TestCapabilitiesExchange(c *C) {
	response, err := http.Get("http://localhost:8080/NetworkDriver.GetCapabilities")
	if err != nil {
		c.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	c.Assert(string(body), Equals, defaultScope+"\n")

}

func (s *MySuite) TestCreateNetworkSuccess(c *C) {
	request := `{"NetworkID":"d76cfa51738e8a12c5eca71ee69e9d65010a4b48eaad74adab439be7e61b9aaf","Options":{"com.docker.network.generic":{}},"IPv4Data":[{"AddressSpace":"","Gateway":"172.18.0.1/16","Pool":"172.18.0.0/16"}],"IPv6Data":[]}`

	response, err := http.Post("http://localhost:8080/NetworkDriver.CreateNetwork",
		sdk.DefaultContentTypeV1_1,
		strings.NewReader(request),
	)
	if err != nil {
		c.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	c.Assert(response.StatusCode, Equals, http.StatusOK)
	c.Assert(string(body), Equals, "{}\n")

}

func (s *MySuite) TestCreateNetworkError(c *C) {
	request := `{"NetworkID":"d76cfa51738e8a12c5eca71ee69e9d65010a4b48eaad74adab439be7e61b9aaf","Options":{"com.docker.network.generic":    {}},"IPv4Data":[{"AddressSpace":"","Gateway":"172.18.0.1/16","Pool":"172.18.0.0/16"}],"IPv6Data":[]}`
	response, err := http.Post("http://localhost:8888/NetworkDriver.CreateNetwork",
		sdk.DefaultContentTypeV1_1,
		strings.NewReader(request))
	if err != nil {
		c.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	c.Assert(response.StatusCode, Equals, http.StatusInternalServerError)
	c.Assert(string(body), Equals, `{"Err":"I CAN HAZ ERRORZ"}`+"\n")

}
