.PHONY: build run clean test

# Переменные
BINARY_NAME=gateway
BUILD_DIR=build
MAIN_PATH=cmd/main.go
GOPATH= $(shell go env GOPATH)


PROTOC_GEN_OPENAPIV2=$(GOPATH)/bin/protoc-gen-openapiv2

install-deps:
	GOBIN=$(GOPATH)/bin go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0


gen-openapi:
	protoc -I . \
		-I $(GOPATH)/pkg/mod \
		--openapiv2_out ./docs \
		--openapiv2_opt logtostderr=true \
		auth/api/*.proto event/api/*.proto user-interact/api/*.proto

# Сборка
build:
	@echo "Building gateway service..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Запуск
run: build
	@echo "Running gateway service..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Запуск в режиме разработки
dev:
	@echo "Running gateway service in development mode..."
	@go run $(MAIN_PATH)

# Очистка
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

# Тесты
test:
	@echo "Running tests..."
	@go test ./...

# Установка зависимостей
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Проверка кода
lint:
	@echo "Running linter..."
	@golangci-lint run

# Генерация документации
docs:
	@echo "Generating documentation..."
	@swag init -g $(MAIN_PATH) -o docs 