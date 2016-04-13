package graphdriver

// See https://github.com/docker/docker/blob/master/experimental/plugins_graphdriver.md

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	// DefaultDockerRootDirectory is the default directory where graph drivers will be created.
	DefaultDockerRootDirectory = "/var/lib/docker/graph"

	manifest        = `{"Implements": ["GraphDriver"]}`
	initPath        = "/GraphDriver.Init"
	createPath      = "/GraphDriver.Create"
	createRWPath    = "/GraphDriver.CreateReadWrite"
	removePath      = "/GraphDriver.Remove"
	getPath         = "/GraphDriver.Get"
	putPath         = "/GraphDriver.Put"
	existsPath      = "/GraphDriver.Exists"
	statusPath      = "/GraphDriver.Status"
	getMetadataPath = "/GraphDriver.GetMetadata"
	cleanupPath     = "/GraphDriver.Cleanup"
	diffPath        = "/GraphDriver.Diff"
	changesPath     = "/GraphDriver.Changes"
	applyDiffPath   = "/GraphDriver.ApplyDiff"
	diffSizePath    = "/GraphDriver.DiffSize"
)

// Init

// InitRequest is the structure that docker's init requests are deserialized to.
type InitRequest struct {
	Home    string
	Options []string `json:"Opts"`
}

// InitResponse is the strucutre that docker's init responses are serialized to.
type InitResponse struct {
	Err string
}

// Create

// CreateRequest is the structure that docker's create requests are deserialized to.
type CreateRequest struct {
	ID         string
	Parent     string
	MountLabel string
}

// CreateResponse is the strucutre that docker's create responses are serialized to.
type CreateResponse struct {
	Err string
}

// Remove

// RemoveRequest is the structure that docker's remove requests are deserialized to.
type RemoveRequest struct {
	ID string
}

// RemoveResponse is the strucutre that docker's remove responses are serialized to.
type RemoveResponse struct {
	Err string
}

// Get

// GetRequest is the structure that docker's get requests are deserialized to.
type GetRequest struct {
	ID         string
	MountLabel string
}

// GetResponse is the strucutre that docker's remove responses are serialized to.
type GetResponse struct {
	Dir string
	Err string
}

// Put

// PutRequest is the structure that docker's put requests are deserialized to.
type PutRequest struct {
	ID string
}

// PutResponse is the strucutre that docker's put responses are serialized to.
type PutResponse struct {
	Err string
}

// Exists

// ExistsRequest is the structure that docker's exists requests are deserialized to.
type ExistsRequest struct {
	ID string
}

// ExistsResponse is the structure that docker's exists responses are serialized to.
type ExistsResponse struct {
	Exists bool
}

// Status

// StatusRequest is the structure that docker's status requests are deserialized to.
type StatusRequest struct{}

// StatusResponse is the structure that docker's status responses are serialized to.
type StatusResponse struct {
	Status [][2]string
}

// GetMetadata

// GetMetadataRequest is the structure that docker's getMetadata requests are deserialized to.
type GetMetadataRequest struct {
	ID string
}

// GetMetadataResponse is the structure that docker's getMetadata responses are serialized to.
type GetMetadataResponse struct {
	Metadata map[string]string
	Err      string
}

// Cleanup

// CleanupRequest is the structure that docker's cleanup requests are deserialized to.
type CleanupRequest struct{}

// CleanupResponse is the structure that docker's cleanup responses are serialized to.
type CleanupResponse struct {
	Err string
}

// Diff

// DiffRequest is the structure that docker's diff requests are deserialized to.
type DiffRequest struct {
	ID     string
	Parent string
}

// DiffResponse is the structure that docker's diff responses are serialized to.
type DiffResponse struct {
	Stream io.ReadCloser // TAR STREAM
}

// Changes

// ChangesRequest is the structure that docker's changes requests are deserialized to.
type ChangesRequest struct {
	ID     string
	Parent string
}

