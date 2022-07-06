package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type catchStatusWriter struct {
	http.ResponseWriter
	status int
}

func (w *catchStatusWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func NewLoggerMiddleware() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				catch := catchStatusWriter{ResponseWriter: w}

				defer func() {
					log.Printf("%s %s %s ==> %d %v ", r.RemoteAddr, r.Method, r.URL, catch.status, time.Now().Sub(start))
				}()

				h.ServeHTTP(&catch, r)
			},
		)
	}
}
