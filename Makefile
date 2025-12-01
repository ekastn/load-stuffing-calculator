.PHONY: all build run test clean fmt swag swagfmt

APP_NAME=api
BUILD_DIR=bin
CMD_DIR=cmd/api

all: fmt swagfmt swag test build

fmt:
	go fmt ./...

swagfmt:
	swag fmt -g cmd/api/main.go

swag:
	swag init -g cmd/api/main.go -o internal/docs

test:
	go test -v ./...

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)

run: build
	./$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)
