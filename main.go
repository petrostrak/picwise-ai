package main

import (
	"embed"
	"github.com/petrostrak/picwise-ai/db"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/petrostrak/picwise-ai/handler"
	"github.com/petrostrak/picwise-ai/pkg/sb"
)

//go:embed public
var FS embed.FS

func Init() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}

func main() {
	if err := Init(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handler.WithUser)

	router.Handle("/*", public())
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.HandleSignInIndex))
	router.Get("/login/provider/google", handler.Make(handler.HandleLoginWithGoogle))
	router.Get("/signup", handler.Make(handler.HandleSignupIndex))
	router.Get("/auth/callback", handler.Make(handler.HandleAuthCallback))
	router.Post("/logout", handler.Make(handler.HandleLogoutCreate))
	router.Post("/login", handler.Make(handler.HandleLoginCreate))
	router.Post("/signup", handler.Make(handler.HandleSignupCreate))
	router.Post("/replicate/callback/{userID}/{batchID}", handler.Make(handler.HandleReplicateCallback))

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth)
		auth.Get("/account/setup", handler.Make(handler.HandleAccountSetupIndex))
		auth.Post("/account/setup", handler.Make(handler.HandleAccountSetupCreate))
	})

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth, handler.WithAccountSetup)
		auth.Get("/settings", handler.Make(handler.HandleSettingsIndex))
		auth.Put("/settings/account/profile", handler.Make(handler.HandleSettingsUsernameUpdate))

		auth.Post("/auth/reset-password", handler.Make(handler.HandleResetPasswordCreate))
		auth.Put("/auth/reset-password", handler.Make(handler.HandleResetPasswordUpdate))
		auth.Get("/auth/reset-password", handler.Make(handler.HandleResetPasswordIndex))

		auth.Get("/generate", handler.Make(handler.HandleGenerateIndex))
		auth.Post("/generate", handler.Make(handler.HandleGenerateCreate))

		auth.Get("/generate/image/status/{id}", handler.Make(handler.HandleGenerateImageStatus))
	})

	port := os.Getenv("HTTP_LISTEN_ADDRESS")
	slog.Info("application running on", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
