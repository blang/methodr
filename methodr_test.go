package methodr

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Status codes depending on method, for tests only
const (
	statusGET = iota + 200
	statusHEAD
	statusPOST
	statusPUT
	statusDELETE
	statusTRACE
	statusOPTIONS
	statusCONNECT
	statusPATCH
)

// Success handler, writes status code depending on request method
var succHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	status := 0
	switch r.Method {
	case methodGET:
		status = statusGET
	case methodHEAD:
		status = statusHEAD
	case methodPOST:
		status = statusPOST
	case methodPUT:
		status = statusPUT
	case methodDELETE:
		status = statusDELETE
	case methodTRACE:
		status = statusTRACE
	case methodOPTIONS:
		status = statusOPTIONS
	case methodCONNECT:
		status = statusCONNECT
	case methodPATCH:
		status = statusPATCH
	}
	w.WriteHeader(status)
})

// Handler returning 404
var deadHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
})

type routeTest struct {
	reqMethod string
	handler   http.Handler
	code      int
}

var routeTests = []routeTest{
	// Chain constructors
	{"GET", GET(succHandler), statusGET},
	{"HEAD", HEAD(succHandler), statusHEAD},
	{"POST", POST(succHandler), statusPOST},
	{"PUT", PUT(succHandler), statusPUT},
	{"DELETE", DELETE(succHandler), statusDELETE},
	{"TRACE", TRACE(succHandler), statusTRACE},
	{"OPTIONS", OPTIONS(succHandler), statusOPTIONS},
	{"CONNECT", CONNECT(succHandler), statusCONNECT},
	{"PATCH", PATCH(succHandler), statusPATCH},

	// Chained - Call chained method
	{"GET", POST(deadHandler).GET(succHandler), statusGET},
	{"HEAD", GET(deadHandler).HEAD(succHandler), statusHEAD},
	{"POST", GET(deadHandler).POST(succHandler), statusPOST},
	{"PUT", GET(deadHandler).PUT(succHandler), statusPUT},
	{"DELETE", GET(deadHandler).DELETE(succHandler), statusDELETE},
	{"TRACE", GET(deadHandler).TRACE(succHandler), statusTRACE},
	{"OPTIONS", GET(deadHandler).OPTIONS(succHandler), statusOPTIONS},
	{"CONNECT", GET(deadHandler).CONNECT(succHandler), statusCONNECT},
	{"PATCH", GET(deadHandler).PATCH(succHandler), statusPATCH},

	//Miss -> Defaulthandler
	{"GET", POST(succHandler), defaultHandlerStatusCode},

	// Special: HEAD uses GET if no HEAD handler is set and GET is available
	{"HEAD", GET(succHandler), statusHEAD},
	{"HEAD", POST(succHandler), defaultHandlerStatusCode},
	{"POST", GET(succHandler), defaultHandlerStatusCode},
	{"PUT", GET(succHandler), defaultHandlerStatusCode},
	{"DELETE", GET(succHandler), defaultHandlerStatusCode},
	{"TRACE", GET(succHandler), defaultHandlerStatusCode},
	{"OPTIONS", GET(succHandler), defaultHandlerStatusCode},
	{"CONNECT", GET(succHandler), defaultHandlerStatusCode},
	{"PATCH", GET(succHandler), defaultHandlerStatusCode},

	//Custom default handler
	{"GET", DEFAULT(succHandler), statusGET},
	{"POST", GET(deadHandler).DEFAULT(succHandler), statusPOST},
	{"POST", DEFAULT(succHandler).GET(deadHandler), statusPOST},

	//Unknown method
	{"UNKNOWN", DEFAULT(succHandler), 0},

	//Default handlers
	{"GET", DefaultHandlerMethodNotAllowed, http.StatusMethodNotAllowed},
}

// Test routeTest table
func TestRouting(t *testing.T) {
	for _, rt := range routeTests {
		resp := httptest.NewRecorder()
		req, err := http.NewRequest(rt.reqMethod, "", nil)
		if err != nil {
			t.Fatal(err)
		}

		rt.handler.ServeHTTP(resp, req)
		if resp.Code != rt.code {
			t.Errorf("Wrong statuscode, expected %d got %d", rt.code, resp.Code)
		}
	}
}

// Benchmark stdhandler without routing
func BenchmarkNoRoutingReference(b *testing.B) {
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		handler.ServeHTTP(resp, req)
	}
}

// Benchmark routing with Hit in routing table
func BenchmarkRoutingHit(b *testing.B) {
	resp := httptest.NewRecorder()
	endhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		b.Fatal(err)
	}
	handler := GET(endhandler)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		handler.ServeHTTP(resp, req)
	}
}

// Benchmark routing with Miss in routing table -> Default handler
func BenchmarkRoutingMissToDefault(b *testing.B) {
	resp := httptest.NewRecorder()
	endhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err := http.NewRequest("POST", "", nil)
	if err != nil {
		b.Fatal(err)
	}
	handler := GET(endhandler)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		handler.ServeHTTP(resp, req)
	}
}

// Benchmark routing with Miss in routing table -> Custom default handler
func BenchmarkRoutingMissToCustom(b *testing.B) {
	resp := httptest.NewRecorder()
	endhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err := http.NewRequest("POST", "", nil)
	if err != nil {
		b.Fatal(err)
	}
	handler := GET(endhandler).DEFAULT(endhandler)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		handler.ServeHTTP(resp, req)
	}
}

// func BenchmarkMissCustom(b *testing.B) {
//   resp := httptest.NewRecorder()
//   endhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//     w.WriteHeader(http.StatusOK)
//   })
//   req, err := http.NewRequest("POST", "", nil)
//   handler := GET(endhandler)

//   b.ResetTimer()
//   for n := 0; n < b.N; n++ {
//     handler.ServeHttp(resp, req)
//   }
// }
