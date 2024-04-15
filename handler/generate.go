package handler

import (
	"github.com/petrostrak/picwise-ai/types"
	"github.com/petrostrak/picwise-ai/view/generate"
	"net/http"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	images := make([]types.Image, 20)
	data := generate.ViewData{
		Images: images,
	}
	return render(w, r, generate.Index(data))
}
