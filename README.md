# Zeus

Zeus is a very simple and fast HTTP router for Go, nothing more, nothing less.

#### Install

`go get github.com/daryl/zeus`

#### Usage

```go
package main

import (
    "fmt"
    "github.com/daryl/zeus"
    "net/http"
)

func main() {
    mux := zeus.New()
    // Supports named parameters
    mux.GET("/users/:id", showUser)
    // Custom 404 handler
    mux.NotFound = notFound
    // Listen and serve
    mux.Listen(":4545")
}

func showUser(w http.ResponseWriter, r *http.Request) {
    // Extract parameter value
    id := r.URL.Query().Get("id")

    fmt.Fprintf(w, "User ID: %s", id)
}

func notFound(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Nothing to see here")
}
```

#### Documentation

For further documentation, check out [GoDoc](http://godoc.org/github.com/daryl/zeus).
