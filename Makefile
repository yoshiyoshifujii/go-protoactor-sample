APP_NAME := go-protoactor-sample

.PHONY: help build run test fmt tidy clean

help:
	@echo "Targets:"
	@echo "  make build  Build the binary"
	@echo "  make run    Run the sample"
	@echo "  make test   Run tests"
	@echo "  make fmt    Format Go files"
	@echo "  make tidy   Tidy module dependencies"
	@echo "  make clean  Remove build artifacts"

build:
	go build -o bin/$(APP_NAME) ./cmd

run:
	go run ./cmd

test:
	go test ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

clean:
	rm -rf bin
