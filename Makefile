SHELL := /bin/bash
.DEFAULT_GOAL := help

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk ' \
		BEGIN { \
			FS = ":.*##"; \
			printf "\n\033[1mUsage:\033[0m\n  make \033[36m<target>\033[0m\n" \
		} \
		/^[a-zA-Z_0-9%-]+:.*?##/ { \
			printf "  \033[36m%-20s\033[0m \033[2;37m%-20s\033[0m\n", $$1, $$2 \
		} \
		/^##@/ { \
			printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
		} ' $(MAKEFILE_LIST)

##@ Test

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	go test ./aspect/... -v -race -cover

.PHONY: bench
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	go test ./aspect/... -bench=. -benchmem

.PHONY: converage
coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test ./aspect/... -coverprofile=coverage.out
	go tool cover -html=coverage.out

##@ Lint

.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	go clean
	rm -f coverage.out

##@ Example

.PHONY: ex-basic
ex-basic: ## Run example: Basic Usage (Before, After, AfterReturning)
	@echo "Running basic usage example..."
	go run docs/examples/01_basic_usage/main.go

.PHONY: ex-cache
ex-cache: ## Run example: Caching Pattern (Around, AfterReturning)
	@echo "Running caching pattern example..."
	go run docs/examples/02_caching_pattern/main.go

.PHONY: ex-auth
ex-auth: ## Run example: Authentication (Before, After, AfterReturning)
	@echo "Running authentication example..."
	go run docs/examples/03_authentication/main.go

.PHONY: ex-circuit-br
ex-circuit-br: ## Run example: Circuit Breaker (Before, Around, After)
	@echo "Running circuit breaker example..."
	go run docs/examples/04_circuit_breaker/main.go

.PHONY: ex-retry
ex-retry: ## Run example: Retry Pattern (Before, After)
	@echo "Running retry pattern example..."
	go run docs/examples/05_retry_pattern/main.go

.PHONY: ex-fluent-api
ex-fluent-api: ## Run example: Fluent API (Before, After)
	@echo "Running fluent-api example..."
	go run docs/examples/06_fluent_api/main.go

.PHONY: ex-real
ex-real: ## Run example: Real World Example (all cross-cutting concerns)
	@echo "Running real world example..."
	go run docs/examples/07_real_world_example/*.go

.PHONY: ex-context
ex-context: ## Run example: Context Propagation (context cancellation, deadlines, values)
	@echo "Running context propagation example..."
	go run docs/examples/08_context_example/main.go

.PHONY: ex-all
ex-all: ex-basic ex-cache ex-auth ex-circuit-br ex-retry ex-fluent-api ex-real ex-context ## Run all examples.

##@ Docs (Jekyll)

.PHONY: docs-deps
docs-deps: ## deps
	cd docs && \
	bundle config set --local path 'vendor/bundle' && \
	bundle install

.PHONY: docs-build
docs-build: ## Build static files from markdown.
	cd docs && bundle exec jekyll build

.PHONY: docs-serve
docs-serve: ## docs serve
	cd docs && bundle exec jekyll serve --livereload

##@ Git

.PHONY: rebase
rebase-%: ## Rebase current branch to the specified number of commits. Usage: make rebase-n
	@git rebase -i HEAD~$*