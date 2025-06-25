# Default
.DEFAULT_GOAL := help

run: ## Run the Go application
	go run main.go

deps: ## Download dependencies
	go mod tidy

lint: ## Run gofmt on all files
	gofmt -s -w .

docker-up: ## Start Docker containers (Postgres)
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

clean: ## Clean binary and cache
	go clean -cache -testcache

# Help - List all tasks
help: ## Display all available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
