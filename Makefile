.PHONY: vet lint test build shortenertest

VETTOOL ?= $(shell which statictest)
SOURCE_PATH ?= ${CURDIR}
BINARY_PATH ?= cmd/shortener/shortener

vet:
	@go vet -vettool=$(VETTOOL) ./...

lint:
	@../bin/golangci-lint run ./...

test:
	@go test -v ./...

build:
	@rm -f ./cmd/shortener/shortener
	@go build -buildvcs=false -o cmd/shortener/shortener ./cmd/shortener

shortenertest: build
	@echo "Running increment1 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration1\$$ -source-path=$(SOURCE_PATH) -binary-path=$(BINARY_PATH)
	@echo "Running increment2 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration2\$$ -source-path=$(SOURCE_PATH)
	@echo "Running increment3 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration3\$$ -source-path=$(SOURCE_PATH)
	@echo "Running increment4 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration4\$$ -binary-path=$(BINARY_PATH) -server-port=8989
	@echo "Running increment5 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration5\$$ -binary-path=$(BINARY_PATH) -server-port=8787
	@echo "Running increment6 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration6\$$ -source-path=$(SOURCE_PATH)
	@echo "Running increment7 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration7\$$ -source-path=$(SOURCE_PATH) -binary-path=$(BINARY_PATH)
	@echo "Running increment8 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration8\$$ -binary-path=$(BINARY_PATH)