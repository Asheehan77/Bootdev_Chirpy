package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Asheehan77/Bootdev_Chirpy.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	secret         string
	plat           string
	polka_key      string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	plat := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polka_key := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	log.Println("Starting Server...")
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	server_mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))
	cfg := apiConfig{
		queries:   dbQueries,
		plat:      plat,
		secret:    secret,
		polka_key: polka_key,
	}

	server_mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", handler)))
	server_mux.HandleFunc("GET /admin/metrics", cfg.totalHits)
	server_mux.HandleFunc("POST /admin/reset", cfg.reset)
	server_mux.HandleFunc("GET /api/healthz", readinessHandler)
	server_mux.HandleFunc("GET /api/chirps", cfg.getchirpHandler)
	server_mux.HandleFunc("GET /api/chirps/{id}", cfg.getchirpidHandler)
	server_mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	server_mux.HandleFunc("POST /api/chirps", cfg.makechirpHandler)
	server_mux.HandleFunc("POST /api/login", cfg.loginHandler)
	server_mux.HandleFunc("POST /api/refresh", cfg.refreshHandler)
	server_mux.HandleFunc("POST /api/revoke", cfg.revokeHandler)
	server_mux.HandleFunc("PUT /api/users", cfg.updateUserHandler)
	server_mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirpHandler)
	server_mux.HandleFunc("POST /api/polka/webhooks", cfg.polkaHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: server_mux,
	}
	server.ListenAndServe()

}
