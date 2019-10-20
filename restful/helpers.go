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
	switch value := model.(type) {
	case error:
		model = NewErrorDTO(value)
	}

	var bytes []byte
	switch value := model.(type) {
	case nil:
	case []byte:
		bytes = value
		if response.Header().Get(ContentType) == "" {
			response.Header().Set(ContentType, JSONContentType)
		}
	default:
		switch request.Header.Get(ContentType) {
		case XMLContentType:
			response.Header().Set(ContentType, XMLContentType)
			bytes, _ = xml.Marshal(model)
		case JSONContentType, "":
			response.Header().Set(ContentType, JSONContentType)
			bytes, _ = json.Marshal(model)
		}
	}

	response.WriteHeader(statusCode)
	response.Write(bytes)
}

// AddResponseHeaders add header list to response
func AddResponseHeaders(response http.ResponseWriter, headers http.Header) {
	for key, value := range headers {
		response.Header()[key] = value
	}
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
