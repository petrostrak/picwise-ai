package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/petrostrak/picwise-ai/db"
	"github.com/petrostrak/picwise-ai/types"
	"github.com/uptrace/bun"
	"net/http"
)

const succeded = "succeded"

type ReplicateResponse struct {
	Input struct {
		Prompt string `json:"prompt"`
	} `json:"input"`
	Status string   `json:"status"`
	Output []string `json:"output"`
}

func HandleReplicateCallback(w http.ResponseWriter, r *http.Request) error {
	var resp ReplicateResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}

	if resp.Status != succeded {
		return fmt.Errorf("replicate callback responsed with: %s", resp.Status)
	}

	batchID, err := uuid.Parse(chi.URLParam(r, "batchID"))
	if err != nil {
		return err
	}

	images, err := db.GetImagesByBatchID(batchID)
	if err != nil {
		return err
	}

	if len(images) != len(resp.Output) {
		return fmt.Errorf("replicate callback output doesn't match the number of images fetched")
	}

	err = db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for i, imageURL := range resp.Output {
			images[i].Status = types.ImageStatusCompleted
			images[i].ImageLocation = imageURL
			images[i].Prompt = resp.Input.Prompt
			if err := db.UpdateImage(&images[i]); err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}
