package dkauthz

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/go-plugins-sdk/sdk"
)

const (
	defaultContentTypeV1_1        = "application/vnd.docker.plugins.v1.1+json"
	defaultImplementationManifest = `{"Implements": ["` + authorization.AuthZApiImplements + `"]}`

	activatePath = "/Plugin.Activate"
	reqPath      = "/" + authorization.AuthZApiRequest
	resPath      = "/" + authorization.AuthZApiResponse
)

type Request authorization.Request

type Response authorization.Response

// Plugin represent the interface a plugin must fulfill.
type Plugin interface {
	AuthZReq(Request) Response
	AuthZRes(Request) Response
}

// Handler forwards requests and responses between the docker daemon and the plugin.
type Handler struct {
	plugin Plugin
	mux    *http.ServeMux
}

// NewHandler initializes the request handler with a plugin implementation.
func NewHandler(plugin Plugin) *Handler {
	h := &Handler{plugin, http.NewServeMux()}
	h.initMux()
	return h
}

func (h *Handler) initMux() {
	h.mux.HandleFunc(activatePath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", defaultContentTypeV1_1)
		fmt.Fprintln(w, defaultImplementationManifest)
	})

	h.handle(reqPath, func(req Request) Response {
		return h.plugin.AuthZReq(req)
	})

	h.handle(resPath, func(req Request) Response {
		return h.plugin.AuthZRes(req)
	})
}

type actionHandler func(Request) Response

func (h *Handler) handle(name string, actionCall actionHandler) {
	h.mux.HandleFunc(name, func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(w, r)
		if err != nil {
			return
		}

		res := actionCall(req)

		encodeResponse(w, res)
	})
}

// ServeTCP makes the handler to listen for request in a given TCP address.
// It also writes the spec file on the right directory for docker to read.
func (h *Handler) ServeTCP(pluginName, addr string) error {
	return h.listenAndServe("tcp", addr, pluginName)
}

// ServeUnix makes the handler to listen for requests in a unix socket.
// It also creates the socket file on the right directory for docker to read.
func (h *Handler) ServeUnix(systemGroup, addr string) error {
	return h.listenAndServe("unix", addr, systemGroup)
}

func (h *Handler) listenAndServe(proto, addr, group string) error {
	var (
		l    net.Listener
		err  error
		spec string
	)

	server := http.Server{
		Addr:    addr,
		Handler: h.mux,
	}

	switch proto {
	case "tcp":
		l, spec, err = sdk.NewTCPListener(addr, group)
	case "unix":
		l, spec, err = sdk.NewUnixListener(addr, group)
	}

	if spec != "" {
		defer os.Remove(spec)
	}
	if err != nil {
		return err
	}

	return server.Serve(l)
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (req Request, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}

func encodeResponse(w http.ResponseWriter, res Response) {
	w.Header().Set("Content-Type", defaultContentTypeV1_1)
	if res.Err != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}
