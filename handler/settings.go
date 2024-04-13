package handler

import (
	"github.com/petrostrak/picwise-ai/db"
	"github.com/petrostrak/picwise-ai/pkg/kit/validate"
	"github.com/petrostrak/picwise-ai/view/settings"
	"net/http"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(w, r, settings.Index(user))
}

func HandleSettingsUsernameUpdate(w http.ResponseWriter, r *http.Request) error {
	params := settings.ProfileParams{
		Username: r.FormValue("username"),
	}
	var errors settings.ProfileErrors
	if ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(3), validate.Max(40)),
	}).Validate(&errors); !ok {
		return render(w, r, settings.ProfileForm(params, errors))
	}
	user := getAuthenticatedUser(r)
	user.Username = params.Username
	if err := db.UpdateAccount(&user.Account); err != nil {
		return err
	}
	params.Success = true
	return render(w, r, settings.ProfileForm(params, settings.ProfileErrors{}))
}
