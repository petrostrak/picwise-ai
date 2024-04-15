package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/petrostrak/picwise-ai/db"
	"github.com/petrostrak/picwise-ai/types"
	"github.com/petrostrak/picwise-ai/view/generate"
	"net/http"
	"strconv"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	images, err := db.GetImagesByUserID(user.ID)
	if err != nil {
		return err
	}
	data := generate.ViewData{
		Images: images,
	}
	return render(w, r, generate.Index(data))
}

func HandleGenerateCreate(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	prompt := "a gopher fighting a crub"
	image := types.Image{
		UserID: user.ID,
		Prompt: prompt,
		Status: types.ImageStatusPending,
	}
	if err := db.CreateImage(&image); err != nil {
		return err
	}
	return render(w, r, generate.GalleryImage(image))
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	image, err := db.GetImageByID(id)
	if err != nil {
		return err
	}
	return render(w, r, generate.GalleryImage(image))
}
