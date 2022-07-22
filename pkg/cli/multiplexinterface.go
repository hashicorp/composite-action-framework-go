package cli

import (
	"reflect"
)

type collector[T any] interface {
	Collect(T)
}

func multiplexInterface[T any, C collector[T]](maybe any, collection C) T {
	tType := reflect.TypeOf(new(T)).Elem()
	out, ok := any(collection).(T)
	if !ok {
		panic("collection must implement T")
	}
	t := reflect.TypeOf(maybe)
	if t.Kind() != reflect.Pointer {
		return out
	}
	v := reflect.ValueOf(maybe).Elem()
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return out
	}

	// Collect the struct fields in the order specified in code.
	for i := 0; i < t.NumField(); i++ {
		fieldVal := v.Field(i)
		if fieldVal.Kind() == reflect.Pointer {
			fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		} else {
			fieldVal = fieldVal.Addr()
		}
		if fieldVal.Type().AssignableTo(tType) {
			collection.Collect(fieldVal.Interface().(T))
		}
	}

	// If maybe is also a T, add it to the collection as well.
	if f, ok := maybe.(T); ok {
		collection.Collect(f)
	}

	return out
}
