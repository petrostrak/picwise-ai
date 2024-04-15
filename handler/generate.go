package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/petrostrak/picwise-ai/types"
	"github.com/petrostrak/picwise-ai/view/generate"
	"log/slog"
	"net/http"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	data := generate.ViewData{
		Images: []types.Image{},
	}
	return render(w, r, generate.Index(data))
}

func HandleGenerateCreate(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, generate.GalleryImage(types.Image{Status: types.ImageStatusPending}))
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	slog.Info("checking image status", "id", id)
	return render(w, r, generate.GalleryImage(types.Image{Status: types.ImageStatusPending}))
}
