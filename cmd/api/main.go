package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"example.com/pz8-mongo/internal/db"
	"example.com/pz8-mongo/internal/notes"
)

func main() {
	uri := getenv("MONGO_URI", "mongodb://root:secret@mongo:27017/?authSource=admin")
	dbName := getenv("MONGO_DB", "pz8")
	addr := getenv("HTTP_ADDR", ":8080") // слушаем на всех интерфейсах

	deps, err := db.ConnectMongo(context.Background(), uri, dbName)
	if err != nil { log.Fatal("mongo connect:", err) }
	defer deps.Client.Disconnect(context.Background())

	repo, err := notes.NewRepo(deps.Database)
	if err != nil { log.Fatal("notes repo:", err) }
	h := notes.NewHandler(repo)

	r := chi.NewRouter()
	// Простые CORS, чтобы удобно тестировать из браузера/извне
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.Mount("/api/v1/notes", h.Routes())

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}
