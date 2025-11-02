//go:build darwin
// +build darwin

package config

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	launchAgentDir  = "Library/LaunchAgents"
	launchAgentPlist = "com.smartclipboard.plist"
)

func updateStartupSettings(enable bool) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	launchAgentPath := filepath.Join(homeDir, launchAgentDir, launchAgentPlist)

	if !enable {
		// Remove the launch agent if it exists
		if _, err := os.Stat(launchAgentPath); err == nil {
			return os.Remove(launchAgentPath)
		}
		return nil
	}

	// Create launch agent directory if it doesn't exist
	launchAgentDirPath := filepath.Dir(launchAgentPath)
	if err := os.MkdirAll(launchAgentDirPath, 0755); err != nil {
		return err
	}

	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Create the launch agent plist content
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.smartclipboard</string>
    <key>ProgramArguments</key>
    <array>
        <string>` + execPath + `</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <false/>
</dict>
</plist>`

	// Write the launch agent plist file
	return os.WriteFile(launchAgentPath, []byte(plistContent), 0644)
}

// loadStartupSettings loads the current startup settings from the system
func loadStartupSettings() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	launchAgentPath := filepath.Join(homeDir, launchAgentDir, launchAgentPlist)
	if _, err := os.Stat(launchAgentPath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Check if the plist is actually a symlink to our current executable
	_, err = exec.Command("launchctl", "print", "gui/$UID/", "com.smartclipboard").Output()
	return err == nil, nil
}
