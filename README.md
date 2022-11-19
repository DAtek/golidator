[![codecov](https://codecov.io/gh/DAtek/golidator/branch/main/graph/badge.svg?token=1QYUBN9NDN)](https://codecov.io/gh/DAtek/golidator)

# Golidator
Lightweight, extensible validation library inspired by [Pydantic](https://github.com/pydantic/pydantic)

## Features
- You can use any context in the validation logic e.g. database connection
- Fields can have multiple errors
- Returned errors show the exact location of the erroneous fileds, you can also provide error context
- During validation you can mutate the struct e.g. use `time.Parse()` and store the result into a different field


## Usage
See `validator_test.go`