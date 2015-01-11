package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func handler(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

//a fake server for test_test.go
func main() {
	goji.Get("/*", handler)
	goji.Serve()
}
