package fyne_app

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"smart-clipboard-2/internal/config"
)

// Структура для управления приложением
type App struct {
	App           fyne.App
	Window        fyne.Window
	WindowVisible bool
	Clipboard     *ClipboardMonitor
}

// Функция для создания нового приложения
func NewApp() (*App, error) {
	myApp := app.New()
	myWindow := myApp.NewWindow("Smart Clipboard Manager")

	app := &App{
		App:           myApp,
		Window:        myWindow,
		WindowVisible: false,
	}

	// Настройки окна для убирания кнопок управления
	myWindow.SetMaster()        // Делает окно основным, убирая некоторые кнопки управления
	myWindow.SetFixedSize(true) // Запрещает изменение размера окна

	// Перехватываем событие закрытия окна - скрываем вместо закрытия
	myWindow.SetCloseIntercept(func() {
		app.HideWindow()
	})

	// Создаем монитор буфера обмена
	app.Clipboard = NewClipboardMonitor(app)

	// Создаем интерфейс
	app.createUI()

	return app, nil
}

// Функция для показа окна
func (a *App) ShowWindow() {
	a.Window.Show()
	a.WindowVisible = true
}

// Функция для скрытия окна
func (a *App) HideWindow() {
	a.Window.Hide()
	a.WindowVisible = false
}

// Функция для переключения видимости окна
func (a *App) ToggleWindow() {
	if a.WindowVisible {
		a.HideWindow()
	} else {
		a.ShowWindow()
	}
}

// Функция для запуска приложения
func (a *App) Run() {
	a.App.Run()
}

// Функция для создания интерфейса приложения
func (a *App) createUI() {
	// Получаем текущую конфигурацию
	cfg := config.GetConfig()

	// Создаем элементы интерфейса для настройки конфигурации
	titleLabel := widget.NewLabel("Application Settings")
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Поле для размера истории буфера обмена
	historySizeEntry := widget.NewEntry()
	historySizeEntry.SetText(strconv.Itoa(cfg.ClipboardHistorySize))

	// Чекбокс для режима дебага
	debugCheck := widget.NewCheck("", func(checked bool) {
		cfg.DebugMode = checked
	})

	// Устанавливаем начальное значение чекбокса
	debugCheck.SetChecked(cfg.DebugMode)

	// Кнопка сохранения
	saveButton := widget.NewButton("Save Settings", func() {
		// Валидация и сохранение размера истории
		if size, err := strconv.Atoi(historySizeEntry.Text); err == nil && size > 0 {
			config.SetClipboardHistorySize(size)
		}

		// Сохраняем конфигурацию в файл
		if err := config.SaveConfig(); err != nil {
			dialog.ShowError(err, a.Window)
		} else {
			a.HideWindow()
		}
	})

	// Создаем форму с лейблами слева от элементов управления
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Clipboard History Size:", Widget: historySizeEntry},
			{Text: "Debug Mode:", Widget: debugCheck},
		},
	}

	// Создаем контейнер с привязкой кнопки к нижнему краю
	content := container.NewBorder(
		container.NewVBox( // верх - заголовок и разделитель
			titleLabel,
			widget.NewSeparator(),
		),
		saveButton, // низ - кнопка сохранения
		nil,        // лево
		nil,        // право
		form,       // центр - форма с настройками
	)

	a.Window.SetContent(content)
	a.Window.Resize(fyne.NewSize(400, 300))
	a.Window.CenterOnScreen() // Центрируем окно на экране
}
