package fyne_app

import (
	"fmt"
	"smart-clipboard-2/internal/config"

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
	if desk, ok := tmm.app.app.(desktop.App); ok {
		cfg := config.GetConfig()
		menuItems := make([]*fyne.MenuItem, 0)
		for _, item := range clipboardItems {
			preview := item.Preview
			if cfg.DebugMode {
				preview = fmt.Sprintf("[%d] %s", item.ClickCount, preview)
			}
			menuItems = append(menuItems, fyne.NewMenuItem(preview, func() {
				tmm.app.app.Clipboard().SetContent(item.Content)
			}))
		}
		menuItems = append(menuItems, fyne.NewMenuItemSeparator())
		menuItems = append(menuItems, fyne.NewMenuItem("Preferences", func() {
			tmm.app.ToggleWindow()
		}))
		menuItems = append(menuItems, fyne.NewMenuItem("Exit", func() {
			tmm.app.app.Quit()
		}))
		tmm.currentMenu = fyne.NewMenu("", menuItems...)

		if tmm.iconResource != nil {
			desk.SetSystemTrayIcon(tmm.iconResource)
		}

		desk.SetSystemTrayMenu(tmm.currentMenu)
	}
}
