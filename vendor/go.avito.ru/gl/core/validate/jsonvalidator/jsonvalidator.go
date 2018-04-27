package jsonvalidator

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/xeipuuv/gojsonschema"
)

//Validator валидирует json по json schema
type Validator struct {
	Schema *gojsonschema.Schema
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

//NewValidatorByString создаёт валидатор из строки
func NewValidatorByString(schema string) (*Validator, error) {
	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schema))
	if err != nil {
		return nil, err
	}

	return &Validator{
		Schema: s,
	}, nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

//NewValidatorByFile создаёт валидатор из файла
func NewValidatorByFile(filename string) (*Validator, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(string(b)))
	if err != nil {
		return nil, err
	}

	return &Validator{
		Schema: s,
	}, nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

//ValidateReader валидирует объект io.ReadCloser, возвращает вычитанные байты и ошибку
func (validator *Validator) ValidateReader(reader io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(reader)
	defer reader.Close()
	if err != nil {
		return []byte{}, err
	}
	return data, validator.ValidateBytes(data)
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

//ValidateBytes валидируются raw bytes
func (validator *Validator) ValidateBytes(data []byte) error {
	result, err := validator.Schema.Validate(gojsonschema.NewStringLoader(string(data)))
	if err != nil {
		return err
	}
	valid := result.Valid()
	if valid == true {
		return nil
	}

	if len(result.Errors()) > 0 {
		return NewValidationError(result.Errors())
	}
	return errors.New("Unknown error")
}
