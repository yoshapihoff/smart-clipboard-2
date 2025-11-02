//go:build windows
// +build windows

package config

import (
	"os"
	"path/filepath"
)

const (
	autostartFolder = "AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"
	shortcutName   = "Smart Clipboard.lnk"
)


// updateStartupSettings manages the Windows startup settings
func updateStartupSettings(enable bool) error {
	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Get the path to the user's startup folder
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	startupPath := filepath.Join(homeDir, autostartFolder)
	shortcutPath := filepath.Join(startupPath, shortcutName)

	if !enable {
		// Remove the shortcut if it exists
		if _, err := os.Stat(shortcutPath); err == nil {
			return os.Remove(shortcutPath)
		}
		return nil
	}

	// Create startup directory if it doesn't exist
	if err := os.MkdirAll(startupPath, 0755); err != nil {
		return err
	}

	// Create a Windows shortcut
	err = createShortcut(execPath, shortcutPath, "")
	if err != nil {
		return err
	}

	return nil
}

// loadStartupSettings checks if the application is set to run at startup
func loadStartupSettings() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	shortcutPath := filepath.Join(homeDir, autostartFolder, shortcutName)
	if _, err := os.Stat(shortcutPath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// createShortcut creates a Windows shortcut (.lnk file)
func createShortcut(target, shortcutPath, args string) error {
	// This is a simplified version - in a real application, you might want to use a library
	// like github.com/go-ole/go-ole/oleutil for more reliable shortcut creation
	// For now, we'll use a simple implementation that creates a basic shortcut

	// Create a simple batch file as a fallback
	batchContent := `@echo off
start "" "` + target + `" ` + args + `
`

	// Use .bat extension for the shortcut
	batchPath := shortcutPath + ".bat"
	
	// Write the batch file
	return os.WriteFile(batchPath, []byte(batchContent), 0644)
}
