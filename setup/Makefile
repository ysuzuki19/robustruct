.PHONY: build

RM := rm -f
CMP := cmp -s
GO := go

BIN_NAME = golangci-lint
BUILD_PATH = bin/$(BIN_NAME)
INSTALL_PATH = $(HOME)/.local/bin/$(BIN_NAME)
COMPARE_BINS = $(CMP) $(BUILD_PATH) $(INSTALL_PATH)

build:
	@echo "Building $(BUILD_PATH)..."
	@$(GO) run build/main.go
	@echo "Built    $(BUILD_PATH)!"

install: check-built
	@cp $(BUILD_PATH) $(INSTALL_PATH)

uninstall:
	@$(RM) $(INSTALL_PATH)

check: check-path check-built check-bin

check-built:
	@if [ -f $(BUILD_PATH) ]; then \
		echo "Built binary exists"; \
	else \
		echo "Error: Built binary does not exist"; \
		echo "Please run 'make build'"; \
		exit 1; \
	fi

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

check-bin:
	@if $(COMPARE_BINS); then \
		echo "Binary files are identical"; \
	else \
		echo "Binary files are different"; \
	fi

clean:
	@$(RM) $(BUILD_PATH)
