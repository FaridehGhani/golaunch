package restful

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"
)

const (
	// ContentType key in header
	ContentType = "Content-Type"
	// XMLContentType xml value of content type
	XMLContentType = "application/xml"
	// JSONContentType json value of content type
	JSONContentType = "application/json"
)

var (
	errorInvalidContentType   = errors.New("Invalid request content type")
	errorInvalidAuthorization = errors.New("Invalid bearer authorization header")
)

// SendResponse by model and status code
func SendResponse(response http.ResponseWriter, request *http.Request, statusCode int, model interface{}) {
	var bytes []byte
	if model != nil {
		switch request.Header.Get(ContentType) {
		case XMLContentType:
			response.Header().Set(ContentType, XMLContentType)
			bytes, _ = ConvertToBytes(model, xml.Marshal)
		case JSONContentType, "":
			response.Header().Set(ContentType, JSONContentType)
			bytes, _ = ConvertToBytes(model, json.Marshal)
		}
	}
	response.WriteHeader(statusCode)
	response.Write(bytes)
}

// ConvertToBytes convert model to response byte array
func ConvertToBytes(model interface{}, convert func(interface{}) ([]byte, error)) ([]byte, error) {
	switch value := model.(type) {
	case []byte:
		return value, nil
	case error:
		model = NewErrorDTO(value)
	}
	return convert(model)
}

// ParseRequestBody parsing request body to model
func ParseRequestBody(request *http.Request, model interface{}) error {
	switch request.Header.Get(ContentType) {
	case XMLContentType:
		return xml.NewDecoder(request.Body).Decode(&model)
	case JSONContentType:
		return json.NewDecoder(request.Body).Decode(&model)
	default:
		return errorInvalidContentType
	}
}

// GetAuthBearerHeader get Bearer from header
func GetAuthBearerHeader(request *http.Request) (string, error) {
	authorization := request.Header.Get("Authorization")
	dump := strings.Split(authorization, "Bearer")
	if len(dump) != 2 {
		return "", errorInvalidAuthorization
	}
	return strings.TrimSpace(dump[1]), nil
}

// IsSuccessStatusCode return true if http status code is in success
func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
