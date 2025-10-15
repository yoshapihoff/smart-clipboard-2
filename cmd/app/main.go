package main

import (
	"log"
	"smart-clipboard-2/internal/fyne_app"
)

func main() {
	// Создаем приложение
	app, err := fyne_app.NewApp()
	if err != nil {
		log.Fatalf("Ошибка создания приложения: %v", err)
	}

	// Base64 данные простой иконки
	base64Icon := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="

	// Создаем ресурс изображения для трея
	iconResource := fyne_app.CreateImageResourceFromBase64(base64Icon)

	// Настраиваем системный трей
	fyne_app.SetupSystemTray(app, iconResource)

	// Запускаем мониторинг буфера обмена
	app.Clipboard.SetCallback(func(content string) {
		log.Printf("Новый контент в буфере обмена: %s", content)
	})
	app.Clipboard.Start()

	// Запускаем приложение
	app.Run()
}
