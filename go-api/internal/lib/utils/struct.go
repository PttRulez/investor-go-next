package utils

import "reflect"

func StructForEach(s any, fn func(key string, value any)) {
	v := reflect.ValueOf(s)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fn(t.Field(i).Name, v.Field(i).Interface())
	}
}

func IsNilPointer(v any) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}
