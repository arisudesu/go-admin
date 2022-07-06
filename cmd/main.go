package main

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/arisudesu/go-admin/cmd/handler"
	"github.com/arisudesu/go-admin/cmd/handler/index"
	"github.com/arisudesu/go-admin/cmd/handler/login"
	"github.com/arisudesu/go-admin/cmd/handler/logout"
)

const Version = "0.1.0"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

func main() {
	// Initialize templates and html handler

	var resolver interface{ Get(string) *mux.Route }

	resolveUrl := func(name string, args ...string) string {
		url, err := resolver.Get(name).URL(args...)
		if err != nil {
			panic(err)
		}
		return url.String()
	}

	tmplFuncs := template.FuncMap{
		"url": resolveUrl,
	}

	tmpl, err := template.
		New("__templates__").Funcs(tmplFuncs).
		ParseGlob(path.Join("templates", "*.gohtml"))
	if err != nil {
		log.Fatal(errors.WithStack(err))
	}

	ctxAddVersion := func(r *http.Request, ctx handler.HtmlCtx) { ctx["Version"] = Version }

	htmlHandler := handler.NewHtmlHandler(tmpl)
	htmlHandler.Use(ctxAddVersion)

	// Build front router and middleware

	loggerMw := handler.NewLoggerMiddleware()

	reqAuthMw := handler.NewRequireAuthMiddleware(
		func(r *http.Request) bool { return r.URL.Path == resolveUrl("login") },
		resolveUrl,
	)

	router := mux.NewRouter()
	router.Use(loggerMw, reqAuthMw)

	// Register routes wrapped with middleware

	root := router.PathPrefix("/").Subrouter()
	root.NotFoundHandler = handler.NewNotFoundHandler(htmlHandler)

	root.Handle("/", index.NewHandler(htmlHandler)).Name("index")
	root.Handle("/login/", login.NewHandler(htmlHandler, resolveUrl)).Name("login")
	root.Handle("/logout/", logout.NewHandler(htmlHandler, resolveUrl)).Name("logout")

	resolver = root

	// Build std mux with front router and assets

	assetsHandler := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/", router)
	mux.Handle("/assets/", http.StripPrefix("/assets/", assetsHandler))
	mux.Handle("/favicon.ico", assetsHandler)

	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = mux
	log.Println(server.ListenAndServe())
}
