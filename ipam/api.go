package ipam

import (
	"net/http"

	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	driverName = "IpamDriver"

	capabilitiesPath   = "/IpamDriver.GetCapabilities"
	addressSpacesPath  = "/IpamDriver.GetDefaultAddressSpaces"
	requestPoolPath    = "/IpamDriver.RequestPool"
	releasePoolPath    = "/IpamDriver.ReleasePool"
	requestAddressPath = "/IpamDriver.RequestAddress"
	releaseAddressPath = "/IpamDriver.ReleaseAddress"
)

// Ipam represent the interface a driver must fulfill.
type Ipam interface {
	GetCapabilities() (*CapabilitiesResponse, error)
	GetDefaultAddressSpaces() (*AddressSpacesResponse, error)
	RequestPool(*RequestPoolRequest) (*RequestPoolResponse, error)
	ReleasePool(*ReleasePoolRequest) error
	RequestAddress(*RequestAddressRequest) (*RequestAddressResponse, error)
	ReleaseAddress(*ReleaseAddressRequest) error
}

// CapabilitiesResponse returns whether or not this IPAM required pre-made MAC
type CapabilitiesResponse struct {
	RequiresMACAddress bool
}

// AddressSpacesResponse returns the default local and global address space names for this IPAM
type AddressSpacesResponse struct {
	LocalDefaultAddressSpace  string
	GlobalDefaultAddressSpace string
}

// RequestPoolRequest is sent by the daemon when a pool needs to be created
type RequestPoolRequest struct {
	AddressSpace string
	Pool         string
	SubPool      string
	Options      map[string]string
	V6           bool
}

// RequestPoolResponse returns a registered address pool with the IPAM driver
type RequestPoolResponse struct {
	PoolID string
	Pool   string
	Data   map[string]string
}

// ReleasePoolRequest is sent when releasing a previously registered address pool
type ReleasePoolRequest struct {
	PoolID string
}

// RequestAddressRequest is sent when requesting an address from IPAM
type RequestAddressRequest struct {
	PoolID  string
	Address string
	Options map[string]string
}

// RequestAddressResponse is formed with allocated address by IPAM
type RequestAddressResponse struct {
	Address string
	Data    map[string]string
}

// ReleaseAddressRequest is sent in order to release an address from the pool
type ReleaseAddressRequest struct {
	PoolID  string
	Address string
}

// ErrorResponse is a formatted error message that libnetwork can understand
type ErrorResponse struct {
	Err string
}

// NewErrorResponse creates an ErrorResponse with the provided message
func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Err: msg}
}

// NewHandler initializes the request handler with a driver implementation.
func NewHandler(ipam Ipam) *sdk.Handler {
	h := sdk.NewHandler()
	RegisterDriver(ipam, h)
	return h
}

// RegisterDriver registers the plugin to the SDK handler.
func RegisterDriver(ipam Ipam, h *sdk.Handler) {
	h.RegisterDriver(driverName, func(h *sdk.Handler) {
		initMux(ipam, h)
	})
}

func initMux(ipam Ipam, h *sdk.Handler) {
	h.HandleFunc(capabilitiesPath, func(w http.ResponseWriter, r *http.Request) {
		res, err := ipam.GetCapabilities()
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(addressSpacesPath, func(w http.ResponseWriter, r *http.Request) {
		res, err := ipam.GetDefaultAddressSpaces()
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(requestPoolPath, func(w http.ResponseWriter, r *http.Request) {
		req := &RequestPoolRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := ipam.RequestPool(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(releasePoolPath, func(w http.ResponseWriter, r *http.Request) {
		req := &ReleasePoolRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = ipam.ReleasePool(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(requestAddressPath, func(w http.ResponseWriter, r *http.Request) {
		req := &RequestAddressRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := ipam.RequestAddress(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(releaseAddressPath, func(w http.ResponseWriter, r *http.Request) {
		req := &ReleaseAddressRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = ipam.ReleaseAddress(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
}
