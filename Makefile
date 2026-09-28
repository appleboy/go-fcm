GO ?= go
TOOLS_MOD := -modfile=go.tools.mod

.PHONY: test
test:
	@$(GO) test -v -cover -coverprofile coverage.out ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

clean:
	go clean -x -i ./...

.PHONY: fmt
fmt: ## Format Go files using golangci-lint
	$(GO) tool $(TOOLS_MOD) golangci-lint fmt

.PHONY: lint
lint: ## Run golangci-lint
	$(GO) tool $(TOOLS_MOD) golangci-lint run

.PHONY: install-tools fmt lint
install-tools: ## Download pinned Go tools
	$(GO) mod download $(TOOLS_MOD)
