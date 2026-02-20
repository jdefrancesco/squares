BINARY := squares
CMD := ./cmd/squares
BIN_DIR := bin
DIST_DIR := dist

APP_NAME := Squares
BUNDLE_ID ?= com.example.squares
VERSION ?= 0.1.0

.PHONY: build clean macos-app macos-app-clean icon

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) $(CMD)

clean:
	rm -rf $(BIN_DIR)

macos-app:
	@mkdir -p $(DIST_DIR)
	APP_NAME=$(APP_NAME) \
	EXECUTABLE_NAME=$(BINARY) \
	CMD_PKG=$(CMD) \
	BUNDLE_ID=$(BUNDLE_ID) \
	VERSION=$(VERSION) \
	DIST_DIR=$(DIST_DIR) \
	bash ./scripts/macos_app.sh

macos-app-clean:
	rm -rf $(DIST_DIR)/$(APP_NAME).app

icon:
	go run ./scripts/gen_icon.go -o assets/icon.png -s 1024
