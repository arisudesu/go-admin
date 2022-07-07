package web

import (
	"net/url"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var ErrNoRoute = errors.New("no route with provided name")

type URLGenerator struct {
	router *mux.Router
}

func NewURLGenerator(router *mux.Router) *URLGenerator {
	return &URLGenerator{
		router: router,
	}
}

func (r *URLGenerator) Generate(name string, args ...string) string {
	return r.GenerateURL(name, args...).String()
}

func (r *URLGenerator) GenerateURL(name string, args ...string) *url.URL {
	u, err := r.maybeGenerateURL(name, args...)
	if err != nil {
		panic(err)
	}
	return u
}

func (r *URLGenerator) maybeGenerateURL(name string, args ...string) (*url.URL, error) {
	route := r.router.Get(name)
	if route == nil {
		return nil, ErrNoRoute
	}

	u, err := route.URL(args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return u, nil
}
