package golidator

import (
	"fmt"
	"reflect"
	"strings"
)

type ErrorType string

type ValidationError struct {
	source interface{}
	Errors []*FieldError
}

func (err *ValidationError) Error() string {
	type_ := reflect.ValueOf(err.source).Elem().Type()
	mainErrorString := fmt.Sprintf("%d validation error(s) for %v", len(err.Errors), type_)
	errorStrings := []string{mainErrorString}

	for _, fieldError := range err.Errors {
		errorStrings = append(errorStrings, fieldError.Error())
	}

	return strings.Join(errorStrings, "\n")
}

type FieldError struct {
	ValueError
	Location string
}

func (err *FieldError) Error() string {
	return err.Location + ": " + string(err.ErrorType)
}

type ValueError struct {
	Context   map[string]any
	ErrorType ErrorType
	Message   string
}
