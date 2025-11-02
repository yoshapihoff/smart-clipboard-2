# Makefile для сборки Fyne приложения

# Имя приложения
APP_NAME := SmartClipboard
BUILD_DIR := build
BUNDLE_DIR := $(BUILD_DIR)/$(APP_NAME).app/Contents
BINARY_NAME := $(APP_NAME)
ICON_PATH := assets/icon.png

# Go параметры
GO := go
GOFLAGS :=

# Цвета для вывода
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Проверка ОС
UNAME_S := $(shell uname -s)

.PHONY: all build clean run help deps mac-bundle

# Цель по умолчанию
all: build

# Создание папки build и сборка приложения
build: deps
	@echo "$(GREEN)Сборка приложения...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/app
	@echo "$(GREEN)Приложение собрано в $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Установка зависимостей
deps:
	@echo "$(YELLOW)Установка зависимостей...$(NC)"
	$(GO) mod tidy
	$(GO) mod download

# Запуск приложения
run: build
	@echo "$(GREEN)Запуск приложения...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME)

# Создание .app бандла для macOS
mac-bundle: build
	@echo "$(GREEN)Создание .app бандла для macOS...$(NC)"
	@# Создаем структуру папок
	@mkdir -p $(BUNDLE_DIR)/MacOS
	@mkdir -p $(BUNDLE_DIR)/Resources

	@# Копируем бинарник
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(BUNDLE_DIR)/MacOS/$(BINARY_NAME)

	@# Создаем Info.plist
	@echo '<?xml version="1.0" encoding="UTF-8"?>' > $(BUNDLE_DIR)/Info.plist
	@echo '<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' >> $(BUNDLE_DIR)/Info.plist
	@echo '<plist version="1.0">' >> $(BUNDLE_DIR)/Info.plist
	@echo '<dict>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <key>CFBundleExecutable</key>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <string>$(BINARY_NAME)</string>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <key>CFBundleIconFile</key>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <string>icon.icns</string>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <key>CFBundleIdentifier</key>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <string>com.smartclipboard.app</string>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <key>CFBundleName</key>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <string>$(APP_NAME)</string>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <key>CFBundleVersion</key>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <string>1.0.0</string>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <key>CFBundleShortVersionString</key>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <string>1.0.0</string>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <key>NSHighResolutionCapable</key>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <true/>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <key>LSUIElement</key>' >> $(BUNDLE_DIR)/Info.plist
	@echo '  <true/>' >> $(BUNDLE_DIR)/Info.plist
	@echo '</dict>' >> $(BUNDLE_DIR)/Info.plist
	@echo '</plist>' >> $(BUNDLE_DIR)/Info.plist

	@# Копируем иконку, если она существует
	@if [ -f "$(ICON_PATH)" ]; then \
		echo "$(YELLOW)Конвертируем иконку в формат .icns...$(NC)"; \
		mkdir -p $(BUILD_DIR)/icon.iconset; \
		sips -z 16 16     $(ICON_PATH) --out $(BUILD_DIR)/icon.iconset/icon_16x16.png > /dev/null; \
		sips -z 32 32     $(ICON_PATH) --out $(BUILD_DIR)/icon.iconset/icon_16x16@2x.png > /dev/null; \
		sips -z 32 32     $(ICON_PATH) --out $(BUILD_DIR)/icon.iconset/icon_32x32.png > /dev/null; \
		sips -z 64 64     $(ICON_PATH) --out $(BUILD_DIR)/icon.iconset/icon_32x32@2x.png > /dev/null; \
		sips -z 128 128   $(ICON_PATH) --out $(BUILD_DIR)/icon.iconset/icon_128x128.png > /dev/null; \
		sips -z 256 256   $(ICON_PATH) --out $(BUILD_DIR)/icon.iconset/icon_128x128@2x.png > /dev/null; \
		sips -z 256 256   $(ICON_PATH) --out $(BUILD_DIR)/icon.iconset/icon_256x256.png > /dev/null; \
		sips -z 512 512   $(ICON_PATH) --out $(BUILD_DIR)/icon.iconset/icon_256x256@2x.png > /dev/null; \
		sips -z 512 512   $(ICON_PATH) --out $(BUILD_DIR)/icon.iconset/icon_512x512.png > /dev/null; \
		iconutil -c icns $(BUILD_DIR)/icon.iconset -o $(BUNDLE_DIR)/Resources/icon.icns; \
		rm -rf $(BUILD_DIR)/icon.iconset; \
	else \
		echo "$(YELLOW)Иконка не найдена по пути $(ICON_PATH). Используется стандартная иконка.$(NC)"; \
	fi

	@# Делаем бинарник исполняемым
	@chmod +x $(BUNDLE_DIR)/MacOS/$(BINARY_NAME)

	@echo "$(GREEN)Создан .app бандл: $(BUILD_DIR)/$(APP_NAME).app$(NC)"

# Установка приложения в /Applications
install: mac-bundle
	@echo "$(GREEN)Установка приложения в /Applications...$(NC)"
	@rm -rf /Applications/$(APP_NAME).app 2>/dev/null || true
	@cp -R $(BUILD_DIR)/$(APP_NAME).app /Applications/
	@echo "$(GREEN)Приложение установлено в /Applications/$(APP_NAME).app$(NC)"

# Очистка
clean:
	@echo "$(YELLOW)Очистка...$(NC)"
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)Папка $(BUILD_DIR) удалена$(NC)"

# Помощь
help:
	@echo "$(GREEN)Доступные команды:$(NC)"
	@echo "  build       - Сборка приложения в папку $(BUILD_DIR)"
	@echo "  mac-bundle  - Создание .app бандла для macOS"
	@echo "  install     - Установка приложения в /Applications"
	@echo "  clean       - Удаление папки $(BUILD_DIR)"
	@echo "  run         - Запуск приложения"
	@echo "  deps        - Установка зависимостей"
	@echo "  help        - Показать это сообщение"
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
