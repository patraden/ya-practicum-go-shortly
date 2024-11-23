VETTOOL ?= $(shell which statictest)
SOURCE_PATH ?= ${CURDIR}
BINARY_PATH ?= cmd/shortener/shortener
TEMP_FILE ?= data/service_storage.json
DATABASE_DSN ?= postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable

.PHONY: vet
vet:
	@go vet -vettool=$(VETTOOL) ./...


.PHONY: lint
lint:
	@goimports -e -w -local "github.com/patraden/ya-practicum-go-shortly" .
	@gofumpt -w ./cmd/shortener ./internal/app
	@golangci-lint run ./...


.PHONY: mocks
mocks:
	@mockgen -source=internal/app/repository/repository.go -destination=internal/app/mock/repository.go -package=mock URLRepository
	@mockgen -source=internal/app/service/urlgenerator/urlgenerator.go -destination=internal/app/mock/urlgenerator.go -package=mock URLGenerator
	@mockgen -source=internal/app/memento/originator.go -destination=internal/app/mock/originator.go -package=mock Originator
	@mockgen -source=internal/app/service/shortener/shortener.go -destination=internal/app/mock/shortener.go -package=mock URLShortener


.PHONY: code
code: mocks
	@easyjson -all internal/app/dto/dto.go
	@easyjson -all internal/app/domain/urlmapping.go


.PHONY: test
test:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: clean
clean:
	@rm -f ./cmd/shortener/shortener
	@rm -f ./coverage.out
	@rm -f ./data/service_storage.json


.PHONY: build
build: clean
	@go build -buildvcs=false -o cmd/shortener/shortener ./cmd/shortener


.PHONY: shortenertest
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
	@echo "Running increment10 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration10\$$ -source-path=$(SOURCE_PATH) -binary-path=$(BINARY_PATH) -database-dsn=$(DATABASE_DSN)
	@echo "Running increment11 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration11\$$ -binary-path=$(BINARY_PATH) -database-dsn=$(DATABASE_DSN)
	@echo "Running increment12 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration12\$$ -binary-path=$(BINARY_PATH) -database-dsn=$(DATABASE_DSN)
	@echo "Running increment13 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration13\$$ -binary-path=$(BINARY_PATH) -database-dsn=$(DATABASE_DSN)
	@echo "Running increment14 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration14\$$ -binary-path=$(BINARY_PATH) -database-dsn=$(DATABASE_DSN)


.PHONY: goose-init
goose-init:
	@goose --dir migrations -s create app_schema sql
	@goose --dir migrations -s create app_repository sql
	@goose --dir migrations -s create app_grants sql


.PHONY: goose-status
goose-status:
	@goose -dir migrations postgres ${DATABASE_DSN} status


.PHONY: goose-up
goose-up:
	@goose -dir migrations postgres ${DATABASE_DSN} up


.PHONY: goose-down
goose-down:
	@goose -dir migrations postgres ${DATABASE_DSN} down


.PHONY: sqlc
sqlc:
	@sqlc generate