package validate

import (
	"io"
	"io/ioutil"

	log "go.avito.ru/gl/core/logger"
)

//Careless валидатор который беспечно валидирует всё подряд
type Careless struct {
}

func NewCareless() *Careless {
	return &Careless{}
}

func (careless *Careless) ValidateReader(reader io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(reader)
	defer reader.Close()
	if err != nil {
		return []byte{}, err
	}
	return data, careless.ValidateBytes(data)
}

func (careless *Careless) ValidateBytes(data []byte) error {
	log.Infof("Validation? Who cares about validation!? body:%v", string(data))
	return nil
}
