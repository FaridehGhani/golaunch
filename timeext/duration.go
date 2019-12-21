package timeext

import (
	"encoding/json"
	"github.com/mohsensamiei/golaunch/errorext"
	"time"
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
	if err != nil {
		return NewDuration(0), errorext.NewValidationError("invalid duration string", err)
	}
	return NewDuration(duration), nil
}

// MarshalYAML convert duration to yaml
func (duration Duration) MarshalYAML() (interface{}, error) {
	return duration.String(), nil
}

// UnmarshalYAML convert yaml to duration
func (duration *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err:= unmarshal(&duration.Duration); err!=nil{
		return errorext.NewValidationError("invalid duration yaml", err)
	}
	return nil
}

// MarshalJSON is implement for converting duration to json string
func (duration Duration) MarshalJSON() ([]byte, error) {
	result, err:= json.Marshal(duration.String())
	if err != nil {
		return nil, errorext.NewInternalError("parse duration to json failed", err)
	}
	return result, nil
}

// UnmarshalJSON is implement for json string to duration
func (duration *Duration) UnmarshalJSON(bytes []byte) error {
	var object interface{}
	if err := json.Unmarshal(bytes, &object); err != nil {
		return errorext.NewValidationError("invalid duration json", err)
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
			return errorext.NewValidationError("invalid duration json", err)
		}
		return nil
	default:
		return errorext.NewValidationError("invalid duration parsing type")
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
	if err != nil {
		return errorext.NewValidationError("invalid duration text", err)
	}
	return nil
}
