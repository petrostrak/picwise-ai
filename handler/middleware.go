package handler

import (
	"net/http"
	"strings"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		println("from the WithUser middleware")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
