package fyne_app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// Функция для создания системного трея с меню (обертка для обратной совместимости)
func SetupSystemTray(app *App, iconResource *fyne.StaticResource) {
	// Создаем менеджер меню
	menuManager := NewTrayMenuManager(app)

	// Добавляем стандартные пункты меню
	menuManager.AddMenuItem("settings", "Settings", func() {
		app.ToggleWindow()
	})

	// Устанавливаем иконку
	menuManager.SetIcon(iconResource)
}

// Структура для управления пунктами меню системного трея
type TrayMenuManager struct {
	app          *App
	menuItems    map[string]*fyne.MenuItem
	currentMenu  *fyne.Menu
	iconResource *fyne.StaticResource
}

// NewTrayMenuManager создает новый менеджер меню системного трея
func NewTrayMenuManager(app *App) *TrayMenuManager {
	return &TrayMenuManager{
		app:       app,
		menuItems: make(map[string]*fyne.MenuItem),
	}
}

// AddMenuItem добавляет пункт меню в системный трей
func (tmm *TrayMenuManager) AddMenuItem(id, label string, action func()) error {
	if _, exists := tmm.menuItems[id]; exists {
		return fmt.Errorf("menu item with id '%s' already exists", id)
	}

	item := fyne.NewMenuItem(label, action)
	tmm.menuItems[id] = item
	tmm.updateTrayMenu()
	return nil
}

// RemoveMenuItem удаляет пункт меню из системного трея
func (tmm *TrayMenuManager) RemoveMenuItem(id string) error {
	if _, exists := tmm.menuItems[id]; !exists {
		return fmt.Errorf("menu item with id '%s' not found", id)
	}

	delete(tmm.menuItems, id)
	tmm.updateTrayMenu()
	return nil
}

// UpdateMenuItem обновляет существующий пункт меню
func (tmm *TrayMenuManager) UpdateMenuItem(id, newLabel string, newAction func()) error {
	if _, exists := tmm.menuItems[id]; !exists {
		return fmt.Errorf("menu item with id '%s' not found", id)
	}

	// Удаляем старый пункт меню
	delete(tmm.menuItems, id)

	// Создаем новый пункт меню с обновленными данными
	newItem := fyne.NewMenuItem(newLabel, newAction)
	tmm.menuItems[id] = newItem
	tmm.updateTrayMenu()
	return nil
}

// SetIcon устанавливает иконку для системного трея
func (tmm *TrayMenuManager) SetIcon(iconResource *fyne.StaticResource) {
	tmm.iconResource = iconResource
	tmm.updateTrayMenu()
}

// updateTrayMenu обновляет меню в системном трее
func (tmm *TrayMenuManager) updateTrayMenu() {
	if desk, ok := tmm.app.App.(desktop.App); ok {
		// Создаем слайс для всех пунктов меню
		items := make([]*fyne.MenuItem, 0, len(tmm.menuItems))

		// Добавляем все пункты меню в порядке их добавления
		for _, item := range tmm.menuItems {
			items = append(items, item)
		}

		// Создаем новое меню
		tmm.currentMenu = fyne.NewMenu("", items...)

		// Устанавливаем иконку если она есть
		if tmm.iconResource != nil {
			desk.SetSystemTrayIcon(tmm.iconResource)
		}

		// Устанавливаем меню
		desk.SetSystemTrayMenu(tmm.currentMenu)
	}
}

// GetMenuItem возвращает пункт меню по ID
func (tmm *TrayMenuManager) GetMenuItem(id string) (*fyne.MenuItem, error) {
	item, exists := tmm.menuItems[id]
	if !exists {
		return nil, fmt.Errorf("menu item with id '%s' not found", id)
	}
	return item, nil
}

// ListMenuItems возвращает список всех ID пунктов меню
func (tmm *TrayMenuManager) ListMenuItems() []string {
	ids := make([]string, 0, len(tmm.menuItems))
	for id := range tmm.menuItems {
		ids = append(ids, id)
	}
	return ids
}
