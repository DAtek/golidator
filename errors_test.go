package golidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Building struct{}

func TestValidationError(t *testing.T) {
	t.Run("Validation error constructs proper error", func(t *testing.T) {
		valueError := ValueError{
			ErrorType: "MISSING",
		}

		err := &ValidationError{
			source: &Building{},
			Errors: []*FieldError{
				{
					Location:   "window",
					ValueError: valueError,
				},
				{
					Location:   "door",
					ValueError: valueError,
				},
			},
		}

		errorString := err.Error()
		wantedErrorString := "2 validation error(s) for golidator.Building\nwindow: MISSING\ndoor: MISSING"

		assert.Equal(t, wantedErrorString, errorString)
	})
}
