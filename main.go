package main

import (
	"net/http"
)

func main() {
	server_mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))

	server_mux.Handle("/", handler)
	server := http.Server{
		Addr:    ":8080",
		Handler: server_mux,
	}
	server.ListenAndServe()

}
