package mountpoint

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/volume/mountpoint"
	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	manifest = `{"Implements": ["` + mountpoint.MountPointAPIImplements + `"]}`
	propPath = "/" + mountpoint.MountPointAPIProperties
	attPath  = "/" + mountpoint.MountPointAPIAttach
	detPath  = "/" + mountpoint.MountPointAPIDetach
)

// Plugin is the mount point plugin you are implementing
type Plugin interface {
	Properties(mountpoint.PropertiesRequest) mountpoint.PropertiesResponse
	Attach(mountpoint.AttachRequest) mountpoint.AttachResponse
	Detach(mountpoint.DetachRequest) mountpoint.DetachResponse
}

// Handler is a plugin API request receiver object
type Handler struct {
	plugin Plugin
	sdk.Handler
}

// NewHandler creates a new API receiver object from a plugin and
// starts listening for API requests
func NewHandler(plugin Plugin) *Handler {
	h := &Handler{plugin, sdk.NewHandler(manifest)}
	h.initMux()
	return h
}

func (h *Handler) initMux() {
	h.HandleFunc(attPath, func(w http.ResponseWriter, r *http.Request) {
		var req mountpoint.AttachRequest
		d := json.NewDecoder(r.Body)

		if err := d.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		res := h.plugin.Attach(req)

		sdk.EncodeResponse(w, res, res.Err != "")
	})

	h.HandleFunc(detPath, func(w http.ResponseWriter, r *http.Request) {
		var req mountpoint.DetachRequest
		d := json.NewDecoder(r.Body)

		if err := d.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		res := h.plugin.Detach(req)

		sdk.EncodeResponse(w, res, res.Err != "")
	})

	h.HandleFunc(propPath, func(w http.ResponseWriter, r *http.Request) {
		var req mountpoint.PropertiesRequest
		d := json.NewDecoder(r.Body)

		if err := d.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		res := h.plugin.Properties(req)

		sdk.EncodeResponse(w, res, res.Err != "")
	})
}
