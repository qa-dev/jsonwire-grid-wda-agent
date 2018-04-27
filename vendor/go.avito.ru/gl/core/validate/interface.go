package validate

import (
	"io"
)

//Validator объект валидатор должен поодерживать возможность валидации io.ReadCloser(например http.Request.Body)
// И валидацию набора байтов
type Validator interface {
	ValidateReader(reader io.ReadCloser) ([]byte, error)
	ValidateBytes(data []byte) error
}

//ValidateFabric возвращает валидатор по алиасу
type ValidateFabric interface {
	GetValidator(alias string) (Validator, error)
}
