.PHONY: pg mysql_run mysql_start mysql_stop pg pg_start pg_stop pg_connect

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

# https://hub.docker.com/_/mysql
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

mysql_start:
	@$(DOCKER) start $(CONTAINER_MYSQL)

mysql_stop:
	@$(DOCKER) stop $(CONTAINER_MYSQL)

mysql_connect:
	@mysql -h 127.0.0.1 -P 3306 -u $(MYSQL_USER) -p$(MYSQL_PASSWORD) $(MYSQL_DATABASE)

pg:
	@$(DOCKER) run --name $(CONTAINER_POSTGRES) \
    -e POSTGRES_USER=$(POSTGRES_USER) \
    -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
    -e POSTGRES_DB=$(POSTGRES_DB) \
    -p 5432:5432 \
    -d postgres:15.1

pg_start:
	@$(DOCKER) start $(CONTAINER_POSTGRES)

pg_stop:
	@$(DOCKER) stop $(CONTAINER_POSTGRES)

pg_connect:
	@psql -h 127.0.0.1 -p 5432 -U $(POSTGRES_USER) -d $(POSTGRES_DB)