package sdk

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
)

const activatePath = "/Plugin.Activate"

// Mux exposes the minimal methods needed to set handler functions.
type Mux interface {
	HandleFunc(path string, fn func(http.ResponseWriter, *http.Request))
	// Add an implementation to Mux's route for Plugin.Activate
	AddImplementation(implementation string)
	http.Handler
}

// PluginManifest implements the Plugin.Activate JSON response.
type PluginManifest struct {
	Implements []string
}

// Handler is the base to create plugin handlers.
// It initializes connections and sockets to listen to.
type Handler struct {
	mux      *http.ServeMux
	manifest *PluginManifest
}

// NewHandler creates a new Handler with an http mux.
func NewHandler() Handler {
	mux := http.NewServeMux()
	internalManifest := &PluginManifest{}

	mux.HandleFunc(activatePath, func(w http.ResponseWriter, r *http.Request) {
		EncodeResponse(w, internalManifest, "")
	})

	return Handler{mux: mux, manifest: internalManifest}
}

// Serve sets up the handler to serve requests on the passed in listener
func (h Handler) Serve(l net.Listener) error {
	server := http.Server{
		Addr:    l.Addr().String(),
		Handler: h.mux,
	}
	return server.Serve(l)
}

// ServeTCP makes the handler to listen for request in a given TCP address.
// It also writes the spec file on the right directory for docker to read.
func (h Handler) ServeTCP(pluginName, addr string, tlsConfig *tls.Config) error {
	l, spec, err := newTCPListener(addr, pluginName, tlsConfig)
	if err != nil {
		return err
	}
	if spec != "" {
		defer os.Remove(spec)
	}
	return h.Serve(l)
}

// ServeUnix makes the handler to listen for requests in a unix socket.
// It also creates the socket file on the right directory for docker to read.
func (h Handler) ServeUnix(addr string, gid int) error {
	l, spec, err := newUnixListener(addr, gid)
	if err != nil {
		return err
	}
	if spec != "" {
		defer os.Remove(spec)
	}
	return h.Serve(l)
}

// ServeHTTP implements http.Handler and passes the request through to the contained mux.
// This method allows plugin handlers to be used with other Go HTTP frameworks.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// HandleFunc registers a function to handle a request path with.
func (h Handler) HandleFunc(path string, fn func(w http.ResponseWriter, r *http.Request)) {
	h.mux.HandleFunc(path, fn)
}

// AddImplementation adds the given implementation string to the manifest of the plugin handler.
// Unique implmentation names are only added once.
func (h Handler) AddImplementation(implementation string) {
	// Check the impl
	for _, v := range h.manifest.Implements {
		if v == implementation {
			return
		}
	}
	h.manifest.Implements = append(h.manifest.Implements, implementation)
}
