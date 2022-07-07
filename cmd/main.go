package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/arisudesu/go-admin/cmd/web"
	"github.com/arisudesu/go-admin/cmd/web/index"
	"github.com/arisudesu/go-admin/cmd/web/login"
	"github.com/arisudesu/go-admin/cmd/web/logout"
)

func main() {
	// Initialize router and url generator very early.
	// We need them available in template functions.

	router := mux.NewRouter()
	urlgen := web.NewURLGenerator(router)

	// Initialize templates, providing function to generate urls from the router.

	tmplFuncs := template.FuncMap{
		"url": urlgen.Generate,
	}

	templates := template.New("").Funcs(tmplFuncs)

	// Load templates recursively using our recursive loader.

	if err := web.LoadTemplates(templates, "templates", "*.gohtml"); err != nil {
		log.Fatal(err)
	}

	// Initialize html handler that is able to render responses from these templates.
	// Provide to it context processors which extend render context with useful info.

	ctxAddVersion := func(r *http.Request, ctx web.HtmlCtx) { ctx["Version"] = version }
	ctxAddRequest := func(r *http.Request, ctx web.HtmlCtx) { ctx["Request"] = r }

	htmlHandler := web.NewHtmlHandler(
		templates,
		ctxAddVersion,
		ctxAddRequest,
	)

	// Setup application middlewares and route handles.

	loggingMw := web.NewLoggingMiddleware()
	reqAuthMw := web.NewRequireAuthMiddleware(
		urlgen, func(r *http.Request) bool { return r.URL.Path == urlgen.Generate("login") },
	)

	router.Use(loggingMw)
	router.Use(reqAuthMw)

	router.Handle("/", index.Handler(htmlHandler)).Name("index")
	router.Handle("/login/", login.Handler(htmlHandler, urlgen)).Name("login")
	router.Handle("/logout/", logout.Handler(urlgen)).Name("logout")

	// Register NotFoundHandler after all other handlers, see issue:
	// https://github.com/gorilla/mux/issues/416#issuecomment-600074549.
	// This forces router to execute middlewares before NotFoundHandler.

	router.NotFoundHandler = router.NewRoute().Handler(
		web.NewNotFoundHandler(htmlHandler)).GetHandler()

	// Build std mux, mount page router, assets and favicon to it.
	// This is the lowest level of request processing chain in app.

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
