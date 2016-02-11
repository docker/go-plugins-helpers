package graph

import (
	"io"
	"net/http"

	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	manifest        = `{"Implements": ["GraphDriver"]}`
	initPath        = "/GraphDriver.Init"
	createPath      = "/GraphDriver.Create"
	removePath      = "/GraphDriver.Remove"
	getPath         = "/GraphDriver.Get"
	putPath         = "/GraphDriver.Put"
	existsPath      = "/GraphDriver.Exists"
	statusPath      = "/GraphDriver.Status"
	getMetadataPath = "/GraphDriver.GetMetdata"
	cleanupPath     = "/GraphDriver.Cleanup"
	diffPath        = "/GraphDriver.Diff"
	changesPath     = "/GraphDriver.Changes"
	applyDiffPath   = "/GraphDriver.ApplyDiff"
	diffSizePath    = "/GraphDriver.DiffSize"
)

// InitRequest initializes the plugin
type InitRequest struct {
	Home string
	Opts []string
}

// IDParentRequest is a request with the ID and Parent Fields
type IDParentRequest struct {
	ID     string
	Parent string
}

// OperationRequest is used for Put, Exists operations
type OperationRequest struct {
	ID string
}

// GetRequest is used to request a Get operation
type GetRequest struct {
	ID         string
	MountLabel string
}

// GetResponse is the response of a Get operation
type GetResponse struct {
	Dir string
}

// ExistsResponse is used to indicate if the requested item exists
type ExistsResponse struct {
	Exists bool
}

// MetadataResponse contains the metadata of requested item
type MetadataResponse struct {
	Metadata map[string]interface{}
}

// StatusResponse is a status
type StatusResponse struct {
	Status [][]string
}

// ChangeKind represents the type of change mage
type ChangeKind int

const (
	// Modified is a ChangeKind used when an item has been modified
	Modified ChangeKind = iota
	// Added is a ChangeKind used when an item has been added
	Added
	// Deleted is a ChangeKind used when an item has been deleted
	Deleted
)

// Change represents a single Change made to a directory
type Change struct {
	Path string
	Kind ChangeKind
}

// ChangesResponse contains a list of Change that were made to a directory
type ChangesResponse struct {
	Changes []Change
}

// SizeResponse returns the Size of an item
type SizeResponse struct {
	Size int
}

// ErrorResponse is a formatted error message that libnetwork can understand
type ErrorResponse struct {
	Err string
}

// NewErrorResponse creates an ErrorResponse with the provided message
func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Err: msg}
}

// Driver represents the interface a driver must fulfill
type Driver interface {
	Init(*InitRequest) error
	Create(*IDParentRequest) error
	Remove(*OperationRequest) error
	Get(*GetRequest) (*GetResponse, error)
	Put(*OperationRequest) error
	Exists(*OperationRequest) (*ExistsResponse, error)
	Status() *StatusResponse
	GetMetadata(*OperationRequest) (*MetadataResponse, error)
	Cleanup() error
	Diff(*IDParentRequest) io.ReadCloser
	Changes(*IDParentRequest) (*ChangesResponse, error)
	ApplyDiff(io.Reader) (*SizeResponse, error)
	DiffSize(*IDParentRequest) (*SizeResponse, error)
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
	h.HandleFunc(initPath, func(w http.ResponseWriter, r *http.Request) {
		req := &InitRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.Init(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(createPath, func(w http.ResponseWriter, r *http.Request) {
		req := &IDParentRequest{}
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
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(removePath, func(w http.ResponseWriter, r *http.Request) {
		req := &OperationRequest{}
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
		sdk.EncodeResponse(w, make(map[string]string), "")

	})
	h.HandleFunc(getPath, func(w http.ResponseWriter, r *http.Request) {
		req := &GetRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.Get(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(putPath, func(w http.ResponseWriter, r *http.Request) {
		req := &OperationRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.Put(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")

	})
	h.HandleFunc(existsPath, func(w http.ResponseWriter, r *http.Request) {
		req := &OperationRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.Exists(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")

	})
	h.HandleFunc(statusPath, func(w http.ResponseWriter, r *http.Request) {
		res := h.driver.Status()
		sdk.EncodeResponse(w, res, "")

	})
	h.HandleFunc(getMetadataPath, func(w http.ResponseWriter, r *http.Request) {
		req := &OperationRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.GetMetadata(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")

	})
	h.HandleFunc(cleanupPath, func(w http.ResponseWriter, r *http.Request) {
		err := h.driver.Cleanup()
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")

	})
	h.HandleFunc(diffPath, func(w http.ResponseWriter, r *http.Request) {
		req := &IDParentRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res := h.driver.Diff(req)
		sdk.StreamResponse(w, res)

	})
	h.HandleFunc(changesPath, func(w http.ResponseWriter, r *http.Request) {
		req := &IDParentRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.Changes(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")

	})
	h.HandleFunc(applyDiffPath, func(w http.ResponseWriter, r *http.Request) {
		res, err := h.driver.ApplyDiff(r.Body)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")

	})
	h.HandleFunc(diffSizePath, func(w http.ResponseWriter, r *http.Request) {
		req := &IDParentRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.DiffSize(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, res, "")

	})
}
