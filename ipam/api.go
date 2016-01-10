package ipam

import (
	"fmt"
	"net/http"

	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	manifest             = `{"Implements": ["IpamDriver"]}`
	defaultRequiresMac   = `{"RequiresMACAddress": true}`
	defaultAddressSpaces = `{"LocalDefaultAddressSpace": "defaultLocal", "GlobalDefaultAddressSpace": "defaultGlobal"}`

	capabilitiesPath         = "/IpamDriver.GetCapabilities"
	defaultAddressSpacesPath = "/IpamDriver.GetDefaultAddressSpaces"

	requestPoolPath = "/IpamDriver.RequestPool"
)

// Driver represent the interface a driver must fulfill.
type Driver interface {
	RequestPool(*RequestPoolRequest) (*RequestPoolResponse, error)
}

// RequestPoolRequest is sent by the Daemon requesting an address pool
type RequestPoolRequest struct {
	AddressSpace	string
	Pool			string
	SubPool			string
	Options			map[string]interface{}
	V6				bool
}

// RequestPoolResponse is sent in response to RequestPoolRequest
type RequestPoolResponse {
	PoolID	string
	Pool	string
	Data	map[string]string
}

// ErrorResponse is a formatted error message that libnetwork can understand
type ErrorResponse struct {
	Err string
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
		fmt.Fprintln(w, defaultRequiresMac)
	})

	h.HandleFunc(defaultAddressSpacesPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, defaultAddressSpaces)
	})

	h.HandleFunc(requestPoolPath, func(w http.ResponseWriter, r *http.Request) {
		req := &RequestPoolRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.RequestPool(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
		}
		sdk.EncodeResponse(w, res, "")
	})
}
