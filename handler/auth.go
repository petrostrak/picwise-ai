package handler

import (
	"github.com/petrostrak/picwise-ai/pkg/kit/validate"
	"log/slog"
	"net/http"

	"github.com/nedpals/supabase-go"
	"github.com/petrostrak/picwise-ai/pkg/sb"
	"github.com/petrostrak/picwise-ai/view/auth"
)

func HandleSignInIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Login())
}

func HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Signup())
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	errors := auth.SignupErrors{}

	if ok := validate.New(&params, validate.Fields{
		"Email":    validate.Rules(validate.Email),
		"Password": validate.Rules(validate.Password),
		"ConfirmPassword": validate.Rules(
			validate.Equal(params.Password),
			validate.Message("passwords do not match"),
		),
	}).Validate(&errors); !ok {
		return render(w, r, auth.SignupForm(params, errors))
	}

	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}

	return render(w, r, auth.SignupSuccess(user.Email))
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(w, r, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials you provided are invalid",
		}))
	}

	setAuthCookie(w, resp.AccessToken)
	return hxRedirect(w, r, "/")
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(w, r, auth.CallbaclScript())
	}
	setAuthCookie(w, accessToken)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Value:    "",
		Name:     "at",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func setAuthCookie(w http.ResponseWriter, accessToken string) {
	cookie := &http.Cookie{
		Value:    accessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
}

func HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "http://localhost:3000/auth/callback", // TODO: Read url from env var
	})
	if err != nil {
		return err
	}
	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
	return nil
}
