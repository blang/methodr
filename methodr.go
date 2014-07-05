// Package methodr provides http.Handler compliant routing based on the request method.
package methodr

import (
	"net/http"
)

// Routing table with handler for all methods
type Mux struct {
	Get     http.Handler
	Head    http.Handler
	Post    http.Handler
	Put     http.Handler
	Delete  http.Handler
	Trace   http.Handler
	Options http.Handler
	Connect http.Handler
	Patch   http.Handler
	Default http.Handler // Default handler in case of miss
}

const (
	methodGET     = "GET"
	methodHEAD    = "HEAD"
	methodPOST    = "POST"
	methodPUT     = "PUT"
	methodDELETE  = "DELETE"
	methodTRACE   = "TRACE"
	methodOPTIONS = "OPTIONS"
	methodCONNECT = "CONNECT"
	methodPATCH   = "PATCH"
)

var (
	// Changeable global default handler in case of miss on routing table. Default: DefaultHandlerMethodNotAllowed
	DefaultHandler = DefaultHandlerMethodNotAllowed

	// Default handler returning http.StatusMethodNotAllowed
	DefaultHandlerMethodNotAllowed = http.HandlerFunc(defaultHandleMethodNotAllowed)

	defaultHandlerStatusCode = http.StatusMethodNotAllowed // For tests only - consistent with DefaultHandler
)

func defaultHandleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// ServeHTTP routes requests depending on method routing table.
// In case of a miss the custom DEFAULT handler is used, otherwise the global DefaultHandler.
// HEAD requests are delegated to GET if there's a GET handler available and no HEAD handler set.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case methodGET:
		if m.Get != nil {
			m.Get.ServeHTTP(w, r)
		} else {
			m.handleDefault(w, r)
		}
	case methodHEAD:
		if m.Head != nil {
			m.Head.ServeHTTP(w, r)
		} else if m.Get != nil { // Special case: HEAD uses GET if GET available and HEAD not set.
			m.Get.ServeHTTP(w, r)
		} else {
			m.handleDefault(w, r)
		}
	case methodPOST:
		if m.Post != nil {
			m.Post.ServeHTTP(w, r)
		} else {
			m.handleDefault(w, r)
		}
	case methodPUT:
		if m.Put != nil {
			m.Put.ServeHTTP(w, r)
		} else {
			m.handleDefault(w, r)
		}
	case methodDELETE:
		if m.Delete != nil {
			m.Delete.ServeHTTP(w, r)
		} else {
			m.handleDefault(w, r)
		}
	case methodTRACE:
		if m.Trace != nil {
			m.Trace.ServeHTTP(w, r)
		} else {
			m.handleDefault(w, r)
		}
	case methodOPTIONS:
		if m.Options != nil {
			m.Options.ServeHTTP(w, r)
		} else {
			m.handleDefault(w, r)
		}
	case methodCONNECT:
		if m.Connect != nil {
			m.Connect.ServeHTTP(w, r)
		} else {
			m.handleDefault(w, r)
		}
	case methodPATCH:
		if m.Patch != nil {
			m.Patch.ServeHTTP(w, r)
		} else {
			m.handleDefault(w, r)
		}
	default:
		m.handleDefault(w, r)
	}
}

// Use global default handler if there is no custom default handler HDEFAULT
func (m *Mux) handleDefault(w http.ResponseWriter, r *http.Request) {
	if m.Default == nil {
		DefaultHandler.ServeHTTP(w, r)
	} else {
		m.Default.ServeHTTP(w, r)
	}
}

// DEFAULT sets default handler used in case of miss on routing table.
func DEFAULT(h http.Handler) *Mux {
	return &Mux{
		Default: h,
	}
}

// DEFAULT sets default handler in case of miss on routing table
func (m *Mux) DEFAULT(h http.Handler) *Mux {
	m.Default = h
	return m
}

// GET sets the handler used for GET method requests.
// HEAD requests are delegated to GET if there's a GET handler available and no HEAD handler set.
func GET(h http.Handler) *Mux {
	return &Mux{
		Get: h,
	}
}

// GET sets the handler used for GET method requests.
// HEAD requests are delegated to GET if there's a GET handler available and no HEAD handler set.
func (m *Mux) GET(h http.Handler) *Mux {
	m.Get = h
	return m
}

// HEAD sets the handler used for HEAD method requests.
// HEAD requests are delegated to GET if there's a GET handler available and no HEAD handler set.
func HEAD(h http.Handler) *Mux {
	return &Mux{
		Head: h,
	}
}

// HEAD sets the handler used for HEAD method requests.
// HEAD requests are delegated to GET if there's a GET handler available and no HEAD handler set.
func (m *Mux) HEAD(h http.Handler) *Mux {
	m.Head = h
	return m
}

// POST sets the handler used for POST method requests.
func POST(h http.Handler) *Mux {
	return &Mux{
		Post: h,
	}
}

// POST sets the handler used for POST method requests.
func (m *Mux) POST(h http.Handler) *Mux {
	m.Post = h
	return m
}

// PUT sets the handler used for PUT method requests.
func PUT(h http.Handler) *Mux {
	return &Mux{
		Put: h,
	}
}

// PUT sets the handler used for PUT method requests.
func (m *Mux) PUT(h http.Handler) *Mux {
	m.Put = h
	return m
}

// DELETE sets the handler used for DELETE method requests.
func DELETE(h http.Handler) *Mux {
	return &Mux{
		Delete: h,
	}
}

// DELETE sets the handler used for DELETE method requests.
func (m *Mux) DELETE(h http.Handler) *Mux {
	m.Delete = h
	return m
}

// TRACE sets the handler used for TRACE method requests.
func TRACE(h http.Handler) *Mux {
	return &Mux{
		Trace: h,
	}
}

// TRACE sets the handler used for TRACE method requests.
func (m *Mux) TRACE(h http.Handler) *Mux {
	m.Trace = h
	return m
}

// OPTIONS sets the handler used for OPTIONS method requests.
func OPTIONS(h http.Handler) *Mux {
	return &Mux{
		Options: h,
	}
}

// OPTIONS sets the handler used for OPTIONS method requests.
func (m *Mux) OPTIONS(h http.Handler) *Mux {
	m.Options = h
	return m
}

// CONNECT sets the handler used for CONNECT method requests.
func CONNECT(h http.Handler) *Mux {
	return &Mux{
		Connect: h,
	}
}

// CONNECT sets the handler used for CONNECT method requests.
func (m *Mux) CONNECT(h http.Handler) *Mux {
	m.Connect = h
	return m
}

// PATCH sets the handler used for PATCH method requests.
func PATCH(h http.Handler) *Mux {
	return &Mux{
		Patch: h,
	}
}

// PATCH sets the handler used for PATCH method requests.
func (m *Mux) PATCH(h http.Handler) *Mux {
	m.Patch = h
	return m
}
