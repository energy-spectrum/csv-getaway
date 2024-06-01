package error_templates

import (
	"github.com/pkg/errors"
)

type OutputError struct {
	errorMessage   string
	httpStatusCode int
}

func (e *OutputError) Error() string {
	return e.errorMessage
}

func New(msg string, httpCode int) *OutputError {
	return &OutputError{
		errorMessage:   msg,
		httpStatusCode: httpCode,
	}
}

func (e *OutputError) GetHTTP() (int, error) {
	return e.httpStatusCode, errors.New(e.errorMessage)
}
