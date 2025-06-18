//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js

package notifier

import (
	"os/exec"
)

// sendNotification sends a notification using the appropriate method for the current OS
// Use terminal-notifier for macOS, see https://github.com/julienXX/terminal-notifier
func sendNotification(title, message, link string) error {
	_, err := exec.LookPath("terminal-notifier")
	if err != nil {
		return err
	}
	cmd := exec.Command(
		"terminal-notifier",
		"-title", title,
		"-message", message,
		"-open", link,
	)
	return cmd.Run()
}
