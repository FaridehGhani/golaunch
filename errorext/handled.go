package errorext

import (
	"fmt"
	"strings"
)

type IHandledError interface {
	Error() string
	Inners() []error
	String() string
	StatusCode() int
}

type handledError struct {
	Message string `json:"message"`
	Errors  []error `json:"-"`
}

func (err handledError) Error() string {
	return err.Message
}
func (err handledError) Inners() []error {
	return err.Errors
}
func (err handledError) String() string {
	return fmt.Sprint(err.Message, ": ", strings.Join(ToStringSlice(err.Errors...), ", "))
}
func newHandled(message string, errs ...error) handledError {
	return handledError{
		Message: message,
		Errors:  errs,
	}
}
func hasHandled(errs ...error) (bool, IHandledError) {
	for _, err := range errs {
		if handled, ok := err.(IHandledError); ok {
			return true, handled
		}
	}
	return false, nil
}
