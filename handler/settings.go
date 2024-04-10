package handler

import (
	"github.com/petrostrak/picwise-ai/view/settings"
	"net/http"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(w, r, settings.Index(user))
}
