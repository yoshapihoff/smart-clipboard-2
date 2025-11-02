package fyne_app

import (
	"strings"
	"time"
)

type ClipboardMonitor struct {
	app           *App
	checkInterval time.Duration
	stopChan      chan bool
	isRunning     bool
	callback      func(string)
	lastContent   string
}

func NewClipboardMonitor(app *App) *ClipboardMonitor {
	return &ClipboardMonitor{
		app:           app,
		checkInterval: time.Millisecond * 100,
		stopChan:      make(chan bool),
		isRunning:     false,
		lastContent:   "",
	}
}

func (cm *ClipboardMonitor) SetCheckInterval(interval time.Duration) {
	cm.checkInterval = interval
}

func (cm *ClipboardMonitor) SetCallback(callback func(string)) {
	cm.callback = callback
}

func (cm *ClipboardMonitor) Start() {
	if cm.isRunning {
		return
	}

	cm.isRunning = true
	go cm.monitorLoop()
}

func (cm *ClipboardMonitor) Stop() {
	if !cm.isRunning {
		return
	}

	cm.isRunning = false
	cm.stopChan <- true
}

func (cm *ClipboardMonitor) IsRunning() bool {
	return cm.isRunning
}

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

func (cm *ClipboardMonitor) checkClipboard() {
	content := strings.TrimSpace(cm.app.app.Clipboard().Content())

	if content != "" && content != cm.lastContent && len(content) > 0 {
		cm.lastContent = content
		if cm.callback != nil {
			cm.callback(content)
		}
	}
}
