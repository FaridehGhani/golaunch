package types

import (
	"net/url"
)

// Address is wrap on url for easy marshaling
type Address struct {
	url.URL
}

// NewAddress create new Address from url
func NewAddress(url url.URL) Address {
	return Address{
		URL: url,
	}
}

// MarshalYAML marshal url to yaml
func (address Address) MarshalYAML() (interface{}, error) {
	return address.String(), nil
}

// UnmarshalYAML yaml to url
func (address *Address) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buffer []byte
	if err := unmarshal(&buffer); err != nil {
		return err
	}
	if err := address.UnmarshalBinary(buffer); err != nil {
		return err
	}
	return nil
}
