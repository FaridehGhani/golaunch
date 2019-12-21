package utils

import (
	"fmt"
	"github.com/mohsensamiei/golaunch/errorext"
)

// GetMapField return value of map field
func GetMapField(vars map[string]string, field string, parser func(string) (interface{}, error)) (interface{}, error) {
	value, ok := vars[field]
	if ok == false {
		return nil, errorext.NewValidationError(fmt.Sprint(field, " field does not exist"))
	}
	result, err := parser(value)
	if err != nil {
		return nil, errorext.NewValidationError(fmt.Sprint(field, " field parsing failed"), err)
	}
	return result, nil
}
