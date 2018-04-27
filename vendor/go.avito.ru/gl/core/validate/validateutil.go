package validate

import (
	"encoding/json"
	"io"

	log "go.avito.ru/gl/core/logger"
)

//ValidateTo валидирует объект io.ReadCloser (чаще всего это может быть http.Request.Body или http.Response. Body)
//указанным валидатором и пытается с unmarshalить результат
func ValidateTo(r io.ReadCloser, v Validator, to interface{}) error {
	b, err := v.ValidateReader(r)

	if err != nil {
		log.Errorf("Error while validate err:%v data:%s", err, b)
		return err
	}

	err = json.Unmarshal(b, &to)
	if err != nil {
		log.Errorf("Error while unmarshal data err:%v", err)
		return err
	}

	return nil
}
