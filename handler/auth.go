package handler

import (
	"net/http"

	"github.com/petrostrak/picwise-ai/view/auth"
)

func HandleSignInIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Login().Render(r.Context(), w)
}
