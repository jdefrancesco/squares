BINARY := squares
CMD := ./cmd/squares
BIN_DIR := bin

.PHONY: build clean

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) $(CMD)

clean:
	rm -rf $(BIN_DIR)
