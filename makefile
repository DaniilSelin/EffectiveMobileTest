BIN_DIR=bin
BIN_SERVER=$(BIN_DIR)/server
MAIN_FILE=cmd/server/main.go

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	go build -o $(BIN_SERVER) $(MAIN_FILE)
	@echo "complete!"

run: build
	./$(BIN_SERVER)

stop:
	kill -SIGINT $(shell lsof -ti:8080)

clean:
	rm -rf $(BIN_DIR)
	@echo "complete!"

help:
	@echo "Available commands:"
	@echo "  make build   - Build server and client"
	@echo "  make run     - Start the server"
	@echo "	 make stop 	  - Send SIGINT signal!!! Server Shutdown"
	@echo "  make clean   - Remove binaries and Unix socket"
