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

		switch valueType := value.(type) {
		case map[string]interface{}:
			if field.Kind() == reflect.Struct {
				// recurse
				field.Set(reflect.ValueOf(MapToStruct(valueType, field.Addr().Interface())))
			} else if field.Kind() == reflect.Map {
				// convert map
				mapValue := reflect.MakeMap(field.Type())
				for key, value := range valueType {
					mapValue.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
				}
				field.Set(mapValue)
			}
		case []interface{}:
			// convert slice
			if field.Kind() == reflect.Slice {
				sliceLen := len(valueType)
				slice := reflect.MakeSlice(field.Type(), 0, sliceLen)
				for i := 0; i < sliceLen; i++ {
					slice = reflect.Append(slice, reflect.ValueOf(valueType[i]))
				}
				field.Set(slice)
			}
		default:
			// try to set field
			if reflect.TypeOf(value).AssignableTo(field.Type()) {
				field.Set(reflect.ValueOf(value))
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

func MergeDeep(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	for key, value := range b {
		switch value.(type) {
		case map[string]interface{}:
			if _, ok := a[key]; !ok {
				// a does not have this key
				a[key] = value
			} else if _, ok := a[key].(map[string]interface{}); ok {
				// a has this key and is is and map too
				a[key] = MergeDeep(a[key].(map[string]interface{}), value.(map[string]interface{}))
			} else {
				// a has this key but is not a map - replace this value
				a[key] = value
			}
		case []interface{}:
			if _, ok := a[key]; !ok {
				// a does not have this key
				a[key] = value
			} else if _, ok := a[key].([]interface{}); ok {
				// a has this key and is is and slice too
				a[key] = append(a[key].([]interface{}), value.([]interface{})...)
			} else {
				// a has this key but is not a slice - replace this value
				a[key] = value
			}
		default:
			a[key] = value
		}
	}
	return a
}
