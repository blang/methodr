package main

import (
	"fmt"
	"github.com/blang/methodr"
	"net/http"
)

var (
	getHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Requested method: %s\n", r.Method)
	})
	postHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Post me something")
	})
	notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
)

const listen = ":9000"

func main() {

	// Setup 1: Restrict route to single method
	// POST/PUT/... will result in StatusMethodNotAllowed
	// Note: HEAD will use GET unless HEAD handler was set.
	http.Handle("/foo1", methodr.GET(getHandler))

	// Setup 2: Restrict route with custom default handler
	// POST/PUT/... will result in StatusNotFound
	http.Handle("/foo2", methodr.GET(getHandler).DEFAULT(notFoundHandler))

	//NOTE: If you're not happy with the global default handler, you can set it:
	// methodr.DefaultHandler = myBadRequestHandler

	// Setup 3: Route depending on method
	// Use chains to specify different routes, order is nonrelevant.
	// PATCH/PUT/... will result in StatusMethodNotAllowed
	http.Handle("/foo3", methodr.GET(getHandler).POST(postHandler))
	// Equivalent to: http.Handle("/foo3", methodr.POST(postHandler).GET(getHandler))

	// Setup 4: Setup more complicated routes
	mux := &methodr.Mux{
		Get:     getHandler,
		Post:    postHandler,
		Patch:   postHandler,
		Default: notFoundHandler,
	}
	http.Handle("/foo4", mux)
	// Equivalent to:
	// http.Handle("/foo4", methodr.GET(getHandler).POST(postHandler).PATCH(postHandler).DEFAULT(notFoundHandler))

	fmt.Printf("HTTP Server listening on %s, exit with ctrl+c\n", listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		fmt.Printf("HTTP Server exited with: %s\n", err)
	}
}
