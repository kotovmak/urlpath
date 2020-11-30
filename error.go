package urlpath

var _ error = Error{}

type Error struct {
	ErrorType
	error
}

func (e Error) Error() string {
	return "urlpath: " + e.error.Error()
}

type ErrorType = uint8

const InvalidFormatError ErrorType = iota

func newError(t ErrorType, err error) Error { return Error{t, err} }
