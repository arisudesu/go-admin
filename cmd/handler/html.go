package handler

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/pkg/errors"
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

func NewHtmlHandler(template *template.Template) *HtmlHandler {
	return &HtmlHandler{
		template: template,
	}
}

func (h *HtmlHandler) Use(prf ...HtmlCtxProcFunc) {
	for _, fn := range prf {
		h.processors = append(h.processors, fn)
	}
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
		h.error500(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

	buf, err := h.render("error404.gohtml", ctx)
	if err != nil {
		log.Printf("Error rendering error404.gohtml: %+v", errors.WithStack(err))
		h.error500(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write(buf); err != nil {
		log.Printf("Error writing response: %+v", errors.WithStack(err))
	}
}

func (h *HtmlHandler) Error(w http.ResponseWriter, r *http.Request, err error) {
	h.error500(w, r)
}

func (h *HtmlHandler) render(template string, ctx HtmlCtx) ([]byte, error) {
	var buf bytes.Buffer
	if err := h.template.ExecuteTemplate(&buf, template, ctx); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}

func (h *HtmlHandler) error500(w http.ResponseWriter, r *http.Request) {
	buf, err := h.render("error500.gohtml", nil)
	if err != nil {
		log.Printf("Error rendering error500.gohtml: %+v", errors.WithStack(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(500)
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
