package url

import (
	"net/url"
)

// URL is wrap on net/url for easy marshaling
type URL struct {
	*url.URL
}

// New create new instance of URL from net/url
func New(url *url.URL) *URL {
	return &URL{
		URL: url,
	}
}

// Parse parses rawurl into a URL structure
func Parse(rawurl string) (*URL, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return New(url), nil
}

// MarshalYAML marshal url to yaml
func (url URL) MarshalYAML() (interface{}, error) {
	return url.String(), nil
}

// UnmarshalYAML yaml to url
func (url *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buffer []byte
	if err := unmarshal(&buffer); err != nil {
		return err
	}
	if err := url.UnmarshalBinary(buffer); err != nil {
		return err
	}
	return nil
}
