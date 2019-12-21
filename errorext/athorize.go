package errorext

import "net/http"

const AuthorizeMessage = "unauthorized request"

type AuthorizeError struct {
	handledError
}

func (AuthorizeError) StatusCode() int {
	return http.StatusUnauthorized
}
func NewAuthorize(errs ...error) IHandledError {
	return NewAuthorizeError(AuthorizeMessage, errs...)
}
func NewAuthorizeError(message string, errs ...error) IHandledError {
	if ok, handled := hasHandled(errs...); ok {
		return handled
	}
	return AuthorizeError{newHandled(message, errs...)}
}
