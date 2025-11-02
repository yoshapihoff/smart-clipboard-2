//go:build linux
// +build linux

package config

import (
	"os"
	"path/filepath"
)

const (
	autostartDir = ".config/autostart"
	desktopFile  = "smart-clipboard.desktop"
)

func updateStartupSettings(enable bool) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	autostartPath := filepath.Join(homeDir, autostartDir)
	desktopFilePath := filepath.Join(autostartPath, desktopFile)

	if !enable {
		// Remove the desktop file if it exists
		if _, err := os.Stat(desktopFilePath); err == nil {
			return os.Remove(desktopFilePath)
		}
		return nil
	}

	// Create autostart directory if it doesn't exist
	if err := os.MkdirAll(autostartPath, 0755); err != nil {
		return err
	}

	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Create the desktop file content
	desktopContent := `[Desktop Entry]
Type=Application
Name=Smart Clipboard
Exec=` + execPath + `
Comment=Smart Clipboard Manager
X-GNOME-Autostart-enabled=true
NoDisplay=false
`

	// Write the desktop file
	return os.WriteFile(desktopFilePath, []byte(desktopContent), 0644)
}

func loadStartupSettings() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	desktopFilePath := filepath.Join(homeDir, autostartDir, desktopFile)
	if _, err := os.Stat(desktopFilePath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
