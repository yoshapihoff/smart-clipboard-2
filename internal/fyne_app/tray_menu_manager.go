package fyne_app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type TrayMenuManager struct {
	app          *App
	currentMenu  *fyne.Menu
	iconResource *fyne.StaticResource
}

func NewTrayMenuManager(app *App) *TrayMenuManager {
	return &TrayMenuManager{
		app: app,
	}
}

func (tmm *TrayMenuManager) SetIcon(iconResource *fyne.StaticResource) {
	tmm.iconResource = iconResource
}

func (tmm *TrayMenuManager) UpdateTrayMenu(clipboardItems []ClipboardItem) {
	if desk, ok := tmm.app.App.(desktop.App); ok {
		menuItems := make([]*fyne.MenuItem, 0)
		for _, item := range clipboardItems {
			if item.Content == "" {
				continue
			}
			menuItems = append(menuItems, fyne.NewMenuItem(item.Preview, func() {
				tmm.app.App.Clipboard().SetContent(item.Content)
			}))
		}
		menuItems = append(menuItems, fyne.NewMenuItemSeparator())
		menuItems = append(menuItems, fyne.NewMenuItem("Settings", func() {
			tmm.app.ToggleWindow()
		}))
		menuItems = append(menuItems, fyne.NewMenuItem("Exit", func() {
			tmm.app.App.Quit()
		}))
		tmm.currentMenu = fyne.NewMenu("", menuItems...)

		if tmm.iconResource != nil {
			desk.SetSystemTrayIcon(tmm.iconResource)
		}

		desk.SetSystemTrayMenu(tmm.currentMenu)
	}
}
