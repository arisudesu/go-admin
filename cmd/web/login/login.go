package login

import (
	"net/http"
	"strings"

	"github.com/arisudesu/go-admin/cmd/web"
)

func Handler(html *web.HtmlHandler, urlgen *web.URLGenerator) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if web.IsAuthenticated(r) {
				http.Redirect(w, r, urlgen.Generate("index"), http.StatusFound)
				return
			}

			var username, password, err string

			if r.Method == http.MethodPost {
				username = r.PostFormValue("username")
				password = r.PostFormValue("password")

				if username == "admin" && password == "12345" {
					web.LogIn(w)

					redir := r.URL.Query().Get("next")
					if redir == "" || !strings.HasPrefix(redir, "/") || strings.HasPrefix(redir, "//") {
						redir = urlgen.Generate("index")
					}

					http.Redirect(w, r, redir, http.StatusFound)
					return
				}

				err = "Неправильное имя пользователя или пароль."
			}

			tmplData := web.HtmlCtx{
				"Username": username,
				"Error":    err,
			}
			html.Success(w, r, "login.gohtml", tmplData)
		},
	)
}
