package urlpath

import (
	"reflect"
	"strings"
)

/*
	Tag options:
	1) Field int `urlpath:"-"` - expilicit ignore field for marshal/unmarshal
	2) Field int `urlpath:"name"` - appears in path as key `name`
	3) Field int `urlpath:"name;required"` -
	4) Field int `urlpath:"name;omitempty"` -
	5) Field int `urlpath:"name;default=12345` -
*/

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
