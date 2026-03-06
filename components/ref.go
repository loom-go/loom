package components

import (
	"fmt"
	"reflect"
)

// Ref assigns the parent node to the provided pointer when applied.
type Ref struct {
	Ptr any // *T
	Fn  any // func(T)
}

func (r Ref) Apply(parent any) (func() error, error) {
	// todo: should be able to do all this without reflection and with compile-time type safety
	// once this is added: https://github.com/golang/go/issues/61731

	parentType := reflect.TypeOf(parent)

	if r.Ptr != nil {
		ptrType := reflect.TypeOf(r.Ptr)
		if ptrType.Kind() != reflect.Pointer || ptrType.Elem() != parentType {
			return nil, fmt.Errorf("Ref: %w: the given Ptr type (%s) does not match the parent node type (%s)", ErrNodeRefMissMatch, ptrType, parentType)
		}

		reflect.ValueOf(r.Ptr).Elem().Set(reflect.ValueOf(parent))
	}

	if r.Fn != nil {
		fnType := reflect.TypeOf(r.Fn)
		if fnType.Kind() != reflect.Func || fnType.NumIn() != 1 || fnType.In(0) != parentType {
			return nil, fmt.Errorf("Ref: %w: the given Fn type (%s) does not match the parent node type (%s)", ErrNodeRefMissMatch, fnType, parentType)
		}

		reflect.ValueOf(r.Fn).Call([]reflect.Value{reflect.ValueOf(parent)})
	}

	return nil, nil
}
