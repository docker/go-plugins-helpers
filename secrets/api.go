package secrets

import (
	"github.com/docker/go-plugins-helpers/sdk"
	"net/http"
)

const (
	manifest = `{"Implements": ["secretprovider"]}`
	getPath  = "/SecretProvider.GetSecret"
)

// Request is the plugin secret request
type Request struct {
	SecretName          string            `json:",omitempty"` // SecretName is the name of the secret to request from the plugin
	ServiceHostname     string            `json:",omitempty"` // ServiceHostname is the hostname of the service, can be used for x509 certificate
	ServiceName         string            `json:",omitempty"` // ServiceName is the name of the service that requested the secret
	ServiceLabels       map[string]string `json:",omitempty"` // ServiceLabels capture environment names and other metadata
	ServiceEndpointSpec *EndpointSpec     `json:",omitempty"` // ServiceEndpointSpec holds the specification for endpoints
}

// Response contains the plugin secret value
type Response struct {
	Value []byte `json:",omitempty"` // Value is the value of the secret
	Err   string `json:",omitempty"` // Err is the error response of the plugin
}

// EndpointSpec represents the spec of an endpoint.
type EndpointSpec struct {
	Mode  int32        `json:",omitempty"`
	Ports []PortConfig `json:",omitempty"`
}

// PortConfig represents the config of a port.
type PortConfig struct {
	Name     string `json:",omitempty"`
	Protocol int32  `json:",omitempty"`
	// TargetPort is the port inside the container
	TargetPort uint32 `json:",omitempty"`
	// PublishedPort is the port on the swarm hosts
	PublishedPort uint32 `json:",omitempty"`
	// PublishMode is the mode in which port is published
	PublishMode int32 `json:",omitempty"`
}

// Driver represent the interface a driver must fulfill.
type Driver interface {
	// Get gets a secret from a remote secret store
	Get(Request) Response
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
	h.HandleFunc(getPath, func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := sdk.DecodeRequest(w, r, &req); err != nil {
			return
		}
		res := h.driver.Get(req)
		sdk.EncodeResponse(w, res, res.Err)
	})
}
