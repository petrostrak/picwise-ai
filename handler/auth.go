package handler

import (
	"log/slog"
	"net/http"

	"github.com/nedpals/supabase-go"
	"github.com/petrostrak/picwise-ai/pkg/sb"
	"github.com/petrostrak/picwise-ai/pkg/sb/util"
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

	if !util.IsValidEmail(credentials.Email) {
		return render(w, r, auth.LoginForm(credentials, auth.LoginErrors{
			Email: "Please enter a valid email",
		}))
	}

	if reason, ok := util.ValidatePassword(credentials.Password); !ok {
		return render(w, r, auth.LoginForm(credentials, auth.LoginErrors{
			Password: reason,
		}))
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(w, r, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials you provided are invalid",
		}))
	}

	cookie := &http.Cookie{
		Value:    resp.AccessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
