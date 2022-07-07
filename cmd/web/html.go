package web

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

const (
	tplError404 = "error404.gohtml"
	tplError500 = "error500.gohtml"

	contentTypeHeader = "Content-Type"
	contentTypeHtml   = "text/html; charset=utf-8"
)

type HtmlHandler struct {
	template   *template.Template
	processors []htmlCtxProc
}

type HtmlCtxProcFunc func(*http.Request, HtmlCtx)

type htmlCtxProc interface {
	Process(r *http.Request, ctx HtmlCtx)
}

func (fn HtmlCtxProcFunc) Process(r *http.Request, ctx HtmlCtx) {
	fn(r, ctx)
}

type HtmlCtx map[string]any

func NewHtmlHandler(template *template.Template, prf ...HtmlCtxProcFunc) *HtmlHandler {
	h := &HtmlHandler{
		template: template,
	}
	for _, fn := range prf {
		h.processors = append(h.processors, fn)
	}
	return h
}

func (h *HtmlHandler) Success(w http.ResponseWriter, r *http.Request, template string, ctx HtmlCtx) {
	if ctx == nil {
		ctx = HtmlCtx{}
	}
	for _, fn := range h.processors {
		fn.Process(r, ctx)
	}

	buf, err := h.render(template, ctx)
	if err != nil {
		log.Printf("Error rendering %s: %+v", template, errors.WithStack(err))
		h.error500(w)
		return
	}

	w.Header().Set(contentTypeHeader, contentTypeHtml)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(buf); err != nil {
		log.Printf("Error writing response: %+v", errors.WithStack(err))
	}
}

func (h *HtmlHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	ctx := HtmlCtx{}
	for _, fn := range h.processors {
		fn.Process(r, ctx)
	}

	buf, err := h.render(tplError404, ctx)
	if err != nil {
		log.Printf("Error rendering %s: %+v", tplError404, errors.WithStack(err))
		h.error500(w)
		return
	}

	w.Header().Set(contentTypeHeader, contentTypeHtml)
	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write(buf); err != nil {
		log.Printf("Error writing response: %+v", errors.WithStack(err))
	}
}

func (h *HtmlHandler) Error(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Handler error: %+v", err)
	h.error500(w)
}

func (h *HtmlHandler) render(template string, ctx HtmlCtx) ([]byte, error) {
	var buf bytes.Buffer
	if err := h.template.ExecuteTemplate(&buf, template, ctx); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}

func (h *HtmlHandler) error500(w http.ResponseWriter) {
	buf, err := h.render(tplError500, nil)
	if err != nil {
		log.Printf("Error rendering %s: %+v", tplError500, errors.WithStack(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeHeader, contentTypeHtml)
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write(buf); err != nil {
		log.Printf("Error writing response: %+v", errors.WithStack(err))
	}
}

type NotFoundHandler struct {
	html *HtmlHandler
}

func NewNotFoundHandler(html *HtmlHandler) *NotFoundHandler {
	return &NotFoundHandler{
		html: html,
	}
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.html.NotFound(w, r)
}
