VETTOOL ?= $(shell which statictest)
SOURCE_PATH ?= ${CURDIR}
BINARY_PATH ?= cmd/shortener/shortener
TEMP_FILE ?= data/service_storage.json
DATABASE_DSN ?= postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable

DOCKER := $(shell which docker)
CONTAINER_MYSQL ?= ya_mysql
MYSQL_DATABASE ?= mysql
MYSQL_USER ?= mysql
MYSQL_PASSWORD ?= mysql
MYSQL_ROOT_PASSWORD ?= mysql
CONTAINER_POSTGRES ?= ya_postgres
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= postgres
POSTGRES_DB ?= praktikum

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
	@mockgen -source=internal/app/service/remover/remover.go -destination=internal/app/mock/remover.go -package=mock URLRemover


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
	@echo "Running increment15 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration15\$$ -binary-path=$(BINARY_PATH) -database-dsn=$(DATABASE_DSN)
	@echo "Running increment16 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration16\$$ -source-path=$(SOURCE_PATH)
	@echo "Running increment17 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration17\$$ -source-path=$(SOURCE_PATH)
	@echo "Running increment18 test"
	@shortenertestbeta -test.v -test.run=\^TestIteration18\$$ -source-path=$(SOURCE_PATH)


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


# https://hub.docker.com/_/mysql
.PHONY: mysql
mysql:
	@$(DOCKER) run --name $(CONTAINER_MYSQL) \
		-e MYSQL_DATABASE=$(MYSQL_DATABASE) \
		-e MYSQL_USER=$(MYSQL_USER) \
		-e MYSQL_PASSWORD=$(MYSQL_PASSWORD) \
		-e MYSQL_ROOT_PASSWORD=$(MYSQL_ROOT_PASSWORD) \
		-p 3306:3306 \
		-d mysql:9.1.0 \
		--character-set-server=utf8mb4 \
		--collation-server=utf8mb4_unicode_ci

.PHONY: mysql_start
mysql_start:
	@$(DOCKER) start $(CONTAINER_MYSQL)

.PHONY: mysql_stop
mysql_stop:
	@$(DOCKER) stop $(CONTAINER_MYSQL)

.PHONY: mysql_connect
mysql_connect:
	@mysql -h 127.0.0.1 -P 3306 -u $(MYSQL_USER) -p$(MYSQL_PASSWORD) $(MYSQL_DATABASE)

.PHONY: pg
pg:
	@$(DOCKER) run --name $(CONTAINER_POSTGRES) \
    -e POSTGRES_USER=$(POSTGRES_USER) \
    -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
    -e POSTGRES_DB=$(POSTGRES_DB) \
    -p 5432:5432 \
    -d postgres:15.1

.PHONY: pg_start
pg_start:
	@$(DOCKER) start $(CONTAINER_POSTGRES)

.PHONY: pg_stop
pg_stop:
	@$(DOCKER) stop $(CONTAINER_POSTGRES)

.PHONY: pg_connect
pg_connect:
	@psql -h 127.0.0.1 -p 5432 -U $(POSTGRES_USER) -d $(POSTGRES_DB)