package fyne_app

import (
	"smart-clipboard-2/internal/config"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *App) NewSettingsContainer() *fyne.Container {
	cfg := config.GetConfig()

	historySizeEntry := widget.NewEntry()
	historySizeEntry.SetText(strconv.Itoa(cfg.ClipboardHistorySize))

	debugCheck := widget.NewCheck("", func(checked bool) {
		cfg.DebugMode = checked
	})
	debugCheck.SetChecked(cfg.DebugMode)

	runAtStartupCheck := widget.NewCheck("", func(checked bool) {
		if err := config.SetRunAtStartup(checked); err != nil {
			dialog.ShowError(err, app.window)
		}
	})
	runAtStartup := config.GetRunAtStartup()
	runAtStartupCheck.SetChecked(runAtStartup)

	saveButton := widget.NewButtonWithIcon("Save Settings", theme.DocumentSaveIcon(), func() {
		fyne.Do(func() {
			if size, err := strconv.Atoi(historySizeEntry.Text); err == nil && size > 0 {
				config.SetClipboardHistorySize(size)
			}

			if err := config.SaveConfig(); err != nil {
				dialog.ShowError(err, app.window)
			} else {
				app.HideWindow()
			}
		})
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Clipboard History Size:", Widget: historySizeEntry},
			{Text: "Debug Mode:", Widget: debugCheck},
			{Text: "Run at Startup:", Widget: runAtStartupCheck, HintText: "Start app automatically when you log in"},
		},
	}

	settingsContainer := container.NewBorder(
		nil,
		container.NewVBox(
			widget.NewSeparator(),
			container.NewHBox(
				layout.NewSpacer(),
				saveButton,
			),
		),
		nil,
		nil,
		container.NewVScroll(form),
	)

	return settingsContainer
}
