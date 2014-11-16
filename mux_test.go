package zeus

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	m := New()

	// Default Mux has an emtpy map of handlers
	handlersCount(t, m, 0)

	// Default Mux does not have a NotFound handler func
	if got := m.NotFound; got != nil {
		t.Errorf("m.NotFound = %v, want %v", got, nil)
	}
}

func TestAdd(t *testing.T) {
	m, h := New(), &Handler{}

	m.add("FOO", h)

	// Mux has a map of handlers with one handler
	handlersCount(t, m, 1)

	// Mux has a FOO handler
	if got := m.handlers["FOO"][0]; got != h {
		t.Errorf("m.handlers[\"FOO\"][0] = %v, want %v", got, h)
	}
}

func TestGET(t *testing.T) {
	m, p := New(), "/"

	// Register GET handler
	m.GET(p, http.NotFound)

	// Mux has a map of handlers with two handlers
	handlersCount(t, m, 2)

	// Mux has a GET handler
	hasPattern(t, m, "GET", p)

	// Mux has a HEAD handler
	hasPattern(t, m, "HEAD", p)
}

func TestHEAD(t *testing.T) {
	m := New()
	defaultMethodTest(t, m, m.HEAD, "HEAD")
}

func TestPOST(t *testing.T) {
	m := New()
	defaultMethodTest(t, m, m.POST, "POST")
}

func TestPUT(t *testing.T) {
	m := New()
	defaultMethodTest(t, m, m.PUT, "PUT")
}

func TestDELETE(t *testing.T) {
	m := New()
	defaultMethodTest(t, m, m.DELETE, "DELETE")
}

func TestOPTIONS(t *testing.T) {
	m := New()
	defaultMethodTest(t, m, m.OPTIONS, "OPTIONS")
}

func TestPATCH(t *testing.T) {
	m := New()
	defaultMethodTest(t, m, m.PATCH, "PATCH")
}

type fields map[string]string

var tryTests = []struct {
	pattern string
	path    string
	fields  fields
}{
	{"/foo", "/foo", fields{}},
	{"/foo/:bar", "/foo/xyz", fields{"bar": "xyz"}},
	{"/foo/:bar/:baz", "/foo/xyz/123", fields{"bar": "xyz", "baz": "123"}},
}

func TestTry(t *testing.T) {
	for _, tt := range tryTests {
		h := &Handler{tt.pattern, http.NotFound}
		values, ok := h.try(tt.path)

		for key, value := range tt.fields {
			if ok != true || values.Get(key) != value {
				t.Fatalf("h.try(\"%s\") = %v, %v, wanted map[%s:[%s]], true",
					tt.path, values, ok, key, value)
			}
		}
	}
}

type methFunc func(string, http.HandlerFunc)

func defaultMethodTest(t *testing.T, m *Mux, fn methFunc, meth string) {
	fn("/", http.NotFound)

	handlersCount(t, m, 1)
	hasPattern(t, m, meth, "/")
}

func handlersCount(t *testing.T, m *Mux, c int) {
	if got, out := len(m.handlers), c; got != out {
		t.Errorf("len(m.handlers) = %v, want %v", got, out)
	}
}

func hasPattern(t *testing.T, m *Mux, meth, patt string) {
	if len(m.handlers[meth]) == 0 {
		t.Fatalf("missing handler for method %s", meth)
	}

	if got := m.handlers[meth][0].patt; got != patt {
		t.Errorf("m.handlers[\"%s\"][0].patt = %v, want %v", meth, got, patt)
	}
}
