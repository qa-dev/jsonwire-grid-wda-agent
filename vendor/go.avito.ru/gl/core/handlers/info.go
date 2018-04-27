package handlers

import (
	"net/http"

	"go.avito.ru/gl/core/context"
	"go.avito.ru/gl/core/web"
)

// InfoResultData описывает структуру данных, возвращаемых в поле result ответа.
type InfoResultData struct {
	OK bool `json:"ok"`
}

// InfoResult описывает структуру данных ответа.
type InfoResult struct {
	Result InfoResultData `json:"result"`
}

func newInfoResult(ok bool) InfoResult {
	return InfoResult{
		Result: InfoResultData{
			OK: ok,
		},
	}
}

// InfoHandler обрабатывает запросы на получение общей информации о приложении.
type InfoHandler struct {
	ctx *context.BaseContext
}

// NewInfoHandler возвраащет новый InfoHandler.
func NewInfoHandler(ctx *context.BaseContext) *InfoHandler {
	return &InfoHandler{
		ctx: ctx,
	}
}

// ServeHTTP обрабатывает HTTP-запросы.
func (h *InfoHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	result := newInfoResult(true)

	web.JSONResponse(resp, result, http.StatusOK)
}
