package handlers

import (
	"net/http"

	"go.avito.ru/gl/core/context"
	"go.avito.ru/gl/core/web"
)

// DefaultHandler обрабатывает запросы по-умолчанию.
type DefaultHandler struct {
	ctx *context.BaseContext
}

// NewDefaultHandler возвраащет новый DefaultHandler.
func NewDefaultHandler(ctx *context.BaseContext) *DefaultHandler {
	return &DefaultHandler{
		ctx: ctx,
	}
}

// ServeHTTP обрабатывает HTTP-запросы.
func (h *DefaultHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	web.JSONResponse(resp, web.ErrNotFound, http.StatusNotFound)
}
