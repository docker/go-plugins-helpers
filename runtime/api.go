package runtime

// See https://github.com/docker/docker/blob/master/experimental/plugins_graphdriver.md

import (
	"net/http"
	
	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	// DefaultDockerRootDirectory is the default directory where runtime drivers will be created.
	DefaultDockerRootDirectory = "/var/lib/docker/runtime"

	manifest = `{"Implements": ["RuntimeDriver"]}`
	pathPath = "/RuntimeDriver.Path"
	argsPath = "/RuntimeDriver.Args"
)

// PathResponse is the Path response
type PathResponse struct {
	Path string `json:"path,omitempty"`
}

// ArgsResponse is the Args response
type ArgsResponse struct {
	Args []string `json:"args,omitempty"`
}

// Plugin represent the interface a driver must fulfill.
type Plugin interface {
	Path() PathResponse
	Args() ArgsResponse
}

// Handler forwards requests and responses between the docker daemon and the plugin.
type Handler struct {
	plugin Plugin
	sdk.Handler
}

// NewHandler initializes the request handler with a plugin implementation.
func NewHandler(plugin Plugin) *Handler {
	h := &Handler{plugin, sdk.NewHandler(manifest)}
	h.initMux()
	return h
}

func (h *Handler) initMux() {
	h.handlePath(pathPath, func() PathResponse {
		return h.plugin.Path()
	})

	h.handleArgs(argsPath, func() ArgsResponse {
		return h.plugin.Args()
	})
}

func (h *Handler) handlePath(name string, pathCall func() PathResponse) {
	h.HandleFunc(name, func(w http.ResponseWriter, r *http.Request) {
		res := pathCall()

		sdk.EncodeResponse(w, res, "")
	})
}

func (h *Handler) handleArgs(name string, argsCall func() ArgsResponse) {
	h.HandleFunc(name, func(w http.ResponseWriter, r *http.Request) {
		res := argsCall()

		sdk.EncodeResponse(w, res, "")
	})
}
