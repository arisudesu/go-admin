package logout

import (
	"net/http"

	"github.com/arisudesu/go-admin/cmd/handler"
)

func NewHandler(
	html *handler.HtmlHandler,
	resolveUrl func(name string, args ...string) string,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookie := &http.Cookie{
				Name:   "session",
				Value:  "",
				Path:   resolveUrl("index"),
				MaxAge: -1,
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		},
	)
}
