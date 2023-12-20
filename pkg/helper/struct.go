package helper

import (
	"reflect"
)

func MapToStruct[T any](data map[string]interface{}, target T) T {
	targetValue := reflect.ValueOf(&target)

	for key, value := range data {
		field := FindJsonField(targetValue, key)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		switch value.(type) {
		case map[string]interface{}:
			// recurse
			panic("not implemented")
		default:
			// try to set field
			if field.IsValid() && field.CanSet() {
				if reflect.TypeOf(value).AssignableTo(field.Type()) {
					field.Set(reflect.ValueOf(value))
				}
			}
		}
	}
	return target
}

func FindJsonField(value reflect.Value, key string) reflect.Value {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)

		if fieldType.Tag.Get("json") == key {
			return field
		}
	}

	return reflect.Value{}

}
