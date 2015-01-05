package zeus

import (
	"fmt"
	"net/http"
	"strings"
)

// Mux contains a map of handlers and the NotFound handler func.
type Mux struct {
	handlers map[string][]*Handler
	NotFound http.HandlerFunc
}

// Handler contains the pattern and handler func.
type Handler struct {
	patt string
	vars bool
	wild bool
	http.HandlerFunc
}

var vars = map[*http.Request]map[string]string{}

// New returns a new Mux instance.
func New() *Mux {
	return &Mux{make(map[string][]*Handler), nil}
}

// Get route variables for the current request.
func Vars(r *http.Request) map[string]string {
	if v, ok := vars[r]; ok {
		return v
	}
	return nil
}

// Get route variable from the current request.
func Var(r *http.Request, n string) string {
	var v string
	if m := Vars(r); m != nil {
		v, _ = m[n]
	}
	return v
}

// Listen is a shorthand way of doing http.ListenAndServe.
func (m *Mux) Listen(port string) {
	fmt.Printf("Listening: %s\n", port[1:])
	http.ListenAndServe(port, m)
}

func (m *Mux) add(meth, patt string, handler http.HandlerFunc) {
	h := &Handler{patt, false, false, handler}
	for _, v := range patt {
		if v == ':' {
			h.vars = true
		} else if v == '*' {
			h.wild = true
		}
	}
	m.handlers[meth] = append(
		m.handlers[meth],
		h,
	)
}

// GET adds a new route for GET requests.
func (m *Mux) GET(patt string, handler http.HandlerFunc) {
	m.add("GET", patt, handler)
	m.add("HEAD", patt, handler)
}

// HEAD adds a new route for HEAD requests.
func (m *Mux) HEAD(patt string, handler http.HandlerFunc) {
	m.add("HEAD", patt, handler)
}

// POST adds a new route for POST requests.
func (m *Mux) POST(patt string, handler http.HandlerFunc) {
	m.add("POST", patt, handler)
}

// PUT adds a new route for PUT requests.
func (m *Mux) PUT(patt string, handler http.HandlerFunc) {
	m.add("PUT", patt, handler)
}

// DELETE adds a new route for DELETE requests.
func (m *Mux) DELETE(patt string, handler http.HandlerFunc) {
	m.add("DELETE", patt, handler)
}

// OPTIONS adds a new route for OPTIONS requests.
func (m *Mux) OPTIONS(patt string, handler http.HandlerFunc) {
	m.add("OPTIONS", patt, handler)
}

// PATCH adds a new route for PATCH requests.
func (m *Mux) PATCH(patt string, handler http.HandlerFunc) {
	m.add("PATCH", patt, handler)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := len(r.URL.Path)
	// Redirect trailing slash URL's.
	if l > 1 && r.URL.Path[l-1:] == "/" {
		http.Redirect(w, r, r.URL.Path[:l-1], 301)
		return
	}
	// Map over the registered handlers for
	// the current request (if there is any).
	for _, handler := range m.handlers[r.Method] {
		// If the route doesn't have any
		// named parameters or wildcards.
		if !handler.vars && !handler.wild {
			if handler.patt == r.URL.Path {
				handler.ServeHTTP(w, r)
				return
			}
			continue
		}
		// Compare pattern to URL.
		if ok := handler.try(r); ok {
			handler.ServeHTTP(w, r)
			delete(vars, r)
			return
		}
		delete(vars, r)
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

func (h *Handler) try(r *http.Request) bool {
	up := r.URL.Path
	us := strings.Split(up[1:], "/")
	ps := strings.Split(h.patt[1:], "/")
	pl := len(ps)

	if pl > len(us) {
		return false
	}

	if h.vars {
		vars[r] = map[string]string{}
	}

	var cs string

	for idx, part := range ps {
		// Wildcard segment.
		if h.wild && part == "*" {
			cs += "/" + us[idx]
			continue
		}
		// Named parameter segment.
		if h.vars && part[:1] == ":" {
			cs += "/" + us[idx]
			vars[r][part[1:]] = us[idx]
			continue
		}
		// Regular.
		cs += "/" + part
	}

	// If the pattern ends with *
	if h.wild && h.patt[len(h.patt):] == "*" {
		return up[0:len(cs)] == cs
	}

	return cs == up
}
