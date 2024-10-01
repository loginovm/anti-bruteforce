package app

func NewAppError(msg string) *Error {
	return &Error{err: msg}
}

type Error struct {
	err string
}

func (e Error) Error() string {
	return e.err
}
