package main

import(
	"net/http"
)

func main(){
	server_mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: server_mux,
	}

	server.ListenAndServe()

	return
}