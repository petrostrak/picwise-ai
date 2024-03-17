package handler

import (
	"net/http"

	"github.com/nedpals/supabase-go"
	"github.com/petrostrak/picwise-ai/view/auth"
)

func HandleSignInIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Login())
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	return render(w, r, auth.LoginForm(credentials, auth.LoginErrors{
		InvalidCredentials: "The credentials you provided are invalid",
	}))
}
