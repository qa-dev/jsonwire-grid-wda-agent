package handlers

import (
	"net/http"

	"go.avito.ru/gl/core/context"
	"go.avito.ru/gl/core/web"
)

// ErrorHandler задает обработчик, который всегда возвращает 500.
type ErrorHandler struct {
	ctx *context.BaseContext
}

// NewErrorHandler возвраащет новый ErrorHandler.
func NewErrorHandler(ctx *context.BaseContext) *ErrorHandler {
	return &ErrorHandler{
		ctx: ctx,
	}
}

// ServeHTTP обрабатывает HTTP-запросы.
func (h *ErrorHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	result := web.ErrorResult{
		Error: web.ErrorData{
			Code:    http.StatusInternalServerError,
			Message: "What we're dealing with here is a total lack of respect for the law.",
		}}

	web.JSONResponse(resp, result, http.StatusInternalServerError)
}
