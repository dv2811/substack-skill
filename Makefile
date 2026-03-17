# Substack Skill - Build and Development Commands
# Delegates to tool_build.sh for all tool operations

-include .env

TOOL_BUILD := ./tool_build.sh

.PHONY: help
help: ## Show available tools
	@$(TOOL_BUILD)

.PHONY: build
build: ## Build and install a tool (usage: make build TOOL=<tool-name>)
	@if [ -z "$(TOOL)" ]; then \
		echo "Error: TOOL not specified"; \
		echo "Usage: make build TOOL=<tool-name>"; \
		echo "Available tools:"; \
		$(TOOL_BUILD) 2>&1 | grep -A100 "Available tools:" | tail -n +2 | head -n -2; \
		exit 1; \
	fi
	@$(TOOL_BUILD) $(TOOL)

.PHONY: install
install: ## Install a tool (usage: make install TOOL=<tool-name>)
	@if [ -z "$(TOOL)" ]; then \
		echo "Error: TOOL not specified"; \
		echo "Usage: make install TOOL=<tool-name>"; \
		exit 1; \
	fi
	@$(TOOL_BUILD) $(TOOL)

.PHONY: setup
setup: ## Run full setup for a tool (usage: make setup TOOL=<tool-name>)
	@if [ -z "$(TOOL)" ]; then \
		echo "Error: TOOL not specified"; \
		echo "Usage: make setup TOOL=<tool-name>"; \
		exit 1; \
	fi
	@$(TOOL_BUILD) $(TOOL)

.PHONY: uninstall
uninstall: ## Remove installed binary (usage: make uninstall TOOL=<tool-name>)
	@if [ -z "$(TOOL)" ]; then \
		echo "Error: TOOL not specified"; \
		echo "Usage: make uninstall TOOL=<tool-name>"; \
		exit 1; \
	fi
	@echo "Removing $(TOOL) binary..."
	@if [ -f $$HOME/bin/substack ]; then \
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

.PHONY: clean-tool
clean-tool: ## Remove build artifacts for a specific tool (usage: make clean-tool TOOL=<tool-name>)
	@if [ -z "$(TOOL)" ]; then \
		echo "Error: TOOL not specified"; \
		echo "Usage: make clean-tool TOOL=<tool-name>"; \
		exit 1; \
	fi
	rm -rf tools/$(TOOL)/bin
	@echo "✓ Cleaned $(TOOL) build artifacts"

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
