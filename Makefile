VETTOOL ?= $(shell which statictest)
SOURCE_PATH ?= ${CURDIR}
BINARY_PATH ?= cmd/shortener/shortener
TEMP_FILE ?= data/service_storage.json
DATABASE_DSN ?= postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable

SERVER_PORT ?= 8443
SERVER_ADDRESS ?= 0.0.0.0:${SERVER_PORT}
BASE_URL ?= https://localhost:${SERVER_PORT}/
ENABLE_HTTPS ?= true
DOCKER := $(shell which docker)
DOCKER_COMPOSE_PATH := ./deployments/docker-compose.yml
CONTAINER_MYSQL ?= ya_mysql
MYSQL_DATABASE ?= mysql
MYSQL_USER ?= mysql
MYSQL_PASSWORD ?= mysql
MYSQL_ROOT_PASSWORD ?= mysql
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= postgres
POSTGRES_DB ?= praktikum
BUILD_DATE := $(shell date -u +"%d.%m.%Y")
BUILD_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "N/A")
VERSION_PACKAGE := github.com/patraden/ya-practicum-go-shortly/internal/app/version

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
	@mockgen -source=internal/app/service/statsprovider/statsprovider.go -destination=internal/app/mock/statsprovider.go -package=mock StatsProvider


.PHONY: code
code: mocks
	@easyjson -all internal/app/config/config.go
	@easyjson -all internal/app/dto/dto.go
	@easyjson -all internal/app/domain/urlmapping.go


.PHONY: test
test:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: clean
clean:
	@rm -f ./cmd/shortener/shortener
	@rm -f ./cmd/staticlint/staticlint
	@rm -f ./coverage.out
	@rm -f ./data/service_storage.json


.PHONY: build
build: clean
	@go build \
		-ldflags="-s -w -X $(VERSION_PACKAGE).buildVersion=$(BUILD_VERSION) -X $(VERSION_PACKAGE).buildDate=$(BUILD_DATE) -X $(VERSION_PACKAGE).buildCommit=$(BUILD_COMMIT)" \
		-o cmd/shortener/shortener ./cmd/shortener

.PHONY: run
run: 
	@go run \
		-ldflags="-X $(VERSION_PACKAGE).buildVersion=$(BUILD_VERSION) -X $(VERSION_PACKAGE).buildDate=$(BUILD_DATE) -X $(VERSION_PACKAGE).buildCommit=$(BUILD_COMMIT)" \
		cmd/shortener/main.go


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

.PHONY: docker\:pg
docker\:pg:
	@SERVER_ADDRESS=$(SERVER_ADDRESS) \
	POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_DB=$(POSTGRES_DB) \
	SERVER_PORT=$(SERVER_PORT) \
	docker-compose -f $(DOCKER_COMPOSE_PATH) up -d postgres

.PHONY: pg\:connect
pg\:connect:
	@psql -h 127.0.0.1 -p 5432 -U $(POSTGRES_USER) -d $(POSTGRES_DB)


.PHONY: staticlint
staticlint:
	@echo "Running staticlint..."
	@./cmd/staticlint/staticlint ./...


.PHONY: staticlint\:help
staticlint\:help:
	@./cmd/staticlint/staticlint --help

.PHONY: staticlint\:build
staticlint\:build:
	@echo "Building staticcheck binary..."
	@go build -ldflags="-s -w" -o cmd/staticlint/staticlint ./cmd/staticlint/
	@chmod +x cmd/staticlint/staticlint


.PHONY: docker\:up 
docker\:up:
	@SERVER_ADDRESS=$(SERVER_ADDRESS) \
	BASE_URL=$(BASE_URL) \
	FILE_STORAGE_PATH=/app/data/service_storage.json \
	ENABLE_HTTPS=$(ENABLE_HTTPS) \
	POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_DB=$(POSTGRES_DB) \
	SERVER_PORT=$(SERVER_PORT) \
	docker-compose -f $(DOCKER_COMPOSE_PATH) up -d

.PHONY: docker\:down
docker\:down:
	@SERVER_ADDRESS=$(SERVER_ADDRESS) \
	BASE_URL=$(BASE_URL) \
	FILE_STORAGE_PATH=/app/data/service_storage.json \
	ENABLE_HTTPS=$(ENABLE_HTTPS) \
	POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_DB=$(POSTGRES_DB) \
	SERVER_PORT=$(SERVER_PORT) \
	docker-compose -f $(DOCKER_COMPOSE_PATH) down -v

.PHONY: docker\:stop
docker\:stop:
	@SERVER_ADDRESS=$(SERVER_ADDRESS) \
	BASE_URL=$(BASE_URL) \
	FILE_STORAGE_PATH=/app/data/service_storage.json \
	ENABLE_HTTPS=$(ENABLE_HTTPS) \
	POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_DB=$(POSTGRES_DB) \
	SERVER_PORT=$(SERVER_PORT) \
	docker-compose -f $(DOCKER_COMPOSE_PATH) stop

.PHONY: docker\:build
docker\:build: docker\:down
	@BUILD_DATE=$(BUILD_DATE) \
	BUILD_COMMIT=$(BUILD_COMMIT) \
	BUILD_VERSION=$(BUILD_VERSION) \
	VERSION_PACKAGE=$(VERSION_PACKAGE) \
	POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_DB=$(POSTGRES_DB) \
	SERVER_PORT=$(SERVER_PORT) \
	docker-compose -f $(DOCKER_COMPOSE_PATH) build \
		--build-arg VERSION_PACKAGE=$(VERSION_PACKAGE) \
		--build-arg BUILD_VERSION=$(BUILD_VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg BUILD_COMMIT=$(BUILD_COMMIT) \
		--no-cache
	$(MAKE) docker\:up
