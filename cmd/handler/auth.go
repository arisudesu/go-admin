package handler

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func NewRequireAuthMiddleware(
	skipFunc func(r *http.Request) bool,
	resolveUrl func(name string, args ...string) string,
) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if skipFunc(r) || IsAuthenticated(r) {
					next.ServeHTTP(w, r)
					return
				}

				queryStr := url.Values{}
				queryStr.Set("to", r.URL.String())

				redir := resolveUrl("login") + "?" + queryStr.Encode()
				http.Redirect(w, r, redir, http.StatusFound)
				return
			},
		)
	}
}

func IsAuthenticated(r *http.Request) bool {
	_, err := r.Cookie("session")
	return err != http.ErrNoCookie
}
