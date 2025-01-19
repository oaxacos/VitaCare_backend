
## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## dev: run the development server
.PHONY: dev
dev:
	air -c .air.toml

## test: run the tests
.PHONY: test
test:
	go test -v ./...
## fmt: format the code
.PHONY: fmt
fmt:
	go fmt ./...

## build: build the binary
.PHONY: build
build:
	go build -o bin/ ./...
