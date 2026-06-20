package optional

import (
	"encoding/json"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type Optional[T any] struct {
	Set   bool
	Value T
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Set = true
	return json.Unmarshal(data, &o.Value)
}

func (o Optional[T]) Schema(r huma.Registry) *huma.Schema {
	s := huma.SchemaFromType(r, reflect.TypeOf(o.Value))
	s.Nullable = isNilable(reflect.TypeOf(o.Value))
	return s
}

func (o Optional[T]) IsSet() bool {
	return o.Set
}

func (o Optional[T]) Or(defaultValue T) T {
	if !o.Set {
		return defaultValue
	}

	return o.Value
}

func isNilable(t reflect.Type) bool {
	if t == nil {
		return false
	}

	switch t.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return true
	default:
		return false
	}
}
