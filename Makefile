# Makefile для сборки Fyne приложения

# Имя приложения
APP_NAME := app
BUILD_DIR := build

# Go параметры
GO := go
GOFLAGS :=

# Цвета для вывода
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: all build clean run help deps

# Цель по умолчанию
all: build

# Создание папки build и сборка приложения
build: deps
	@echo "$(GREEN)Сборка приложения...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/app
	@echo "$(GREEN)Приложение собрано в $(BUILD_DIR)/$(APP_NAME)$(NC)"

# Установка зависимостей
deps:
	@echo "$(YELLOW)Установка зависимостей...$(NC)"
	$(GO) mod tidy
	$(GO) mod download

# Запуск приложения
run: build
	@echo "$(GREEN)Запуск приложения...$(NC)"
	./$(BUILD_DIR)/$(APP_NAME)

# Очистка
clean:
	@echo "$(YELLOW)Очистка...$(NC)"
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)Папка $(BUILD_DIR) удалена$(NC)"

# Помощь
help:
	@echo "$(GREEN)Доступные команды:$(NC)"
	@echo "  build  - Сборка приложения в папку $(BUILD_DIR)"
	@echo "  clean  - Удаление папки $(BUILD_DIR)"
	@echo "  run    - Сборка и запуск приложения"
	@echo "  deps   - Установка зависимостей"
	@echo "  help   - Показать эту справку"
	@echo "  all    - Сборка приложения (по умолчанию)"

# Проверка наличия Go
check-go:
	@if ! command -v $(GO) &> /dev/null; then \
		echo "$(RED)Ошибка: Go не установлен$(NC)"; \
		exit 1; \
	fi
