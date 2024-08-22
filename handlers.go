package main

import (
	"fmt"
	"io"
	"net/http"
)

var routeMap = map[string]func(w http.ResponseWriter, req *http.Request){
	"/status": statusHandler,
}

func statusHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, fmt.Sprintf("Server is Up and Running"))
}
