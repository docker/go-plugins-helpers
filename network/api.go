package network

import (
	"log"
	"net/http"

	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	manifest     = `{"Implements": ["NetworkDriver"]}`
	defaultScope = `{"Scope":"local"}`

	// LocalScope is the correct scope response for a local scope driver
	LocalScope = `local`
	// GlobalScope is the correct scope response for a global scope driver
	GlobalScope = `global`

	capabilitiesPath   = "/NetworkDriver.GetCapabilities"
	createNetworkPath  = "/NetworkDriver.CreateNetwork"
	deleteNetworkPath  = "/NetworkDriver.DeleteNetwork"
	createEndpointPath = "/NetworkDriver.CreateEndpoint"
	endpointInfoPath   = "/NetworkDriver.EndpointOperInfo"
	deleteEndpointPath = "/NetworkDriver.DeleteEndpoint"
	joinPath           = "/NetworkDriver.Join"
	leavePath          = "/NetworkDriver.Leave"
	discoverNewPath    = "/NetworkDriver.DiscoverNew"
	discoverDeletePath = "/NetworkDriver.DiscoverDelete"
)

// Driver represent the interface a driver must fulfill.
type Driver interface {
	GetCapabilities() (*CapabilitiesResponse, error)
	CreateNetwork(*CreateNetworkRequest) error
	DeleteNetwork(*DeleteNetworkRequest) error
	CreateEndpoint(*CreateEndpointRequest) (*CreateEndpointResponse, error)
	DeleteEndpoint(*DeleteEndpointRequest) error
	EndpointInfo(*InfoRequest) (*InfoResponse, error)
	Join(*JoinRequest) (*JoinResponse, error)
	Leave(*LeaveRequest) error
	DiscoverNew(*DiscoveryNotification) error
	DiscoverDelete(*DiscoveryNotification) error
}

// CapabilitiesResponse returns whether or not this network is global or local
type CapabilitiesResponse struct {
	Scope string
}

// CreateNetworkRequest is sent by the daemon when a network needs to be created
type CreateNetworkRequest struct {
	NetworkID string
	Options   map[string]interface{}
	IPv4Data  []*IPAMData
	IPv6Data  []*IPAMData
}

// IPAMData contains IPv4 or IPv6 addressing information
type IPAMData struct {
	AddressSpace string
	Pool         string
	Gateway      string
	AuxAddresses map[string]interface{}
}

// DeleteNetworkRequest is sent by the daemon when a network needs to be removed
type DeleteNetworkRequest struct {
	NetworkID string
}

// CreateEndpointRequest is sent by the daemon when an endpoint should be created
type CreateEndpointRequest struct {
	NetworkID  string
	EndpointID string
	Interface  *EndpointInterface
	Options    map[string]interface{}
}

type CreateEndpointResponse struct {
	Interface *EndpointInterface
}

// EndpointInterface contains endpoint interface information
type EndpointInterface struct {
	Address     string
	AddressIPv6 string
	MacAddress  string
}

// DeleteEndpointRequest is sent by the daemon when an endpoint needs to be removed
type DeleteEndpointRequest struct {
	NetworkID  string
	EndpointID string
}

// InterfaceName consists of the name of the interface in the global netns and
// the desired prefix to be appended to the interface inside the container netns
type InterfaceName struct {
	SrcName   string
	DstPrefix string
}

// InfoRequest is send by the daemon when querying endpoint information
type InfoRequest struct {
	NetworkID string
	EnpointID string
}

// InfoResponse is endpoint information sent in response to an InfoRequest
type InfoResponse struct {
	Value map[string]string
}

// JoinRequest is sent by the Daemon when an endpoint needs be joined to a network
type JoinRequest struct {
	NetworkID  string
	EndpointID string
	SandboxKey string
	Options    map[string]interface{}
}

// StaticRoute contains static route information
type StaticRoute struct {
	Destination string
	RouteType   int
	NextHop     string
}

// JoinResponse is sent in response to a JoinRequest
type JoinResponse struct {
	Gateway       string
	InterfaceName InterfaceName
	StaticRoutes  []*StaticRoute
}

// LeaveRequest is send by the daemon when a endpoint is leaving a network
type LeaveRequest struct {
	NetworkID  string
	EndpointID string
	Options    map[string]interface{}
}

// ErrorResponse is a formatted error message that libnetwork can understand
type ErrorResponse struct {
	Err string
}

type DiscoveryNotification struct {
	DiscoveryType int
	DiscoveryData interface{}
}

// NewErrorResponse creates an ErrorResponse with the provided message
func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Err: msg}
}

// Handler forwards requests and responses between the docker daemon and the plugin.
type Handler struct {
	driver Driver
	sdk.Handler
}

// NewHandler initializes the request handler with a driver implementation.
func NewHandler(driver Driver) *Handler {
	h := &Handler{driver, sdk.NewHandler(manifest)}
	h.initMux()
	return h
}

func (h *Handler) initMux() {
	h.HandleFunc(capabilitiesPath, func(w http.ResponseWriter, r *http.Request) {
		res, err := h.driver.GetCapabilities()
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		if res == nil {
			msg := "Network driver must implement GetCapabilities"
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(createNetworkPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers createnetwork")
		req := &CreateNetworkRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.CreateNetwork(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(deleteNetworkPath, func(w http.ResponseWriter, r *http.Request) {
		req := &DeleteNetworkRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.DeleteNetwork(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(createEndpointPath, func(w http.ResponseWriter, r *http.Request) {
		req := &CreateEndpointRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.CreateEndpoint(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(deleteEndpointPath, func(w http.ResponseWriter, r *http.Request) {
		req := &DeleteEndpointRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.DeleteEndpoint(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(endpointInfoPath, func(w http.ResponseWriter, r *http.Request) {
		req := &InfoRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.EndpointInfo(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(joinPath, func(w http.ResponseWriter, r *http.Request) {
		req := &JoinRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.Join(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(leavePath, func(w http.ResponseWriter, r *http.Request) {
		req := &LeaveRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.Leave(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(discoverNewPath, func(w http.ResponseWriter, r *http.Request) {
		req := &DiscoveryNotification{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.DiscoverNew(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(discoverDeletePath, func(w http.ResponseWriter, r *http.Request) {
		req := &DiscoveryNotification{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.DiscoverDelete(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
}
