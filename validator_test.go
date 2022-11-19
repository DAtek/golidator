package golidator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	t.Run("ValidationError contains proper location", func(t *testing.T) {
		animal := &Animal{}

		err := Validate(animal)

		assert.Equal(1, len(err.Errors))
		assert.Equal("Kind", err.Errors[0].Location)
	})

	t.Run("ValidationError contains proper error type", func(t *testing.T) {
		animal := &Animal{}

		err := Validate(animal)

		assert.Equal(1, len(err.Errors))
		assert.Equal(ErrorType("EMPTY"), err.Errors[0].ErrorType)
	})

	t.Run("Object can be mutated during validation", func(t *testing.T) {
		animal := &Animal{Kind: "snake"}

		err := Validate(animal)

		assert.Nil(err)
		assert.Equal("lizard with no legs", animal.KindAlias)
	})

	t.Run("Validate list of other structs", func(t *testing.T) {
		herd := &Herd{
			Animals: []*Animal{
				{Kind: "dog"},
				{},
			},
			Alpha: &Animal{Kind: "lion"},
		}

		err := Validate(herd)

		assert.Equal(1, len(err.Errors))
		assert.Equal("Animals.1.Kind", err.Errors[0].Location)
	})

	t.Run("Validate child struct", func(t *testing.T) {
		herd := &Herd{
			Alpha: &Animal{},
		}

		err := Validate(herd)

		assert.Equal(1, len(err.Errors))
		assert.Equal("leader.Kind", err.Errors[0].Location)
	})

	t.Run("No error if lists or child structs are nil", func(t *testing.T) {
		herd := &Herd{}

		err := Validate(herd)

		assert.Nil(err)
	})

	t.Run("Validate with context", func(t *testing.T) {
		input := &CreateUserInput{Username: "John", Password: "abcde4"}
		db := &FakeDatabase{}

		err := Validate(input, db)

		assert.Equal(1, len(err.Errors))
		assert.Equal("Username", err.Errors[0].Location)
		assert.Equal(ErrorType("ALREADY_EXISTS"), err.Errors[0].ErrorType)
	})

	t.Run("A field can have multiple errors", func(t *testing.T) {
		input := &CreateUserInput{"Emmett Lathrop Brown, Ph.D.", "abcd"}
		db := &FakeDatabase{}

		err := Validate(input, db)

		assert.Equal(2, len(err.Errors))
		assert.Equal("Password", err.Errors[0].Location)
		assert.Equal("Password", err.Errors[1].Location)
		assert.Equal(ErrorType("TOO_SHORT"), err.Errors[0].ErrorType)
		assert.Equal(ErrorType("MISSING_NUMBER"), err.Errors[1].ErrorType)
	})
}

type Herd struct {
	Animals []*Animal
	Alpha   *Animal
}

func (obj *Herd) getValidators(ctx ...interface{}) ValidatorCollection {
	validarots := GetValidatorsForList("Animals", obj.Animals)
	validarots = append(validarots, GetValidatorsForObject("leader", obj.Alpha)...)
	return validarots

}

type Animal struct {
	Kind      string
	KindAlias string
}

func (obj *Animal) getValidators(ctx ...interface{}) ValidatorCollection {
	return ValidatorCollection{
		{"Kind", func() *ValueError {
			if obj.Kind == "" {
				return &ValueError{
					ErrorType: "EMPTY",
				}
			}

			if obj.Kind == "snake" {
				obj.KindAlias = "lizard with no legs"
			}
			return nil
		}},
	}
}

type CreateUserInput struct {
	Username string
	Password string
}

func (obj *CreateUserInput) getValidators(ctx ...interface{}) ValidatorCollection {
	db, ok := ctx[0].(IDatabase)
	if !ok {
		panic("oops")
	}

	return ValidatorCollection{
		{"Username", func() *ValueError {
			if db.UserExists(obj.Username) {
				return &ValueError{ErrorType: ErrorType("ALREADY_EXISTS")}
			}
			return nil
		}},
		{"Password", func() *ValueError {
			minLength := 5
			if len(obj.Password) < minLength {
				return &ValueError{
					ErrorType: ErrorType("TOO_SHORT"),
					Context:   map[string]any{"min": minLength},
				}
			}
			return nil
		}},
		{"Password", func() *ValueError {
			if !strings.ContainsAny(obj.Password, "0123456789") {
				return &ValueError{
					ErrorType: ErrorType("MISSING_NUMBER"),
				}
			}
			return nil
		}},
	}
}

type IDatabase interface {
	UserExists(username string) bool
}

type FakeDatabase struct{}

func (db *FakeDatabase) UserExists(username string) bool {
	return username == "John"
}
