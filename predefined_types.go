package urlpath

import (
	"encoding/base64"
	"net/url"
)

var (
	_ Marshaler   = Base64("")
	_ Unmarshaler = new(Base64)
)

type Base64 string

func (b Base64) MarshalURLValue() (string, error) {
	return base64.URLEncoding.EncodeToString([]byte(b)), nil
}

func (b *Base64) UnmarshalURLValue(value string) error {
	decoded, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return err
	}
	*b = Base64(decoded)
	return nil
}

var (
	_ Marshaler   = URLEscaped("")
	_ Unmarshaler = new(URLEscaped)
)

type URLEscaped string

func (u URLEscaped) MarshalURLValue() (string, error) {
	return url.QueryEscape(string(u)), nil
}

func (u *URLEscaped) UnmarshalURLValue(value string) error {
	s, err := url.QueryUnescape(value)
	if err != nil {
		return err
	}
	*u = URLEscaped(s)
	return nil
}
