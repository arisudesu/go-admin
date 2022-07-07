package web

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func NewRequireAuthMiddleware(
	urlgen *URLGenerator,
	skipFunc func(r *http.Request) bool,
) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if skipFunc(r) || IsAuthenticated(r) {
					next.ServeHTTP(w, r)
					return
				}

				redirQuery := url.Values{}
				redirQuery.Set("next", r.URL.String())

				redir := urlgen.GenerateURL("login")
				redir.RawQuery = redirQuery.Encode()

				http.Redirect(w, r, redir.String(), http.StatusFound)
				return
			},
		)
	}
}

const cookieName = "auth"

func LogIn(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: "1",
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func LogOut(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func IsAuthenticated(r *http.Request) bool {
	_, err := r.Cookie(cookieName)
	return err != http.ErrNoCookie
}
