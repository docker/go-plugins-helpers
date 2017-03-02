package volume

import (
	"log"
	"net/http"

	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	// DefaultDockerRootDirectory is the default directory where volumes will be created.
	DefaultDockerRootDirectory = "/var/lib/docker-volumes"

	manifest         = `{"Implements": ["VolumeDriver"]}`
	createPath       = "/VolumeDriver.Create"
	getPath          = "/VolumeDriver.Get"
	listPath         = "/VolumeDriver.List"
	removePath       = "/VolumeDriver.Remove"
	hostVirtualPath  = "/VolumeDriver.Path"
	mountPath        = "/VolumeDriver.Mount"
	unmountPath      = "/VolumeDriver.Unmount"
	capabilitiesPath = "/VolumeDriver.Capabilities"
)

// CreateRequest is the structure that docker's requests are deserialized to.
type CreateRequest struct {
	Name    string
	Options map[string]string `json:"Opts,omitempty"`
}

// RemoveRequest structure for a volume remove request
type RemoveRequest struct {
	Name string
}

// MountRequest structure for a volume mount request
type MountRequest struct {
	Name string
	ID   string
}

// MountResponse structure for a volume mount response
type MountResponse struct {
	Mountpoint string
	Err        string
}

// UnmountRequest structure for a volume unmount request
type UnmountRequest struct {
	Name string
	ID   string
}

// PathRequest structure for a volume path request
type PathRequest struct {
	Name string
}

// PathResponse structure for a volume path response
type PathResponse struct {
	Mountpoint string
	Err        string
}

// GetRequest structure for a volume get request
type GetRequest struct {
	Name string
}

// GetResponse structure for a volume get response
type GetResponse struct {
	Err    string
	Volume *Volume
}

// ListResponse structure for a volume list response
type ListResponse struct {
	Err     string
	Volumes []*Volume
}

// CapabilitiesResponse structure for a volume capability response
type CapabilitiesResponse struct {
	Err          string
	Capabilities Capability
}

// Volume represents a volume object for use with `Get` and `List` requests
type Volume struct {
	Name       string
	Mountpoint string
	Status     map[string]interface{}
}

// Capability represents the list of capabilities a volume driver can return
type Capability struct {
	Scope string
}

// ErrorResponse is a formatted error message that docker can understand
type ErrorResponse struct {
	Err string
}

// NewErrorResponse creates an ErrorResponse with the provided message
func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Err: msg}
}

// Driver represent the interface a driver must fulfill.
type Driver interface {
	Create(*CreateRequest) error
	List() (*ListResponse, error)
	Get(*GetRequest) (*GetResponse, error)
	Remove(*RemoveRequest) error
	Path(*PathRequest) (*PathResponse, error)
	Mount(*MountRequest) (*MountResponse, error)
	Unmount(*UnmountRequest) error
	Capabilities() *CapabilitiesResponse
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
	h.HandleFunc(createPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers createPath")
		req := &CreateRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.Create(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, NewErrorResponse(""), "")
	})
	h.HandleFunc(removePath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers removePath")
		req := &RemoveRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.Remove(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, NewErrorResponse(""), "")
	})
	h.HandleFunc(mountPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers mountPath")
		req := &MountRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.Mount(req)
		if err != nil {
			res.Err = err.Error()
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(hostVirtualPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers hostVirtualPath")
		req := &PathRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.Path(req)
		if err != nil {
			res.Err = err.Error()
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(getPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers getPath")
		req := &GetRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.Get(req)
		if err != nil {
			res.Err = err.Error()
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(unmountPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers unmountPath")
		req := &UnmountRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.Unmount(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, NewErrorResponse(""), "")
	})
	h.HandleFunc(listPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers listPath")
		res, err := h.driver.List()
		if err != nil {
			res.Err = err.Error()
		}
		sdk.EncodeResponse(w, res, "")
	})

	h.HandleFunc(capabilitiesPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers capabilitiesPath")
		sdk.EncodeResponse(w, h.driver.Capabilities(), "")
	})
}
