package handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/petrostrak/picwise-ai/db"
	"github.com/petrostrak/picwise-ai/pkg/kit/validate"
	"github.com/petrostrak/picwise-ai/types"
	"github.com/petrostrak/picwise-ai/view/generate"
	"github.com/replicate/replicate-go"
	"github.com/uptrace/bun"
	"log"
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

	genParams := GenerateImageParams{
		Prompt:  params.Prompt,
		Amount:  params.Amount,
		UserID:  user.ID,
		BatchID: uuid.New(),
	}

	if err := generateImages(r.Context(), genParams); err != nil {
		return err
	}

	err = db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for i := 0; i < params.Amount; i++ {
			image := types.Image{
				UserID:  user.ID,
				Status:  types.ImageStatusPending,
				BatchID: genParams.BatchID,
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

type GenerateImageParams struct {
	Prompt  string
	Amount  int
	UserID  uuid.UUID
	BatchID uuid.UUID
}

func generateImages(ctx context.Context, params GenerateImageParams) error {
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	input := replicate.PredictionInput{
		"prompt":      params.Prompt,
		"num_outputs": params.Amount,
	}
	webhook := replicate.Webhook{
		URL:    fmt.Sprintf("https://webhook.site/214305cd-0416-4f6b-890e-69cb47f43c3a/%s/%s", params.UserID, params.BatchID),
		Events: []replicate.WebhookEventType{"completed"},
	}
	_, err = r8.CreatePrediction(ctx, "d70beb400d223e6432425a5299910329c6050c6abcf97b8c70537d6a1fcb269a", input, &webhook, false)
	return err
}
