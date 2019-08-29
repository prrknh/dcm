package main

import (
	"github.com/prrknh/dcm/handler"
	"net/http"
)

func main() {

	// check docker daemon is running
	// check argument image name exists

	http.HandleFunc("/create", handler.CreateContainer())
	http.ListenAndServe(":5000", nil)
}
