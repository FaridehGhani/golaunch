package convert

import (
	"errors"
	"reflect"
	"strings"
)

var (
	errorDTONotPointer         = errors.New("DTO is not a pointer")
	errorDTONotKindModel       = errors.New("DTO is not kind of model")
	errorModelNotStructOrSlice = errors.New("Model is not struct or slice of struct")
)

// MapToDTO convert nested model to flat DTO
func MapToDTO(model interface{}, dto interface{}) error {
	modelReflect := reflect.ValueOf(model)
	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	if modelReflect.Kind() != reflect.Struct &&
		!(modelReflect.Kind() == reflect.Ptr && modelReflect.Type().Elem().Kind() == reflect.Struct) &&
		!(modelReflect.Kind() == reflect.Slice && modelReflect.Type().Elem().Kind() == reflect.Struct) &&
		!(modelReflect.Kind() == reflect.Slice && modelReflect.Type().Elem().Kind() == reflect.Ptr && modelReflect.Type().Elem().Elem().Kind() == reflect.Struct) {
		return errorModelNotStructOrSlice
	}

	dtoPointer := reflect.ValueOf(dto)
	if dtoPointer.Kind() != reflect.Ptr {
		return errorDTONotPointer
	}
	dtoReflect := dtoPointer.Elem()
	if dtoReflect.Kind() != modelReflect.Kind() {
		return errorDTONotKindModel
	}

	dtoType := dtoReflect.Type()
	if dtoType.Kind() == reflect.Slice {
		dtoType = dtoType.Elem()
	}

	if modelReflect.Kind() == reflect.Slice {
		for index := 0; index < modelReflect.Len(); index++ {
			instance := reflect.New(dtoType).Elem()
			fillDTOFields(modelReflect.Index(index), instance, dtoType)
			dtoReflect.Set(reflect.Append(dtoReflect, instance))
		}
	} else {
		fillDTOFields(modelReflect, dtoReflect, dtoType)
	}

	return nil
}

func fillDTOFields(modelReflect reflect.Value, dtoReflect reflect.Value, dtoType reflect.Type) {
	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}
	for index := 0; index < dtoType.NumField(); index++ {
		dtoTypeField := dtoType.Field(index)

		dtoField := dtoReflect.FieldByName(dtoTypeField.Name)

		var modelField reflect.Value
		if name, ok := dtoTypeField.Tag.Lookup("map"); ok {
			modelField = getNestedField(modelReflect, strings.Split(name, "."), 0)
		} else {
			modelField = modelReflect.FieldByName(dtoTypeField.Name)
		}

		if modelField.IsValid() && modelField.Type() == dtoField.Type() {
			dtoField.Set(modelField)
		}
	}
}

func getNestedField(model reflect.Value, nestedNames []string, nameIndex int) reflect.Value {
	if model.IsZero() {
		return model
	} else if model.Kind() == reflect.Ptr {
		model = model.Elem()
	}

	field := model.FieldByName(nestedNames[nameIndex])
	if len(nestedNames) == nameIndex+1 {
		return field
	}

	if field.Type().Kind() == reflect.Slice {
		var values reflect.Value
		for fieldIndex := 0; fieldIndex < field.Len(); fieldIndex++ {
			value := getNestedField(field.Index(fieldIndex), nestedNames, nameIndex+1)
			if !values.IsValid() {
				values = reflect.MakeSlice(reflect.SliceOf(value.Type()), 0, field.Len())
			}
			values = reflect.Append(values, value)
		}
		return values
	}

	return getNestedField(field, nestedNames, nameIndex+1)
}
