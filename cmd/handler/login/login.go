package login

import (
	"net/http"
	"strings"

	"github.com/arisudesu/go-admin/cmd/handler"
)

func NewHandler(
	html *handler.HtmlHandler,
	resolveUrl func(name string, args ...string) string,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if handler.IsAuthenticated(r) {
				http.Redirect(w, r, resolveUrl("index"), http.StatusFound)
				return
			}

			var username, password string

			if r.Method == http.MethodPost {
				username = r.PostFormValue("username")
				password = r.PostFormValue("password")

				if username == "admin" && password == "12345" {
					cookie := &http.Cookie{
						Name:  "session",
						Value: "1",
						Path:  resolveUrl("index"),
					}
					http.SetCookie(w, cookie)

					redir := r.URL.Query().Get("to")
					if redir == "" || !strings.HasPrefix(redir, "/") {
						redir = resolveUrl("index")
					}

					http.Redirect(w, r, redir, http.StatusFound)
					return
				}
			}

			tmplData := handler.HtmlCtx{
				"Username": username,
			}
			html.Success(w, r, "login.gohtml", tmplData)
		},
	)
}
