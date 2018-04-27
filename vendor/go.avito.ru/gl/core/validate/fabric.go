package validate

import (
	"os"

	"go.avito.ru/gl/core/config"
	log "go.avito.ru/gl/core/logger"
)

//Fabric Настраиваемая фабрика, которая выдаёт валидатор по алиасу(имя файла без .json)
//на данный момент может выдавать либо json валидаторы, либо пустышки
//путь до схем задаётся или в конфиге или в параметрах конструктора
type Fabric struct {
	jsonFabric  *JsonFabric
	dummyFabric *DummyFabric
	useDummy    bool
}

//GetValidator возвращет валидатор и ошибку(если валидатора нет или при его создании она была)
func (f *Fabric) GetValidator(alias string) (Validator, error) {
	if f.useDummy == false {
		return f.jsonFabric.GetValidator(alias)
	}
	return f.dummyFabric.GetValidator(alias)
}

//SetUseDummy возвращать ли только валидаторы пустышки
func (f *Fabric) SetUseDummy(use bool) {
	f.useDummy = use
}

//NewValidateFabric создаётся  новая фабрика которая берет путь до схем из конфига
func NewValidateFabric(useDummy bool) *Fabric {
	cfg := &config.BaseConfig{}
	err := config.LoadFromFile(os.Getenv("CONFIG_PATH"), cfg)

	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	log.Infof("Try to create validator for jsonschema path='%s'", cfg.Validator.Path)

	return &Fabric{
		jsonFabric:  NewJsonFabric(cfg.Validator.Path),
		dummyFabric: NewDummyFabric(),
		useDummy:    useDummy,
	}
}

//NewValidateFabricByPath создаётся  новая фабрика которая берет путь до схем из пареметра
func NewValidateFabricByPath(path string, useDummy bool) *Fabric {
	log.Infof("Try to create validator for jsonschema path='%s'", path)
	return &Fabric{
		jsonFabric:  NewJsonFabric(path),
		dummyFabric: NewDummyFabric(),
		useDummy:    useDummy,
	}
}
