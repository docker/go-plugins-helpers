package authorization

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	reqPath = "/" + authorization.AuthZApiRequest
	resPath = "/" + authorization.AuthZApiResponse
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

// NewHandler initializes the request handler with a plugin implementation.
func NewHandler(plugin Plugin) *sdk.Handler {
	h := sdk.NewHandler()
	RegisterDriver(plugin, h)
	return h
}

// RegisterDriver registers the plugin to the SDK handler.
func RegisterDriver(plugin Plugin, h *sdk.Handler) {
	h.RegisterDriver(authorization.AuthZApiImplements, func(h *sdk.Handler) {
		initMux(plugin, h)
	})
}

func initMux(plugin Plugin, h *sdk.Handler) {
	handle(h, reqPath, func(req Request) Response {
		return plugin.AuthZReq(req)
	})

	handle(h, resPath, func(req Request) Response {
		return plugin.AuthZRes(req)
	})
}

type actionHandler func(Request) Response

func handle(h *sdk.Handler, name string, actionCall actionHandler) {
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
