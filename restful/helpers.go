package restful

import (
	"encoding/json"
	"encoding/xml"
	"github.com/mohsensamiei/golaunch/errorext"
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

// SendStandard by model
func SendStandard(response http.ResponseWriter, request *http.Request, model interface{}) {
	var status int
	switch value := model.(type) {
	case nil:
		status = http.StatusNoContent
	case errorext.IHandledError:
		status = value.StatusCode()
	case error:
		status = http.StatusInternalServerError
	default:
		if request.Method == http.MethodPost {
			status = http.StatusCreated
		} else {
			status = http.StatusOK
		}
	}
	Send(response, request, status, model)
}

// Send by model and status code
func Send(response http.ResponseWriter, request *http.Request, statusCode int, model interface{}) {
	switch value := model.(type) {
	case errorext.IHandledError:
	case error:
		model = errorext.NewInternal(value)
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
	_, _ = response.Write(bytes)
}

// AddHeaders add header list to response
func AddHeaders(response http.ResponseWriter, headers http.Header) {
	for key, value := range headers {
		response.Header()[key] = value
	}
}

// ParseBody parsing request body to model
func ParseBody(request *http.Request, model interface{}) error {
	switch request.Header.Get(ContentType) {
	case XMLContentType:
		if err := xml.NewDecoder(request.Body).Decode(&model); err != nil {
			return errorext.NewValidationError("invalid request xml body", err)
		}
	case JSONContentType:
		if err := json.NewDecoder(request.Body).Decode(&model); err != nil {
			return errorext.NewValidationError("invalid request json body", err)
		}
	default:
		return errorext.NewValidationError("invalid request content type")
	}
	return nil
}

// GetBearerToken get bearer token from header
func GetBearerToken(request *http.Request) (string, error) {
	authorization := request.Header.Get("Authorization")
	dump := strings.Split(authorization, "Bearer")
	if len(dump) != 2 {
		return "", errorext.NewAuthorizeError("invalid authorization bearer token")
	}
	return strings.TrimSpace(dump[1]), nil
}

// IsSuccess return true if http status code is in success
func IsSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
