.PHONY: help build build-local up down logs ps migrate test

TAG := latest
build:
	docker build -t my-meal-journal:$(TAG) .

build-local:
	docker compose -f ./deploy/compose.yaml build --no-cache

up:
	docker compose -f ./deploy/compose.yaml up -d

down:
	docker compose -f ./deploy/compose.yaml down

logs:
	docker compose -f ./deploy/compose.yaml logs

ps:
	docker compose -f ./deploy/compose.yaml ps

migrate:
	mysql -u test -ptest -h 127.0.0.1 test < deploy/mysql/schema.sql

test: up
	echo "waiting for database setup"; sleep 20;
	go test -v -race -shuffle=on ./...

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
