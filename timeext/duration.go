package timeext

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	errorInvalidDuration = errors.New("Invalid duration parsing type")
)

// Duration is a wrap of time.duration for supporting Marshal and Unmarshal
type Duration struct {
	time.Duration
}

// NewDuration create new instance of duration
func NewDuration(duration time.Duration) Duration {
	return Duration{Duration: duration}
}

// ParseDuration parses a duration string
func ParseDuration(s string) (Duration, error) {
	duration, err := time.ParseDuration(s)
	return NewDuration(duration), err
}

// MarshalYAML convert duration to yaml
func (duration Duration) MarshalYAML() (interface{}, error) {
	return duration.String(), nil
}

// UnmarshalYAML convert yaml to duration
func (duration *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&duration.Duration)
}

// MarshalJSON is implement for converting duration to json string
func (duration Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(duration.String())
}

// UnmarshalJSON is implement for json string to duration
func (duration *Duration) UnmarshalJSON(bytes []byte) error {
	var object interface{}
	if err := json.Unmarshal(bytes, &object); err != nil {
		return err
	}
	switch value := object.(type) {
	case int64:
		duration.Duration = time.Duration(value)
		return nil
	case float64:
		duration.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		if duration.Duration, err = time.ParseDuration(value); err != nil {
			return err
		}
		return nil
	default:
		return errorInvalidDuration
	}
}

// MarshalText convert duration to text
func (duration Duration) MarshalText() (text []byte, err error) {
	return []byte(duration.String()), nil
}

// UnmarshalText convert text to duration
func (duration *Duration) UnmarshalText(text []byte) error {
	var err error
	duration.Duration, err = time.ParseDuration(string(text))
	return err
}
