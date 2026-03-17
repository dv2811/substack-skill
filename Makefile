# Substack Skill - Build and Development Commands
# Delegates to tool_build.sh for all tool operations

TOOL_BUILD := ./tool_build.sh
TOOL_NAME := substack-reader

.PHONY: help
help: ## Show available tools (delegates to tool_build.sh)
	@$(TOOL_BUILD)

.PHONY: build
build: ## Build and install substack-reader tool
	@$(TOOL_BUILD) $(TOOL_NAME)

.PHONY: install
install: ## Install substack-reader tool (same as build)
	@$(TOOL_BUILD) $(TOOL_NAME)

.PHONY: setup
setup: ## Run full setup for substack-reader (build + auth session)
	@$(TOOL_BUILD) $(TOOL_NAME)

.PHONY: uninstall
uninstall: ## Remove installed binary
	@if [ -f $$HOME/bin/$(TOOL_NAME) ]; then \
		rm -f $$HOME/bin/substack && echo "✓ Removed $$HOME/bin/substack"; \
	elif [ -f $$HOME/.local/bin/substack ]; then \
		rm -f $$HOME/.local/bin/substack && echo "✓ Removed $$HOME/.local/bin/substack"; \
	else \
		echo "Binary not found"; \
	fi

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin
	@echo "✓ Cleaned build artifacts"

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: test
test: ## Run tests
	go test -v ./internal/...

.PHONY: tidy
tidy: ## Tidy go.mod
	go mod tidy

.PHONY: deps
deps: ## Download dependencies
	go get ./...
