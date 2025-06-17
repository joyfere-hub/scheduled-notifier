package notifier

import "fmt"

type Message struct {
	Title   string
	Message string
	Link    string
}

func (m *Message) String() string {
	return fmt.Sprintf("Title: %s, Message: %s, Link: %s", m.Title, m.Message, m.Link)
}

func (m *Message) Send() error {
	return sendNotification(m.Title, m.Message, m.Link)
}
