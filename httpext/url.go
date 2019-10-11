package httpext

import (
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

// ParseURL parses rawurl into a URL structure
func ParseURL(rawurl string) (*URL, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return NewURL(url), nil
}

// MarshalYAML marshal url to yaml
func (url URL) MarshalYAML() (interface{}, error) {
	return url.URL.String(), nil
}

// UnmarshalYAML yaml to url
func (url *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	parsedURL, err := url.Parse(str)
	url.URL = parsedURL
	return err
}

// MarshalText convert url to text
func (url URL) MarshalText() (text []byte, err error) {
	return []byte(url.URL.String()), nil
}

// UnmarshalText convert text to url
func (url *URL) UnmarshalText(text []byte) error {
	var err error
	url.URL, err = url.Parse(string(text))
	return err
}
