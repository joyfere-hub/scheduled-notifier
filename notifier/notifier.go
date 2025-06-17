package notifier

import (
	"os/exec"
	"runtime"
)

// sendNotification sends a notification using the appropriate method for the current OS
// windows11 e.g.: powershell -Command "[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null; [Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] > $null; $template = [Windows.UI.Notifications.ToastTemplateType]::ToastText02; $xml = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent($template); $text = $xml.GetElementsByTagName('text'); $text[0].AppendChild($xml.CreateTextNode('Title')); $text[1].AppendChild($xml.CreateTextNode('Message')); $toast = [Windows.UI.Notifications.ToastNotification]::new($xml); $notifier = [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('Scheduled Notifier'); $notifier.Show($toast)"
func sendNotification(title, message, link string) error {
	switch runtime.GOOS {
	case "darwin":
		// Use terminal-notifier for macOS, see https://github.com/julienXX/terminal-notifier
		cmd := exec.Command("terminal-notifier",
			"-title", title,
			"-message", message,
			"-open", link)
		return cmd.Run()
	case "windows":
		// Use PowerShell for Windows notifications
		script := `[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] > $null

$template = [Windows.UI.Notifications.ToastTemplateType]::ToastText02
$xml = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent($template)
$text = $xml.GetElementsByTagName("text")
$text[0].AppendChild($xml.CreateTextNode("` + title + `"))
$text[1].AppendChild($xml.CreateTextNode("` + message + `"))

$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
$notifier = [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("Scheduled Notifier")
$notifier.Show($toast)`

		cmd := exec.Command("powershell", "-Command", script)
		return cmd.Run()
	default:
		// For other operating systems, we could implement other notification methods
		// or return an error indicating that notifications are not supported
		return nil
	}
}
