package urlpath

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Unmarshaler interface {
	UnmarshalURLValue(value string) error
}

func Unmarshal(urlpath string, out interface{}) error {
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return errors.New("value for unmarshaling is not a pointer")
	}
	if outValue.IsNil() {
		return errors.New("value for unmarshaling is not a nil pointer")
	}
	outValue = outValue.Elem()
	if outValue.Kind() != reflect.Struct {
		return errors.New("value for unmarshaling is not of struct kind")
	}

	args, err := parse(strings.TrimPrefix(urlpath, "/"))
	if err != nil {
		return err
	}
	return decode(args, outValue)
}

func parse(s string) (map[string]string, error) {
	tokens := strings.Split(s, "/")
	if len(tokens)%2 != 0 {
		return nil, newError(InvalidFormatError, errors.New("odd number of elements"))
	}

	var kv = make(map[string]string, len(tokens)/2)
	for i := 0; i < len(tokens); i += 2 {
		var (
			key   = tokens[i]
			value = tokens[i+1]
		)
		if kv[key] != "" {
			return nil, newError(InvalidFormatError, fmt.Errorf("duplicate declaration of key %s", key))
		}
		kv[key] = value
	}
	return kv, nil
}

func decode(args map[string]string, v reflect.Value) error {
	fields := parseFields(v)
	for _, field := range fields {
		value, exists := args[field.tags.name]
		if !exists && field.tags.required {
			return newError(InvalidFormatError, fmt.Errorf("required key %s is missing", field.tags.name))
		}

		if value == "" {
			value = field.tags.defaultValue
			if value == "" {
				continue
			}
		}
		if len(field.tags.gt) > 0 {
			floatValue, _ := strconv.ParseFloat(value, 64)
			floatGT, _ := strconv.ParseFloat(field.tags.gt, 64)
			if floatValue <= floatGT {
				return newError(InvalidFormatError, fmt.Errorf("key %s must be biger than %s", field.tags.name, field.tags.gt))
			}
		}

		err := decodeField(field.Value, value)
		if err != nil {
			return newError(InvalidFormatError, fmt.Errorf("decode value of field %s failed: %v", field.tags.name, err))
		}
	}
	return nil
}

func decodeField(v reflect.Value, value string) error {
	unmarshaler, ok := v.Interface().(Unmarshaler)
	if ok {
		return unmarshaler.UnmarshalURLValue(value)
	}

	if v.CanAddr() {
		unmarshaler, ok = v.Addr().Interface().(Unmarshaler)
		if ok {
			return unmarshaler.UnmarshalURLValue(value)
		}
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(parsed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(parsed)
	case reflect.Float32:
		parsed, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		v.SetFloat(parsed)
	case reflect.Float64:
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		v.SetFloat(parsed)
	}
	return nil
}