// ChangesResponse is the structure that docker's changes responses are serialized to.
type ChangesResponse struct {
	Changes []Change
	Err     string
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

// Change is the structure that docker's individual changes are serialized to.
type Change struct {
	Path string
	Kind ChangeKind
}

// ApplyDiff

// ApplyDiffRequest is the structure that docker's applyDiff requests are deserialized to.
type ApplyDiffRequest struct {
	Stream io.Reader // TAR STREAM
	ID     string
	Parent string
}

// ApplyDiffResponse is the structure that docker's applyDiff responses are serialized to.
type ApplyDiffResponse struct {
	Size int64
	Err  string
}

// DiffSize

// DiffSizeRequest is the structure that docker's diffSize requests are deserialized to.
type DiffSizeRequest struct {
	ID     string
	Parent string
}

// DiffSizeResponse is the structure that docker's diffSize responses are serialized to.
type DiffSizeResponse struct {
	Size int64
	Err  string
}

// Driver represent the interface a driver must fulfill.
type Driver interface {
	Init(home string, options []string) error
	Create(id, parent string) error
	CreateReadWrite(id, parent string) error
	Remove(id string) error
	Get(id, mountLabel string) (string, error)
	Put(id string) error
	Exists(id string) bool
	Status() [][2]string
	GetMetadata(id string) (map[string]string, error)
	Cleanup() error
	Diff(id, parent string) io.ReadCloser
	Changes(id, parent string) ([]Change, error)
	ApplyDiff(id, parent string, archive io.Reader) (int64, error)
	DiffSize(id, parent string) (int64, error)
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
		req := InitRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		err = h.driver.Init(req.Home, req.Options)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &InitResponse{Err: msg}, msg)
	})
	h.HandleFunc(createPath, func(w http.ResponseWriter, r *http.Request) {
		req := CreateRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		err = h.driver.Create(req.ID, req.Parent)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &CreateResponse{Err: msg}, msg)
	})
	h.HandleFunc(createRWPath, func(w http.ResponseWriter, r *http.Request) {
		req := CreateRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		err = h.driver.CreateReadWrite(req.ID, req.Parent)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &CreateResponse{Err: msg}, msg)
	})
	h.HandleFunc(removePath, func(w http.ResponseWriter, r *http.Request) {
		req := RemoveRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		err = h.driver.Remove(req.ID)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &RemoveResponse{Err: msg}, msg)

	})
	h.HandleFunc(getPath, func(w http.ResponseWriter, r *http.Request) {
		req := GetRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		dir, err := h.driver.Get(req.ID, req.MountLabel)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &GetResponse{Err: msg, Dir: dir}, msg)
	})
	h.HandleFunc(putPath, func(w http.ResponseWriter, r *http.Request) {
		req := PutRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		err = h.driver.Put(req.ID)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &PutResponse{Err: msg}, msg)
	})
	h.HandleFunc(existsPath, func(w http.ResponseWriter, r *http.Request) {
		req := ExistsRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		exists := h.driver.Exists(req.ID)
		sdk.EncodeResponse(w, &ExistsResponse{Exists: exists}, "")
	})
	h.HandleFunc(statusPath, func(w http.ResponseWriter, r *http.Request) {
		status := h.driver.Status()
		sdk.EncodeResponse(w, &StatusResponse{Status: status}, "")
	})
	h.HandleFunc(getMetadataPath, func(w http.ResponseWriter, r *http.Request) {
		req := GetMetadataRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		metadata, err := h.driver.GetMetadata(req.ID)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &GetMetadataResponse{Err: msg, Metadata: metadata}, msg)
	})
	h.HandleFunc(cleanupPath, func(w http.ResponseWriter, r *http.Request) {
		err := h.driver.Cleanup()
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &CleanupResponse{Err: msg}, msg)
	})
	h.HandleFunc(diffPath, func(w http.ResponseWriter, r *http.Request) {
		req := DiffRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		stream := h.driver.Diff(req.ID, req.Parent)
		sdk.StreamResponse(w, stream)
	})
	h.HandleFunc(changesPath, func(w http.ResponseWriter, r *http.Request) {
		req := ChangesRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		changes, err := h.driver.Changes(req.ID, req.Parent)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &ChangesResponse{Err: msg, Changes: changes}, msg)
	})
	h.HandleFunc(applyDiffPath, func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("id")
		parent := r.Header.Get("parent")
		size, err := h.driver.ApplyDiff(id, parent, r.Body)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &ApplyDiffResponse{Err: msg, Size: size}, msg)
	})
	h.HandleFunc(diffSizePath, func(w http.ResponseWriter, r *http.Request) {
		req := DiffRequest{}
		err := sdk.DecodeRequest(w, r, &req)
		if err != nil {
			return
		}
		size, err := h.driver.DiffSize(req.ID, req.Parent)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		sdk.EncodeResponse(w, &DiffSizeResponse{Err: msg, Size: size}, msg)
	})
}

