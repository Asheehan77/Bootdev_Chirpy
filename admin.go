package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) reset(writer http.ResponseWriter, req *http.Request) {
	if cfg.plat == "dev" {
		writer.WriteHeader(http.StatusOK)
		cfg.fileserverHits.Store(0)
		err := cfg.queries.Reset(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Server Reset")

	} else {
		writer.WriteHeader(http.StatusForbidden)
	}

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) totalHits(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	hitpage := fmt.Sprintf("<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", (int(cfg.fileserverHits.Load())))
	data := []byte(hitpage)
	writer.Write(data)
}
