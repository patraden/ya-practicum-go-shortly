.PHONY: vet lint mocks code test clean build shortenertest

VETTOOL ?= $(shell which statictest)
SOURCE_PATH ?= ${CURDIR}
BINARY_PATH ?= cmd/shortener/shortener
TEMP_FILE ?= data/service_storage.json

vet:
	@go vet -vettool=$(VETTOOL) ./...

lint:
	@goimports -e -w -local "github.com/patraden/ya-practicum-go-shortly" .
	@gofumpt -w ./cmd/shortener ./internal/app
	@golangci-lint run ./...

mocks:
	@mockgen -source=internal/app/repository/repository.go -destination=internal/app/mock/repository.go -package=mock URLRepository
	@mockgen -source=internal/app/urlgenerator/urlgenerator.go -destination=internal/app/mock/urlgenerator.go -package=mock URLGenerator

code: mocks
	@easyjson -all internal/app/dto/url.go
	@easyjson -all internal/app/repository/record.go


test:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

clean:
	@rm -f ./cmd/shortener/shortener
	@rm -f ./coverage.out
	@rm -f ./data/service_storage.json

build: clean
	@go build -buildvcs=false -o cmd/shortener/shortener ./cmd/shortener

shortenertest: build
	@echo "Running increment1 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration1\$$ -source-path=$(SOURCE_PATH) -binary-path=$(BINARY_PATH)
	@echo "Running increment2 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration2\$$ -source-path=$(SOURCE_PATH)
	@echo "Running increment3 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration3\$$ -source-path=$(SOURCE_PATH)
	@echo "Running increment4 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration4\$$ -binary-path=$(BINARY_PATH) -server-port=8181
	@echo "Running increment5 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration5\$$ -binary-path=$(BINARY_PATH) -server-port=8181
	@echo "Running increment6 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration6\$$ -source-path=$(SOURCE_PATH)
	@echo "Running increment7 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration7\$$ -source-path=$(SOURCE_PATH) -binary-path=$(BINARY_PATH)
	@echo "Running increment8 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration8\$$ -source-path=$(SOURCE_PATH) -binary-path=$(BINARY_PATH)
	@echo "Running increment9 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration9\$$ -source-path=$(SOURCE_PATH) -binary-path=$(BINARY_PATH) -file-storage-path=$(TEMP_FILE)