// CallInit is the raw call to the Graphdriver.Init method
func CallInit(url string, client *http.Client, req InitRequest) (*InitResponse, error) {
	method := initPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp InitResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallCreateReadWrite is the raw call to the Graphdriver.CreateReadWrite method
func CallCreateReadWrite(url string, client *http.Client, req CreateRequest) (*CreateResponse, error) {
	method := createRWPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp CreateResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallCreate is the raw call to the Graphdriver.Create method
func CallCreate(url string, client *http.Client, req CreateRequest) (*CreateResponse, error) {
	method := createPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp CreateResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallRemove is the raw call to the Graphdriver.Remove method
func CallRemove(url string, client *http.Client, req RemoveRequest) (*RemoveResponse, error) {
	method := removePath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp RemoveResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallGet is the raw call to the Graphdriver.Get method
func CallGet(url string, client *http.Client, req GetRequest) (*GetResponse, error) {
	method := getPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp GetResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallPut is the raw call to the Graphdriver.Put method
func CallPut(url string, client *http.Client, req PutRequest) (*PutResponse, error) {
	method := putPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp PutResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallExists is the raw call to the Graphdriver.Exists method
func CallExists(url string, client *http.Client, req ExistsRequest) (*ExistsResponse, error) {
	method := existsPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp ExistsResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallStatus is the raw call to the Graphdriver.Status method
func CallStatus(url string, client *http.Client, req StatusRequest) (*StatusResponse, error) {
	method := statusPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp StatusResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallGetMetadata is the raw call to the Graphdriver.GetMetadata method
func CallGetMetadata(url string, client *http.Client, req GetMetadataRequest) (*GetMetadataResponse, error) {
	method := getMetadataPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp GetMetadataResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallCleanup is the raw call to the Graphdriver.Cleanup method
func CallCleanup(url string, client *http.Client, req CleanupRequest) (*CleanupResponse, error) {
	method := cleanupPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp CleanupResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallDiff is the raw call to the Graphdriver.Diff method
func CallDiff(url string, client *http.Client, req DiffRequest) (*DiffResponse, error) {
	method := diffPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return &DiffResponse{Stream: resp.Body}, nil
}

// CallChanges is the raw call to the Graphdriver.Changes method
func CallChanges(url string, client *http.Client, req ChangesRequest) (*ChangesResponse, error) {
	method := changesPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp ChangesResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallApplyDiff is the raw call to the Graphdriver.ApplyDiff method
func CallApplyDiff(url string, client *http.Client, req ApplyDiffRequest) (*ApplyDiffResponse, error) {
	method := applyDiffPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp ApplyDiffResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

// CallDiffSize is the raw call to the Graphdriver.CallDiffSize method
func CallDiffSize(url string, client *http.Client, req DiffSizeRequest) (*DiffSizeResponse, error) {
	method := diffSizePath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp DiffSizeResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}
