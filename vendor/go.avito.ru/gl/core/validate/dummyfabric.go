package validate

//DummyFabric фабрика которая создаёт только careless валидаторов
type DummyFabric struct {
	Validator
}

func NewDummyFabric() *DummyFabric {
	return &DummyFabric{NewCareless()}
}

func (f *DummyFabric) GetValidator(alias string) (Validator, error) {
	return f.Validator, nil
}
