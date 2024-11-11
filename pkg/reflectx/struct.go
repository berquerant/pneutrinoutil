package reflectx

import (
	"errors"
	"reflect"
)

type Field struct {
	Name string
	Type reflect.Type
	Tag  reflect.StructTag
}

var (
	ErrNotStruct = errors.New("NotStruct")
)

func NewFields(v any) ([]Field, error) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	fields := make([]Field, t.NumField())
	for i := range t.NumField() {
		f := t.Field(i)
		fields[i] = Field{
			Name: f.Name,
			Type: f.Type,
			Tag:  f.Tag,
		}
	}

	return fields, nil
}
