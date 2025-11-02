package fyne_app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type App struct {
	app              fyne.App
	window           fyne.Window
	windowVisible    bool
	clipboardMonitor *ClipboardMonitor
	trayMenuManager  *TrayMenuManager
	clipboardManager *ClipboardManager
}

func NewApp() (*App, error) {
	myApp := app.New()
	window := myApp.NewWindow("Smart Clipboard Manager")

	app := &App{
		app:           myApp,
		window:        window,
		windowVisible: false,
	}

	window.SetMaster()
	window.SetCloseIntercept(func() {
		app.HideWindow()
	})

	app.clipboardMonitor = NewClipboardMonitor(app)

	app.clipboardMonitor.SetCallback(func(content string) {
		fmt.Println(content)
		app.clipboardManager.AddItem(content)
		app.UpdateSystemTray()
	})

	app.trayMenuManager = NewTrayMenuManager(app)
	app.clipboardManager = NewClipboardManager()
	app.createWindowUI()
	return app, nil
}

func (a *App) ShowWindow() {
	a.window.Show()
	a.windowVisible = true
}

func (a *App) HideWindow() {
	a.window.Hide()
	a.windowVisible = false
}

func (a *App) ToggleWindow() {
	if a.windowVisible {
		a.HideWindow()
	} else {
		a.ShowWindow()
	}
}

func (a *App) Run() {
	a.clipboardMonitor.Start()
	a.app.Run()
}

func (a *App) UpdateSystemTray() {
	a.trayMenuManager.UpdateTrayMenu(a.clipboardManager.items)
}

func (a *App) SetupSystemTray(iconResource *fyne.StaticResource) {
	a.trayMenuManager.SetIcon(iconResource)
	a.UpdateSystemTray()
}

func (a *App) createWindowUI() {
	a.window.SetContent(a.NewSettingsContainer())
	a.window.Resize(fyne.NewSize(500, 400))
	a.window.CenterOnScreen()
}
