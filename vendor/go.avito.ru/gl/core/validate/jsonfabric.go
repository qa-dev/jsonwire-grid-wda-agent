package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "go.avito.ru/gl/core/logger"
	"go.avito.ru/gl/core/validate/jsonvalidator"
)

//ValCreateInfo хранит информацию о созданных валидаторах
type ValCreateInfo struct {
	Val Validator
	Err error
}

//JsonFabric фабрика json валидаторов
type JsonFabric struct {
	validators map[string]ValCreateInfo
	path       string
}

//NewJsonFabric инициализирует фабрику путём до схем
func NewJsonFabric(path string) *JsonFabric {
	path = strings.TrimPrefix(path, "./")
	f := &JsonFabric{path: path}
	f.validators = make(map[string]ValCreateInfo)
	filepath.Walk(f.path, f.loadSchemes)

	return f
}

//GetValidator возвращает валидатор по алиасу
// алиасом является относительный путь до схемы без .json
// Например если схемы лежат /etc/jsonschemas/ -этот путь передается в конструкторе
// /etc/jsonschemas
// |-send_request.json
// |-blablabla
// |--foobar.json
//то алиасами будут send_request и blablabla/foobar
func (f *JsonFabric) GetValidator(alias string) (Validator, error) {
	if v, ok := f.validators[alias]; ok == true {
		return v.Val, v.Err
	}
	return nil, fmt.Errorf("unknown alias:%v", alias)
}

func (f *JsonFabric) loadSchemes(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return nil
	}

	if filepath.Ext(path) == ".json" {

		alias := strings.Replace(path, f.path, "", -1)
		alias = strings.TrimLeft(alias, "/")
		alias = strings.Replace(alias, ".json", "", -1)

		log.Infof("'%v'", alias)

		validator, err := jsonvalidator.NewValidatorByFile(path)

		if err != nil {
			log.Infof("validator creation failed! %v path:%s", err, path)
			f.validators[alias] = ValCreateInfo{nil, err}
			return nil //continue to Walk
		}

		f.validators[alias] = ValCreateInfo{validator, nil}
	}
	return nil
}
