package fyne_app

import (
	"time"

	"fyne.io/fyne/v2"
)

// ClipboardMonitor структура для отслеживания изменений буфера обмена
type ClipboardMonitor struct {
	app           *App
	lastContent   string
	checkInterval time.Duration
	stopChan      chan bool
	isRunning     bool
	callback      func(string) // callback функция для обработки новых данных
}

// NewClipboardMonitor создает новый монитор буфера обмена
func NewClipboardMonitor(app *App) *ClipboardMonitor {
	return &ClipboardMonitor{
		app:           app,
		checkInterval: time.Millisecond * 100, // Проверяем каждую секунду
		stopChan:      make(chan bool),
		isRunning:     false,
	}
}

// SetCheckInterval устанавливает интервал проверки буфера обмена
func (cm *ClipboardMonitor) SetCheckInterval(interval time.Duration) {
	cm.checkInterval = interval
}

// SetCallback устанавливает функцию обратного вызова для новых данных
func (cm *ClipboardMonitor) SetCallback(callback func(string)) {
	cm.callback = callback
}

// Start запускает мониторинг буфера обмена
func (cm *ClipboardMonitor) Start() {
	if cm.isRunning {
		return
	}

	cm.isRunning = true
	go cm.monitorLoop()
}

// Stop останавливает мониторинг буфера обмена
func (cm *ClipboardMonitor) Stop() {
	if !cm.isRunning {
		return
	}

	cm.isRunning = false
	cm.stopChan <- true
}

// IsRunning возвращает статус мониторинга
func (cm *ClipboardMonitor) IsRunning() bool {
	return cm.isRunning
}

// monitorLoop основной цикл мониторинга
func (cm *ClipboardMonitor) monitorLoop() {
	ticker := time.NewTicker(cm.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.stopChan:
			return
		case <-ticker.C:
			cm.checkClipboard()
		}
	}
}

// checkClipboard проверяет содержимое буфера обмена
func (cm *ClipboardMonitor) checkClipboard() {
	// Получаем содержимое буфера обмена через Fyne
	content := fyne.CurrentApp().Clipboard().Content()

	// Если содержимое изменилось, вызываем callback
	if content != "" && content != cm.lastContent {
		cm.lastContent = content
		if cm.callback != nil {
			cm.callback(content)
		}
	}
}
