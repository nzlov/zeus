package zeus

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Mux struct {
	handlers map[string][]*Handler
	NotFound http.HandlerFunc
}

type Handler struct {
	patt string
	http.HandlerFunc
}

// Returns a new Mux instance.
func New() *Mux {
	return &Mux{make(map[string][]*Handler), nil}
}

// Listen is a shorthand way of doing http.ListenAndServe.
func (m *Mux) Listen(port string) {
	fmt.Printf("Listening: %s\n", port[1:])
	http.ListenAndServe(port, m)
}

func (m *Mux) add(meth string, handler *Handler) {
	m.handlers[meth] = append(m.handlers[meth], handler)
}

// GET adds a new route for GET requests.
func (m *Mux) GET(patt string, handler http.HandlerFunc) {
	m.add("GET", &Handler{patt, handler})
	m.add("HEAD", &Handler{patt, handler})
}

// GET adds a new route for GET requests.
func (m *Mux) HEAD(patt string, handler http.HandlerFunc) {
	m.add("HEAD", &Handler{patt, handler})
}

// POST adds a new route for POST requests.
func (m *Mux) POST(patt string, handler http.HandlerFunc) {
	m.add("POST", &Handler{patt, handler})
}

// PUT adds a new route for PUT requests.
func (m *Mux) PUT(patt string, handler http.HandlerFunc) {
	m.add("PUT", &Handler{patt, handler})
}

// DELETE adds a new route for DELETE requests.
func (m *Mux) DELETE(patt string, handler http.HandlerFunc) {
	m.add("DELETE", &Handler{patt, handler})
}

// OPTIONS adds a new route for OPTIONS requests.
func (m *Mux) OPTIONS(patt string, handler http.HandlerFunc) {
	m.add("OPTIONS", &Handler{patt, handler})
}

// PATCH adds a new route for PATCH requests.
func (m *Mux) PATCH(patt string, handler http.HandlerFunc) {
	m.add("PATCH", &Handler{patt, handler})
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Map over the registered handlers for
	// the current request (if there is any).
	for _, handler := range m.handlers[r.Method] {
		// Try the pattern against the URL path.
		if vars, ok := handler.try(r.URL.Path); ok {
			// Add params to URL query.
			r.URL.RawQuery += "&" + vars.Encode()
			// Serve handlers.
			handler.ServeHTTP(w, r)
			return
		}
	}
	// Custom 404 handler?
	if m.NotFound != nil {
		w.WriteHeader(404)
		m.NotFound.ServeHTTP(w, r)
		return
	}
	// Default 404.
	http.NotFound(w, r)
}

func (h *Handler) try(path string) (url.Values, bool) {
	// Pattern and URL Path segments.
	ps := strings.Split(h.patt[1:], "/")
	us := strings.Split(path[1:], "/")

	// If the patt and URL slices
	// have different lengths we
	// already know it's bad.
	if len(ps) != len(us) {
		return nil, false
	}

	// Compiled.
	var cs string
	// Parameters.
	uv := url.Values{}

	for idx, part := range ps {
		// Character is at least :x
		if len(part) > 1 && part[:1] == ":" {
			// Add to parameters.
			uv.Add(part[1:], us[idx])
			// Add URL seg.
			cs += "/" + us[idx]
			continue
		}
		// Add patt seg.
		cs += "/" + part
	}
	// Match?
	if cs == path {
		return uv, true
	}

	return nil, false
}
