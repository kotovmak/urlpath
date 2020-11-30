package urlpath

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	var expected = "" +
		"/S3_key/S3_value/S4_key/S4_value" +
		"/I2_key/0/I3_key/0/I4_key/4" +
		"/U2_key/0/U3_key/0/U4_key/4" +
		"/B2_key/false/B3_key/false/B4_key/true" +
		"/Base64_key/QmFzZTY0X3ZhbHVl"

	var example = struct {
		// string
		S  string `urlpath:""`                        // without tag
		S1 string `urlpath:"-"`                       // explicit ignored
		S2 string `urlpath:"S2_key"`                  // zero-value
		S3 string `urlpath:"S3_key;default=S3_value"` // zero-value with default option
		S4 string `urlpath:"S4_key"`                  // non-zero-value
		S5 string `urlpath:"S5_key;omitempty"`        // zero-value with omitempty
		// int
		I  int    `urlpath:""`                 // without tag
		I1 int    `urlpath:"-"`                // explicit ignored
		I2 int    `urlpath:"I2_key"`           // zero-value
		I3 int    `urlpath:"I3_key;default=3"` // zero-value with default option
		I4 int    `urlpath:"I4_key"`           // non-zero-value
		I5 string `urlpath:"I5_key;omitempty"` // zero-value with omitempty
		// uint
		U  uint   `urlpath:""`                 // without tag
		U1 uint   `urlpath:"-"`                // explicit ignored
		U2 uint   `urlpath:"U2_key"`           // zero-value
		U3 uint   `urlpath:"U3_key;default=3"` // zero-value with default option
		U4 uint   `urlpath:"U4_key"`           // non-zero-value
		U5 string `urlpath:"U5_key;omitempty"` // zero-value with omitempty
		// bool
		B  bool   `urlpath:""`                    // without tag
		B1 bool   `urlpath:"-"`                   // explicit ignored
		B2 bool   `urlpath:"B2_key"`              // zero-value
		B3 bool   `urlpath:"B3_key;default=true"` // zero-value with default option
		B4 bool   `urlpath:"B4_key"`              // non-zero-value
		B5 string `urlpath:"B5_key;omitempty"`    // zero-value with omitempty
		// custom
		Base64 Base64 `urlpath:"Base64_key"`
	}{
		S: "S_value", S1: "S1_value", S4: "S4_value",
		I: 12345, I1: 1, I4: 4,
		U: 12345, U1: 1, U4: 4,
		B: true, B1: true, B4: true,
		Base64: "Base64_value",
	}

	got, err := Marshal(&example)
	if err != nil {
		t.Errorf("Marshal failed: %v", err)
	}

	if got != expected {
		t.Errorf("expected - %s, got - %s", expected, got)
	}
}

func TestUnmasrshal(t *testing.T) {
	t.Run("Unexpected url format", func(t *testing.T) {
		err := Unmarshal("some", &struct{}{})
		t.Logf("Unmarshal - %v", err)
		if err == nil {
			t.Error("Unmarshal unexpected error absence")
		}

		var e Error
		if !errors.As(err, &e) {
			t.Error("Unmarshal unexpected type of error")
		}

		if e.ErrorType != InvalidFormatError {
			t.Error("Unmarshal unexpected type of error")
		}

		if e.Error() != "urlpath: odd number of elements" {
			t.Error("Unmarshal unexpected error message")
		}
	})

	t.Run("Required field missed", func(t *testing.T) {
		err := Unmarshal("/key/value", &struct {
			Field string `urlpath:"test;required"`
		}{})
		t.Logf("Unmarshal - %v", err)
		if err == nil {
			t.Error("Unmarshal unexpected error absence")
		}

		var e Error
		if !errors.As(err, &e) {
			t.Error("Unmarshal unexpected type of error")
		}

		if e.ErrorType != InvalidFormatError {
			t.Error("Unmarshal unexpected type of error")
		}

		if e.Error() != "urlpath: required key test is missing" {
			t.Error("Unmarshal unexpected error message")
		}
	})

	t.Run("Unexpected value type", func(t *testing.T) {
		var v = struct {
			Key int `urlpath:"key"`
		}{}

		err := Unmarshal("/key/value", &v)
		t.Logf("Unmarshal - %v", err)
		if err == nil {
			t.Error("Unmarshal unexpected error absence")
		}

		var e Error
		if !errors.As(err, &e) {
			t.Error("Unmarshal unexpected type of error")
		}

		if e.ErrorType != InvalidFormatError {
			t.Error("Unmarshal unexpected type of error")
		}

		if !strings.HasPrefix(e.Error(), "urlpath: decode value of field key failed") {
			t.Error("Unmarshal unexpected error message")
		}
	})

	t.Run("All values struct", func(t *testing.T) {
		var example = "S2_key/S2_value" +
			"/I2_key/1/I3_key/0" +
			"/U2_key/1/U3_key/0" +
			"/B2_key/true/B3_key/false" +
			"/Base64_key/QmFzZTY0X3ZhbHVl"

		type model struct {
			// string
			S  string `urlpath:"S_key"`                   // not presented field
			S1 string `urlpath:"S1_key;default=S1_value"` // not presented field with default
			S2 string `urlpath:"S2_key"`                  // presented field
			// int
			I  int `urlpath:"I_key"`            // not presented field
			I1 int `urlpath:"I1_key;default=1"` // not presented field with default
			I2 int `urlpath:"I2_key"`           // presented field
			I3 int `urlpath:"I3_key"`           // presented field with zero-value
			// uint
			U  int `urlpath:"U_key"`            // not presented field
			U1 int `urlpath:"U1_key;default=1"` // not presented field with default
			U2 int `urlpath:"U2_key"`           // presented field
			U3 int `urlpath:"U3_key"`           // presented field with zero-value
			// bool
			B  bool `urlpath:"B_key"`               // not presented field
			B1 bool `urlpath:"B1_key;default=true"` // not presented field with default
			B2 bool `urlpath:"B2_key"`              // presented field
			B3 bool `urlpath:"B3_key"`              // presented field with zero-value
			// custom
			Base64 Base64 `urlpath:"Base64_key"`
		}

		var expected = model{
			S1: "S1_value", S2: "S2_value",
			I1: 1, I2: 1,
			U1: 1, U2: 1,
			B1: true, B2: true,
			Base64: "Base64_value",
		}

		var got model
		err := Unmarshal(example, &got)
		if err != nil {
			t.Errorf("Unmarshal failed: %v", err)
		}

		if !reflect.DeepEqual(&got, &expected) {
			t.Errorf("got - %#v, expected - %#v", got, expected)
		}
	})
}
