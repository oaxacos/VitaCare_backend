DB_MIGRATIONS_PATH=migrations
CONFIG_FILE=config/config.yaml
ifndef VERBOSE
	MAKEFLAGS += --no-print-directory
endif

PORT= $(shell yq '.database.port' < $(CONFIG_FILE))
DB_NAME= $(shell yq '.database.dbname' < $(CONFIG_FILE))
DB_USER= $(shell yq '.database.username' < $(CONFIG_FILE))
DB_PASSWORD= $(shell yq '.database.password' < $(CONFIG_FILE))
DB_HOST= $(shell yq '.database.host' < $(CONFIG_FILE))
DSN= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(PORT)/$(DB_NAME)?sslmode=disable

## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## install: install the dependencies
.PHONY: install
install:
	go mod download

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

## shell: run the shell
.PHONY: shell
shell:
	devbox shell

db-create:
	dbmate --url ${DSN} create

## db-up: create the database
.PHONY: db-up
db-up: db-create migrate-up

## db-drop: drop the database
.PHONY: db-drop
db-drop:
	dbmate --url ${DSN} drop

## migrate-up: apply the migrations
.PHONY: migrate-up
migrate-up:
	dbmate --url ${DSN} --schema-file ${DB_MIGRATIONS_PATH} up

## migrate-down: rollback the migrations
.PHONY: migrate-down
migrate-down:
	dbmate --url ${DSN} --migrations-dir ${DB_MIGRATIONS_PATH} down

name=""
## migrate-new name=$1: create a new migration
.PHONY: migrate-new
migrate-new:
	@if [ -z ${name} ]; then \
		echo "name is required"; \
		exit 1; \
	fi
	@echo "Creating migration: ${name}"
	dbmate --url ${DSN} -d ${DB_MIGRATIONS_PATH} new ${name}