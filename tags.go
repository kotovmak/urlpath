package urlpath

import (
	"reflect"
	"strings"
)

type field struct {
	tags
	reflect.Value
}

func parseFields(v reflect.Value) []field {
	elems := []field{}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct && v.Type().Field(i).Anonymous {
			elems = append(elems, parseFields(v.Field(i))...)
			continue
		}

		tags := parseTag(v.Type().Field(i))
		if tags.ignore || (v.Field(i).IsZero() && tags.omitempty) {
			continue
		}

		elems = append(elems, field{tags, v.Field(i)})
	}
	return elems
}

type tags struct {
	ignore       bool
	required     bool
	omitempty    bool
	name         string
	defaultValue string
}

func parseTag(field reflect.StructField) (t tags) {
	value, exists := field.Tag.Lookup("urlpath")
	if !exists || value == "-" || value == "" {
		t.ignore = true
		return
	}

	keys := strings.Split(value, ";")
	for i := range keys {
		switch {
		case i == 0:
			t.name = keys[i]
		case keys[i] == "required":
			t.required = true
		case keys[i] == "omitempty":
			t.omitempty = true
		case strings.HasPrefix(keys[i], "default="):
			t.defaultValue = strings.TrimPrefix(keys[i], "default=")
		}
	}
	return
}
