package handler

import (
	"net/http"

	"github.com/petrostrak/picwise-ai/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) {
	home.Index().Render(r.Context(), w)
}
