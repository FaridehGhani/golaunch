package httpext

import (
	"github.com/mohsensamiei/golaunch/errorext"
	"net/url"
)

// URL is wrap on net/url for easy marshaling
type URL struct {
	*url.URL
}

// NewURL create new instance of URL from net/url
func NewURL(url *url.URL) *URL {
	return &URL{
		URL: url,
	}
}

// ParseURL parses raw into a URL structure
func ParseURL(raw string) (*URL, error) {
	urlAddress, err := url.Parse(raw)
	if err != nil {
		return nil, errorext.NewValidationError("invalid url string", err)
	}
	return NewURL(urlAddress), nil
}

// MarshalYAML marshal url to yaml
func (url URL) MarshalYAML() (interface{}, error) {
	return url.URL.String(), nil
}

// UnmarshalYAML yaml to url
func (url *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return errorext.NewValidationError("invalid url yml", err)
	}
	parsedURL, err := url.Parse(str)
	if err != nil {
		return errorext.NewValidationError("invalid url string", err)
	}
	url.URL = parsedURL
	return nil
}

// MarshalText convert url to text
func (url URL) MarshalText() (text []byte, err error) {
	return []byte(url.URL.String()), nil
}

// UnmarshalText convert text to url
func (url *URL) UnmarshalText(text []byte) error {
	var err error
	url.URL, err = url.Parse(string(text))
	if err != nil {
		return errorext.NewValidationError("invalid url text", err)
	}
	return nil
}
