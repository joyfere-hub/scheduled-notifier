//go:build windows

package notifier

import "git.sr.ht/~jackmordaunt/go-toast"

// sendNotification sends a notification using the appropriate method for the current OS
// Use go-toast for Windows, see https://git.sr.ht/~jackmordaunt/go-toast
func sendNotification(title, message, link string) error {
	n := toast.Notification{
		AppID:               "Scheduled Notifier",
		Title:               title,
		Body:                message,
		ActivationArguments: link,
	}
	return n.Push()
}
