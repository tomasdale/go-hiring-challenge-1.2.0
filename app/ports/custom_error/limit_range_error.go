package custom_error

type Error struct {
	code int
}

func New(message string, code int) *Error {
	return &Error{code: code}
}

func (e *Error) Error() string {
	return "Limit must be between 1 and 100"
}

func (e *Error) Code() int {
	return e.code
}
