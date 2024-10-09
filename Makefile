.PHONY: vet lint test

VETTOOL := $(shell which statictest)

vet:
	@go vet -vettool=$(VETTOOL) ./... || exit 1

lint:
	@../bin/golangci-lint run ./... || exit 1

test:
	@go test -v ./...