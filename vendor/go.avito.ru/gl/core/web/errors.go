package web

import (
	"net/http"
)

// ErrorData описывает формат данных об ошибке для core протокола.
type ErrorData struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Scheme  map[string]int `json:"scheme,omitempty"`
}

// ErrorResult описывает общую структуру ответа об ошибке для core протокола.
type ErrorResult struct {
	Error ErrorData `json:"error"`
}

var (
	// ErrBadRequest описывает ответ об ошибке валидации.
	ErrBadRequest = ErrorResult{
		Error: ErrorData{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
		}}
	// ErrNotFound описывает ответ об ошибке когда запись не найдена.
	ErrNotFound = ErrorResult{
		Error: ErrorData{
			Code:    http.StatusNotFound,
			Message: "Handler not found.",
		}}
	// ErrNotAllowed описывает ответ об ошибке доступа.
	ErrNotAllowed = ErrorResult{
		Error: ErrorData{
			Code:    http.StatusMethodNotAllowed,
			Message: "Method is not allowed.",
		}}
	// ErrInternal описывает ответ о внутренней ошибке сервера..
	ErrInternal = ErrorResult{
		Error: ErrorData{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error.",
		}}
)
