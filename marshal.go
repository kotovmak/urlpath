package urlpath

import (
	"fmt"
	"reflect"
	"strings"
)

type Marshaler interface {
	MarshalURLValue() (value string, err error)
}

func Marshal(v interface{}) (string, error) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return encode(value)
}

func encode(v reflect.Value) (string, error) {
	elems := make([]string, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		tag := parseTag(v.Type().Field(i))
		if tag.ignore {
			continue
		}

		value, err := encodeField(v.Field(i))
		if err != nil {
			return "", fmt.Errorf("encode field %s failed: %v", tag.name, err)
		}

		if value == "" {
			value = tag.defaultValue
		}

		elems = append(elems, tag.name)
		elems = append(elems, value)
	}
	return "/" + strings.Join(elems, "/"), nil
}

func encodeField(v reflect.Value) (string, error) {
	marshaler, ok := v.Interface().(Marshaler)
	if ok {
		return marshaler.MarshalURLValue()
	}
	return fmt.Sprintf("%v", v), nil
}
