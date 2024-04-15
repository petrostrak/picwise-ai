package handler

import (
	"github.com/petrostrak/picwise-ai/view/generate"
	"net/http"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, generate.Index())
}
