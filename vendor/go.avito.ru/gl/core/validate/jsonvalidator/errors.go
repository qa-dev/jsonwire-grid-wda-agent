package jsonvalidator

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//ValidationError класс обёртка на gojsonschema.ResultError
type ValidationError struct {
	Errors []gojsonschema.ResultError
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

//NewValidationError инициализируем ошибку валидации массивом ошибок из библиотеки xeipuuv/gojsonschema
func NewValidationError(errs []gojsonschema.ResultError) *ValidationError {
	return &ValidationError{errs}
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

//ErrorsToCode переводит ошибки gojsonschema.ResultError в коды core протокола
//https://cf.avito.ru/display/BD/Core+Services#CoreServices-Ответ
func (ve *ValidationError) ErrorsToCode() map[string]int {
	result := make(map[string]int)
	for _, err := range ve.Errors {
		code := 0

		switch err.Type() {
		default:
		case "required":
			code = 1201
		case "invalid_type":
			if expected, ok := err.Details()["expected"]; ok == true {
				switch expected {
				case "integer":
					code = 1202
				case "number":
					code = 1203
				case "email":
					code = 1204
				case "uri":
					code = 1206
				case "string":
					code = 1705
				case "array":
					code = 1706
				case "boolean":
					code = 1707
				}
			}
		case "array_max_items":
			code = 1218
		case "array_min_items":
			code = 1222
		case "unique":
			code = 1225
		case "multiple_of":
			code = 1226
		case "number_gt":
			code = 1208
		case "number_gte":
			code = 1223
		case "string_gte":
			code = 1208
		case "number_lt":
			code = 1224
		case "number_lte":
			code = 1209
		case "string_lte":
			code = 1209
		case "EnumError":
			code = 1214
		}

		result[err.Field()] = code
	}
	return result
}

//Error форматированное сообщение об ошибке
func (ve *ValidationError) Error() string {
	r := []string{}
	for _, e := range ve.Errors {
		r = append(r, fmt.Sprintf("error: field %v %v", e.Field(), e.Description()))
	}
	return strings.Join(r, "\r\n")
}
