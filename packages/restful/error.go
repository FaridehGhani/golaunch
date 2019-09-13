package restful

// ErrorDTO dto of error message
type ErrorDTO struct {
	Message string `json:"message"`
}

// NewErrorDTO create new instance of ErrorDTO
func NewErrorDTO(err error) ErrorDTO {
	return ErrorDTO{
		Message: err.Error(),
	}
}
