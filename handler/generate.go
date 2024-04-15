package handler

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/petrostrak/picwise-ai/db"
	"github.com/petrostrak/picwise-ai/pkg/kit/validate"
	"github.com/petrostrak/picwise-ai/types"
	"github.com/petrostrak/picwise-ai/view/generate"
	"github.com/uptrace/bun"
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
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		return err
	}
	params := generate.FormParams{
		Prompt: r.FormValue("prompt"),
		Amount: amount,
	}
	var errors generate.FormErrors
	if amount <= 0 || amount > 4 {
		errors.Amount = "Please enter a valid amount"
		return render(w, r, generate.Form(params, errors))
	}
	if ok := validate.New(params, validate.Fields{
		"Prompt": validate.Rules(validate.Min(5), validate.Max(100)),
	}).Validate(&errors); !ok {
		return render(w, r, generate.Form(params, errors))
	}
	err = db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		batchID := uuid.New()
		for i := 0; i < params.Amount; i++ {
			image := types.Image{
				UserID:  user.ID,
				Prompt:  params.Prompt,
				Status:  types.ImageStatusPending,
				BatchID: batchID,
			}
			if err := db.CreateImage(&image); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return hxRedirect(w, r, "/generate")
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
