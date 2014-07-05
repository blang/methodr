methodr for golang [![Build Status](https://drone.io/github.com/blang/methodr/status.png)](https://drone.io/github.com/blang/methodr/latest) [![GoDoc](https://godoc.org/github.com/blang/methodr?status.png)](https://godoc.org/github.com/blang/methodr) [![Coverage Status](https://img.shields.io/coveralls/blang/methodr.svg)](https://coveralls.io/r/blang/methodr?branch=master)
======

methodr provides routing based on the request method written in golang. Fully http.Handler compliant.

Usage
-----
```bash
$ go get github.com/blang/methodr
```

```go
import github.com/blang/methodr

http.Handle("/foo", methodr.GET(getHandler).POST(postHandler).DEFAULT(notFoundHandler))
```

Also check the [GoDocs](http://godoc.org/github.com/blang/methodr).

Why should I use this lib?
-----

- Simple Interface
- Fully http.Handler compliant (supports every good routing environment)
- No reflection/reflection/nasty stuff
- No allocs as handler
- Fully tested (Coverage >99%)
- Fast (See [Benchmarks](#benchmarks))
- Only Stdlib


Features
-----

- Helper functions for all methods
- Full control if your setup is more complex
- Global or routebased default/catchall handler


Example
-----

Have a look at full and runnable examples in [examples/main.go](examples/main.go).
Also check the [GoDocs](http://godoc.org/github.com/blang/methodr).

```go
import github.com/blang/methodr

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
```

To get a list of all available methods see the [GoDocs](http://godoc.org/github.com/blang/methodr).

Benchmarks
-----

```
BenchmarkNoRoutingReference    500000000             7.46 ns/op         0 B/op          0 allocs/op
BenchmarkRoutingHit             50000000            31.2  ns/op         0 B/op          0 allocs/op
BenchmarkRoutingMissToDefault   50000000            35.0  ns/op         0 B/op          0 allocs/op
BenchmarkRoutingMissToCustom    50000000            36.2  ns/op         0 B/op          0 allocs/op
```

See benchmark cases at [methodr_test.go](methodr_test.go)


Contribution
-----

Feel free to make a pull request. For bigger changes create a issue first to discuss about it.


License
-----

See [LICENSE](LICENSE) file.
