.PHONY: all install test

install:
	@go build -o keyper . && mv keyper ~/.local/bin

test:
	@go test ./... -cover -json | tparse -pass -trimpath github.com/jsnfwlr/keyper-cli

coverage:
	@go test -coverprofile=coverage.out ./... -json | tparse -pass -trimpath github.com/jsnfwlr/keyper-cli
	@go tool cover -html=coverage.out -o coverage.html