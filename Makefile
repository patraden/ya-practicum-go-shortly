.PHONY: vet lint test

VETTOOL := $(shell which statictest)
TEST_RUN=\^TestIteration1\$$
SOURCE_PATH=.
BINARY_PATH=cmd/shortener/shortener

vet:
	@go vet -vettool=$(VETTOOL) ./... || exit 1

lint:
	@../bin/golangci-lint run ./... || exit 1

test:
	@go test -v ./...

shortenertest:
	@go build -buildvcs=false -o cmd/shortener/shortener ./cmd/shortener
	@shortenertestbeta -test.v -test.run=$(TEST_RUN) -source-path=$(SOURCE_PATH) -binary-path=$(BINARY_PATH)