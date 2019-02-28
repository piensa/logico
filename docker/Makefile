
DOCKER_VERSION := $(shell docker --version 2>/dev/null)
DOCKER_COMPOSE_VERSION := $(shell docker-compose --version 2>/dev/null)

DB_USER ?= dbuser
DB_PW ?= secret
DB_NAME ?= accesscontroldb
DB_HOST ?= localhost

ENV_BROWSER_IDP_HOST ?= localhost

# Make sure this port corresponds to the one mapped in the docker-compose file.
DB_PORT ?= 5432

all:
ifndef DOCKER_VERSION
    $(error "command docker is not available, please install Docker")
endif
ifndef DOCKER_COMPOSE_VERSION
    $(error "command docker-compose is not available, please install Docker")
endif

export DATABASE_URL=postgres://${DB_USER}:${DB_PW}@postgresd:${DB_PORT}/${DB_NAME}?sslmode=disable
export PG_USER=${DB_USER}
export PG_PW=${DB_PW}
export PG_DB=${DB_NAME}
export PG_PORT=${DB_PORT}
export PG_HOST=${DB_HOST}

export BROWSER_IDP_HOST=${ENV_BROWSER_IDP_HOST}

# TEST database info.
export DB_USERS=testusers
export DB_TABLE=users

up:
	@echo "[INFO] Initializing docker containers"
	docker-compose up -d

	sleep 3

initdb:
	@echo "[INFO] Creating database ${DB_USERS}"
	docker exec -it postgres sh -c \
	"PGPASSWORD=${PG_PW} psql -U ${PG_USER} postgres -h ${PG_HOST} -p 5432 -c 'CREATE DATABASE ${DB_USERS}'"

	@echo "[INFO] Creating table ${DB_USERS}"
	docker exec -it postgres sh -c \
	"PGPASSWORD=${PG_PW} psql -U ${PG_USER} ${DB_USERS} -h ${PG_HOST} -p 5432 -f /dump/insert.sql"

clients:
	docker exec -it hydra hydra clients import --endpoint http://localhost:4445 \
		/config/hydra/app.json \
		/config/hydra/oathkeeper.json \
		/config/hydra/keto.json \
		/config/hydra/testapp.json \

migrate: clients
	@echo "[INFO] Migrating oathkeeper rules"
	docker exec -it ok-api oathkeeper rules import --endpoint http://localhost:4456 /config/ok-rule.json

down:
	docker-compose down

list:
	docker-compose ps

logs:
	docker-compose logs hydra

issue_token:
	docker exec -it hydra hydra token user --client-id testapp \
	--client-secret testapp \
	--endpoint http://localhost:4444 \
	--port 5555 \
	--scope openid,offline,read,write

start: up initdb

restart: down up initdb

