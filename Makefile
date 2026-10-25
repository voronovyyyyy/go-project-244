.PHONY: build test test-coverage lint lint-fix run

COVERAGE_MIN ?= 80
LIB_PACKAGES := . ./parsers ./internal/diff ./formatters
LIB_COVERAGE := code,code/parsers,code/internal/diff,code/formatters
BINARY := bin/gendiff
ifeq ($(OS),Windows_NT)
BINARY := $(BINARY).exe
endif

build:
	go build -o $(BINARY) ./cmd/gendiff

test:
	go test ./...

test-coverage:
	go test -coverpkg=$(LIB_COVERAGE) -coverprofile=coverage.out $(LIB_PACKAGES)
	go tool cover -func=coverage.out
	go run ./scripts/check-coverage -profile coverage.out -min $(COVERAGE_MIN)

lint:
	golangci-lint run

lint-fix:
	golangci-lint fmt
	golangci-lint run --fix

run: build
	./$(BINARY) $(ARGS)
