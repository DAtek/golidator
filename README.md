[![codecov](https://codecov.io/gh/DAtek/golidator/graph/badge.svg?token=1QYUBN9NDN)](https://codecov.io/gh/DAtek/golidator)

# Golidator
Lightweight, extensible validation library inspired by [Pydantic](https://github.com/pydantic/pydantic)

## Features
- You can use any context in the validation logic e.g. database connection
- Fields can have multiple errors
- Returned errors show the exact location of the erroneous fileds, you can also provide error context
- During validation you can mutate the struct e.g. use `time.Parse()` and store the result into a different field


## Usage
See [`validator_test.go`](https://github.com/DAtek/golidator/blob/main/validator_test.go)

## Example
```go
package main

import (
	"fmt"
	"github.com/DAtek/golidator"
	"strings"
)

type CreateUserInput struct {
	Username string
	Password string
}

func (obj *CreateUserInput) GetValidators(ctx ...interface{}) golidator.ValidatorCollection {
	db, ok := ctx[0].(IDatabase)
	if !ok {
		panic("db not provided")
	}

	return golidator.ValidatorCollection{
		{Field: "Username", Function: func() *golidator.ValueError {
			if db.UserExists(obj.Username) {
				return &golidator.ValueError{ErrorType: golidator.ErrorType("ALREADY_EXISTS")}
			}
			return nil
		}},
		{Field: "Password", Function: func() *golidator.ValueError {
			minLength := 5
			if len(obj.Password) < minLength {
				return &golidator.ValueError{
					ErrorType: golidator.ErrorType("TOO_SHORT"),
					Context:   map[string]any{"min": minLength},
				}
			}
			return nil
		}},
		{Field: "Password", Function: func() *golidator.ValueError {
			if !strings.ContainsAny(obj.Password, "0123456789") {
				return &golidator.ValueError{
					ErrorType: golidator.ErrorType("MISSING_NUMBER"),
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

func main() {
	db := &FakeDatabase{}
	input := &CreateUserInput{"John", "asd"}
	err := golidator.Validate(input, db)
	fmt.Printf("%v\n", err)
}

```
>>>
```
3 validation error(s) for main.CreateUserInput
Username: ALREADY_EXISTS
Password: TOO_SHORT
Password: MISSING_NUMBER
```
