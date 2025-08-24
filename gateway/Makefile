.PHONY: build run clean test

# Переменные
BINARY_NAME=gateway
BUILD_DIR=build
MAIN_PATH=cmd/main.go

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