package main

import (
	"log"
	"smart-clipboard-2/internal/fyne_app"
)

func main() {
	app, err := fyne_app.NewApp()
	if err != nil {
		log.Fatalf("Ошибка создания приложения: %v", err)
	}
	base64Icon := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="
	iconResource := fyne_app.CreateImageResourceFromBase64(base64Icon)
	app.SetupSystemTray(iconResource)
	app.Run()
}
