package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := chi.NewMux()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	port := os.Getenv("HTTP_LISTEN_ADDRESS")
	slog.Info("application running on", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
