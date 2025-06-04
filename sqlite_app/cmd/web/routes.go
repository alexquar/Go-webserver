package main

import (
	"net/http"
)

func setupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getHome)
	return mux
}
