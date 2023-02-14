.PHONY: help build up down logs ps test

TAG := latest
build:
	docker build -t yutaroyamanaka/my-mea-journal:$(TAG) .

up:
	docker compose -f ./deploy/compose.yaml up -d

down:
	docker compose -f ./deploy/compose.yaml down

logs:
	docker compose -f ./deploy/compose.yaml logs

ps:
	docker compose -f ./deploy/compose.yaml ps

migrate:
	mysql -u test -p test -h 127.0.0.1 < deploy/mysql/schema.sql

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
