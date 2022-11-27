package golidator

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidatorFunc func() *ValueError
type Validator struct {
	Field    string
	Function ValidatorFunc
}

type ValidatorCollection []*Validator

type IValidators interface {
	GetValidators(params ...interface{}) ValidatorCollection
}

func Validate(obj IValidators, ctx ...interface{}) *ValidationError {
	validationError := &ValidationError{
		source: obj,
	}
	validatorCollection := obj.GetValidators(ctx...)
	for _, validator := range validatorCollection {
		if err := validator.Function(); err != nil {
			validationError.Errors = append(validationError.Errors, &FieldError{
				ValueError: *err,
				Location:   validator.Field,
			})
		}
	}

	if len(validationError.Errors) > 0 {
		return validationError
	}
	return nil
}

func GetValidatorsForList[T IValidators](fieldName string, list []T, ctx ...interface{}) ValidatorCollection {
	validators := ValidatorCollection{}
	for i, children := range list {
		for _, validator := range children.GetValidators(ctx...) {
			childFieldName := strings.Join([]string{fieldName, fmt.Sprint(i), validator.Field}, ".")
			newValidator := &Validator{Field: childFieldName, Function: validator.Function}
			validators = append(validators, newValidator)
		}
	}
	return validators
}

func GetValidatorsForObject[T IValidators](fieldName string, child T, ctx ...interface{}) ValidatorCollection {
	if reflect.ValueOf(child).IsNil() {
		return nil
	}

	validators := ValidatorCollection{}
	for _, validator := range child.GetValidators(ctx...) {
		childFieldName := strings.Join([]string{fieldName, validator.Field}, ".")
		newValidator := &Validator{Field: childFieldName, Function: validator.Function}
		validators = append(validators, newValidator)
	}
	return validators
}
