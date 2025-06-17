package notifier

import (
	"fmt"
	"log"
)

type Messages []Message

type Message struct {
	Title   string
	Message string
	Link    string
}

func (m *Messages) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *Messages) Send() error {
	for _, msg := range *m {
		log.Printf("Send message: %s", msg.String())
		if err := msg.Send(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Message) String() string {
	return fmt.Sprintf("Title: %s, Message: %s, Link: %s", m.Title, m.Message, m.Link)
}

func (m *Message) Send() error {
	return sendNotification(m.Title, m.Message, m.Link)
}
