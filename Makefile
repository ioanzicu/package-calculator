include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: run
run: ## Run the application
	go run cmd/main.go

.PHONY: build
build: ## Build the application binary
	go build -o bin/app cmd/main.go

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: deps
deps: ## Install dependencies and tools
	go mod download
	go mod tidy
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

.PHONY: generate
generate: ## Generate code from sqlc
	$(LOCAL_BIN)/sqlc generate

##@ Database

.PHONY: migration-create
migration-create: ## Create a new migration file (usage: make migration-create name=create_users_table)
	$(LOCAL_BIN)/goose -dir migrations create $(name) sql

.PHONY: migration-up
migration-up: ## Apply all up migrations
	$(LOCAL_BIN)/goose -dir migrations postgres "$(PG_DSN)" up

.PHONY: migration-down
migration-down: ## Rollback the last migration
	$(LOCAL_BIN)/goose -dir migrations postgres "$(PG_DSN)" down

.PHONY: migration-status
migration-status: ## Check migration status
	$(LOCAL_BIN)/goose -dir migrations postgres "$(PG_DSN)" status

##@ Docker

.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t app .

.PHONY: docker-run
docker-run: ## Run docker container
	docker run -p 8080:8080 --env-file .env app

##@ Docker Compose

.PHONY: compose-up
compose-up: ## Start services with docker-compose
	docker-compose up -d

.PHONY: compose-down
compose-down: ## Stop services with docker-compose
	docker-compose down

.PHONY: compose-logs
compose-logs: ## View logs of services
	docker-compose logs -f