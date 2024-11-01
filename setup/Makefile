RM := rm -f

BIN_NAME = golangci-lint
BUILD_PATH = bin/$(BIN_NAME)
INSTALL_PATH = $(HOME)/.local/bin/$(BIN_NAME)

build:
	golangci-lint custom -v

install:
	@cp $(BUILD_PATH) $(INSTALL_PATH)

uninstall:
	@$(RM) $(INSTALL_PATH)

clean:
	@$(RM) $(BUILD_PATH)