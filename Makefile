.PHONY: vet lint test build shortenertest

VETTOOL ?= $(shell which statictest)
TEST_RUN ?= \^TestIteration1\$$
SOURCE_PATH ?= ${CURDIR}
BINARY_PATH ?= cmd/shortener/shortener

vet:
	@go vet -vettool=$(VETTOOL) ./...

lint:
	@../bin/golangci-lint run ./...

test:
	@go test -v ./...

build:
	@go build -buildvcs=false -o cmd/shortener/shortener ./cmd/shortener

shortenertest: build
	@shortenertestbeta -test.v -test.run=$(TEST_RUN) -source-path=$(SOURCE_PATH) -binary-path=$(BINARY_PATH)