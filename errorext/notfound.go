package errorext

import "net/http"

const NotFoundMessage = "data not found"

type NotFoundError struct {
	handledError
}

func (NotFoundError) StatusCode() int {
	return http.StatusNotFound
}
func NewNotFound(errs ...error) IHandledError {
	return NewNotFoundError(NotFoundMessage, errs...)
}
func NewNotFoundError(message string, errs ...error) IHandledError {
	if ok, handled := hasHandled(errs...); ok {
		return handled
	}
	return NotFoundError{newHandled(message, errs...)}
}
