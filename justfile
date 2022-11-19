coverfile := ".coverage"
pkgs := "."


test:
    go test {{ pkgs }}


test-cover:
    go test -coverprofile {{ coverfile }} {{ pkgs }}


show-coverage:
    go tool cover -html={{ coverfile }}


test-and-show-covarage: test-cover show-coverage
