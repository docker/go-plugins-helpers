package authorization

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	pluginType = authorization.AuthZApiImplements
	reqPath    = "/" + authorization.AuthZApiRequest
	resPath    = "/" + authorization.AuthZApiResponse
)

// Request is the structure that docker's requests are deserialized to.
type Request authorization.Request

// Response is the strucutre that the plugin's responses are serialized to.
type Response authorization.Response

// Plugin represent the interface a plugin must fulfill.
type Plugin interface {
	AuthZReq(Request) Response
	AuthZRes(Request) Response
}

// Handler forwards requests and responses between the docker daemon and the plugin.
type Handler struct {
	plugin Plugin
	sdk.Handler
}

// NewHandler initializes the request handler with a plugin implementation.
func NewHandler(plugin Plugin) *Handler {
	h := &Handler{plugin, sdk.NewHandler()}
	InitMux(h, plugin)
	return h
}

// InitMux initializes a compatible HTTP mux with routes for the specified driver. Can be used
// to combine multiple drivers into a single plugin mux.
func InitMux(h sdk.Mux, plugin Plugin) {
	handle(h, reqPath, func(req Request) Response {
		return plugin.AuthZReq(req)
	})

	handle(h, resPath, func(req Request) Response {
		return plugin.AuthZRes(req)
	})

	h.AddImplementation(pluginType)
}

type actionHandler func(Request) Response

func handle(h sdk.Mux, name string, actionCall actionHandler) {
	h.HandleFunc(name, func(w http.ResponseWriter, r *http.Request) {
		var (
			req Request
			d   = json.NewDecoder(r.Body)
		)
		d.UseNumber()
		if err := d.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		res := actionCall(req)

		sdk.EncodeResponse(w, res, res.Err)
	})
}
