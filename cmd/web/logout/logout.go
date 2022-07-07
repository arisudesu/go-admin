package logout

import (
	"net/http"

	"github.com/arisudesu/go-admin/cmd/web"
)

func Handler(urlgen *web.URLGenerator) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			web.LogOut(w)

			http.Redirect(w, r, urlgen.Generate("index"), http.StatusFound)
		},
	)
}
