.PHONY: build

RM := rm -f
CMP := cmp -s
GO := go

BIN_NAME = golangci-lint
BUILD_PATH = bin/$(BIN_NAME)
INSTALL_DIR = $(HOME)/.local/bin
INSTALL_PATH = $(INSTALL_DIR)/$(BIN_NAME)
COMPARE_BINS = $(CMP) $(BUILD_PATH) $(INSTALL_PATH)

build:
	@echo "Building $(BUILD_PATH)..."
	@$(GO) run build/main.go
	@echo "Built    $(BUILD_PATH)!"

install: check-built
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_PATH) $(INSTALL_PATH)
	@echo "Installed $(BIN_NAME) at $(INSTALL_DIR)"

uninstall:
	@$(RM) $(INSTALL_PATH)

check: check-path check-built check-bin

check-path:
	@which $(BIN_NAME) > /dev/null 2>&1 || { \
		echo "Error: $(BIN_NAME) is not in PATH"; \
		exit 1; \
	}
	@if [ "`which $(BIN_NAME)`" = "$(INSTALL_PATH)" ]; then \
		echo "Detected $(BIN_NAME) at desired place."; \
	else \
		echo "Error: Detected $(BIN_NAME) at undesired place. Please check PATH alignment"; \
	fi

check-built:
	@if [ -f $(BUILD_PATH) ]; then \
		echo "Built binary exists"; \
	else \
		echo "Error: Built binary does not exist"; \
		echo "Please run 'make build'"; \
		exit 1; \
	fi

check-bin:
	@if $(COMPARE_BINS); then \
		echo "Binary files are identical"; \
	else \
		echo "Binary files are different"; \
	fi

clean:
	@$(RM) $(BUILD_PATH)

update:
	@go list -m -f '{{if and (ne .Main true) (not .Indirect)}}{{.Path}}{{end}}' all | xargs -n1 go get -u=patch
	@go mod tidy

outdated: outdated-patch outdated-minor

outdated-patch::
	@echo "================================"
	@echo "- patch update -----------------"
	@echo "                                "
	@go list -m -u -json all | jq -r 'select(.Main != true and .Indirect != true and .Update) | select( (.Version | split(".")[:2]) == (.Update.Version | split(".")[:2])) | "\(.Path): \(.Version) → \(.Update.Version)"'
	@echo "--------------------------------"
	@echo "Run 'make update-patch' to update the dependencies."
	@echo "================================"

outdated-minor:
	@echo "================================"
	@echo "- minor update -----------------"
	@echo "--------------------------------"
	@go list -m -u -f '{{if and (not .Main) (not .Indirect) .Update}}{{.Path}}: {{.Version}} → {{.Update.Version}}{{end}}' all
	@echo "--------------------------------"
	@echo "================================"
