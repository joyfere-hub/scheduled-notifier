//go:build darwin

package notifier

import (
	"os/exec"
)

const defTerminalNotifierPath = "/opt/homebrew/opt/terminal-notifier/bin/terminal-notifier"

// sendNotification sends a notification using the appropriate method for the current OS
// Use terminal-notifier for macOS, see https://github.com/julienXX/terminal-notifier
func sendNotification(title, message, link string) error {
	path, err := exec.LookPath("terminal-notifier")
	if err != nil {
		path, err = exec.LookPath(defTerminalNotifierPath)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(
		path,
		"-title", title,
		"-message", message,
		"-open", link,
	)
	return cmd.Run()
}
