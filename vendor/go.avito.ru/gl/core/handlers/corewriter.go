package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.avito.ru/gl/core/validate/jsonvalidator"
	"go.avito.ru/gl/core/web"
)

//CoreWriter обёртка над http.ResponseWriter для форматированния сообщение в core протокол
type CoreWriter struct {
	http.ResponseWriter
	code int
}

//CoreWriter инициализируется http.ResponseWriter
func NewCoreWriter(w http.ResponseWriter) *CoreWriter {
	return &CoreWriter{
		ResponseWriter: w,
		code:           200,
	}
}

//Header возвращает хидер http.ResponseWriter'a которым был проинициализирован
func (cr *CoreWriter) Header() http.Header {
	return cr.ResponseWriter.Header()
}

//Write Если код отличается от 200 то сообщение оборачивается в нащ протокол
func (cr *CoreWriter) Write(b []byte) (int, error) {
	if cr.code != 200 {
		errStr := string(b)
		errStr = strings.Replace(errStr, "\n", "", -1)

		return cr.WriteJson(web.ErrorResult{
			Error: web.ErrorData{
				Code:    cr.code,
				Message: errStr,
			}}, cr.code)
	} else {
		cr.ResponseWriter.WriteHeader(cr.code)
		return cr.ResponseWriter.Write(b)
	}
}

//WriteJson Повторяет функционал web.JSONResponse
func (cr *CoreWriter) WriteJson(data interface{}, code int) (int, error) {
	cr.Header().Set("Content-Type", "application/json; charset=utf-8")
	cr.ResponseWriter.WriteHeader(code)

	body, err := json.Marshal(data)
	if err != nil {
		body = []byte(
			`{"error": "Unknown and unpredictable error with huge, massive and catastrophic consequences!"}`)
	}

	return cr.ResponseWriter.Write(body)
}

//ValidationError обычно вызывается из хэндлеров с валидацией
//Если ошибка jsonvalidator.ValidationError то добавляется поле schema с расширенным кодом ошибки
func (cr *CoreWriter) ValidationError(err error) {
	switch err.(type) {
	case *jsonvalidator.ValidationError:
		e := err.(*jsonvalidator.ValidationError)
		if e != nil {
			cr.WriteJson(web.ErrorResult{
				Error: web.ErrorData{
					Code:    cr.code,
					Message: err.Error(),
					Scheme:  e.ErrorsToCode(),
				}}, http.StatusBadRequest)

			return
		}
	}
	http.Error(cr, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

//WriteHeader...
func (cr *CoreWriter) WriteHeader(code int) {
	cr.code = code
}
