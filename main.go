package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	server_mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))
	cfg := apiConfig{}

	server_mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", handler)))
	server_mux.HandleFunc("GET /admin/metrics", cfg.totalHits)
	server_mux.HandleFunc("POST /admin/reset", cfg.reset)
	server_mux.HandleFunc("GET /api/healthz", readinessHandler)
	server_mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: server_mux,
	}
	server.ListenAndServe()

}

func readinessHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	data := []byte("OK")
	writer.Write(data)
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

func (cfg *apiConfig) reset(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

func validateHandler(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
	}

	param := parameters{}
	res := response{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&param)
	fmt.Println(param.Body)
	if err != nil {
		fmt.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
		res.Error = "Something went wrong"
	} else {
		if len(param.Body) > 140 {
			writer.WriteHeader(400)
			res.Valid = false
			res.Error = "Chirp is too long"
		} else {
			writer.WriteHeader(200)
			res.Valid = true
		}
	}

	jres, err := json.Marshal(res)

	if err != nil {
		log.Printf("Error encoding response: %s", err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(jres))
}
