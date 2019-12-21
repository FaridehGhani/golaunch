package errorext

import (
	"log"
	"net/http"
	"strings"
)

const InternalMessage = "internal server error"

type InternalError struct {
	handledError
}

func (InternalError) StatusCode() int {
	return http.StatusInternalServerError
}
func (err InternalError) String() string {
	return err.Message
}
func NewInternal(errs ...error) IHandledError {
	return NewInternalError(InternalMessage, errs...)
}
func NewInternalError(message string, errs ...error) IHandledError {
	log.Print(message, ": ", strings.Join(ToStringSlice(errs...), ", "))
	if ok, handled := hasHandled(errs...); ok {
		return handled
	}
	return InternalError{newHandled(message, errs...)}
}
