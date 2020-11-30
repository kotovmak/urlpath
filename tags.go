package urlpath

import (
	"reflect"
	"strings"
)

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
