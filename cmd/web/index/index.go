package index

import (
	"net/http"

	"github.com/arisudesu/go-admin/cmd/web"
)

func Handler(html *web.HtmlHandler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			html.Success(w, r, "index.gohtml", nil)
		},
	)
}
