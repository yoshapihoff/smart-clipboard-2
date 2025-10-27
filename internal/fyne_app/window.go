package fyne_app

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"smart-clipboard-2/internal/config"
)

type App struct {
	App              fyne.App
	Window           fyne.Window
	WindowVisible    bool
	Clipboard        *ClipboardMonitor
	trayMenuManager  *TrayMenuManager
	clipboardManager *ClipboardManager
}

func NewApp() (*App, error) {
	myApp := app.New()
	myWindow := myApp.NewWindow("Smart Clipboard Manager")

	app := &App{
		App:           myApp,
		Window:        myWindow,
		WindowVisible: false,
	}

	myWindow.SetMaster()
	myWindow.SetFixedSize(true)
	myWindow.SetCloseIntercept(func() {
		app.HideWindow()
	})

	app.Clipboard = NewClipboardMonitor(app)

	app.Clipboard.SetCallback(func(content string) {
		fmt.Println(content)
		app.clipboardManager.AddItem(content)
		app.UpdateSystemTray()
	})

	app.createUI()
	app.trayMenuManager = NewTrayMenuManager(app)
	app.clipboardManager = NewClipboardManager()
	return app, nil
}

func (a *App) ShowWindow() {
	a.Window.Show()
	a.WindowVisible = true
}

func (a *App) HideWindow() {
	a.Window.Hide()
	a.WindowVisible = false
}

func (a *App) ToggleWindow() {
	if a.WindowVisible {
		a.HideWindow()
	} else {
		a.ShowWindow()
	}
}

func (a *App) Run() {
	a.Clipboard.Start()
	a.App.Run()
}

func (a *App) UpdateSystemTray() {
	a.trayMenuManager.UpdateTrayMenu(a.clipboardManager.items)
}

func (a *App) SetupSystemTray(iconResource *fyne.StaticResource) {
	a.trayMenuManager.SetIcon(iconResource)
	a.UpdateSystemTray()
}

func (a *App) createUI() {
	cfg := config.GetConfig()

	historySizeEntry := widget.NewEntry()
	historySizeEntry.SetText(strconv.Itoa(cfg.ClipboardHistorySize))

	debugCheck := widget.NewCheck("", func(checked bool) {
		cfg.DebugMode = checked
	})

	debugCheck.SetChecked(cfg.DebugMode)

	saveButton := widget.NewButton("Save Settings", func() {
		if size, err := strconv.Atoi(historySizeEntry.Text); err == nil && size > 0 {
			config.SetClipboardHistorySize(size)
		}

		if err := config.SaveConfig(); err != nil {
			dialog.ShowError(err, a.Window)
		} else {
			a.HideWindow()
		}
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Clipboard History Size:", Widget: historySizeEntry},
			{Text: "Debug Mode:", Widget: debugCheck},
		},
	}

	content := container.NewBorder(
		nil,
		saveButton,
		nil,
		nil,
		form,
	)

	a.Window.SetContent(content)
	a.Window.Resize(fyne.NewSize(400, 300))
	a.Window.CenterOnScreen()
}
