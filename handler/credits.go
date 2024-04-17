package handler

import (
	"github.com/petrostrak/picwise-ai/view/credits"
	"net/http"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, credits.Index())
}
