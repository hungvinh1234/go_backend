package util

// MyError .
type MyError struct {
	Raw       error
	ErrorCode int
	HTTPCode  int
	Message   string
}

func (e MyError) Error() string {
	return e.Message
}

// NewError .
func NewError(err error, httpCode, errCode int, message string) MyError {
	return MyError{
		Raw:       err,
		ErrorCode: errCode,
		HTTPCode:  httpCode,
		Message:   message,
	}
}
