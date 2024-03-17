package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/petrostrak/picwise-ai/handler"
	"github.com/petrostrak/picwise-ai/pkg/supabase"
)

//go:embed public
var FS embed.FS

func Init() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	return supabase.Init()
}

func main() {
	if err := Init(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handler.WithUser)

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.HandleSignInIndex))

	port := os.Getenv("HTTP_LISTEN_ADDRESS")
	slog.Info("application running on", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
