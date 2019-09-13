package types

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	errorInvalidInterval = errors.New("Invalid interval parsing type")
)

// Interval is a wrap of time.duration for supporting json Marshal and Unmarshal
type Interval struct {
	time.Duration
}

// NewInterval create new instance of interval
func NewInterval(duration time.Duration) Interval {
	return Interval{Duration: duration}
}

// MarshalYAML convert interval to yaml
func (interval Interval) MarshalYAML() (interface{}, error) {
	return interval.String(), nil
}

// UnmarshalYAML convert yaml to interval
func (interval *Interval) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&interval.Duration)
}

// MarshalJSON is implement for converting duration to json string
func (interval Interval) MarshalJSON() ([]byte, error) {
	return json.Marshal(interval.String())
}

// UnmarshalJSON is implement for json string to duration
func (interval *Interval) UnmarshalJSON(bytes []byte) error {
	var object interface{}
	if err := json.Unmarshal(bytes, &object); err != nil {
		return err
	}
	switch value := object.(type) {
	case int64:
		interval.Duration = time.Duration(value)
		return nil
	case float64:
		interval.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		if interval.Duration, err = time.ParseDuration(value); err != nil {
			return err
		}
		return nil
	default:
		return errorInvalidInterval
	}
}
