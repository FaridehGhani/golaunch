package errorext

import "net/http"

const ValidationMessage = "invalid data or parameters"

type ValidationError struct {
	handledError
}

func (ValidationError) StatusCode() int {
	return http.StatusBadRequest
}
func NewValidation(errs ...error) IHandledError {
	return NewValidationError(ValidationMessage, errs...)
}
func NewValidationError(message string, errs ...error) IHandledError {
	if ok, handled := hasHandled(errs...); ok {
		return handled
	}
	return ValidationError{newHandled(message, errs...)}
}
