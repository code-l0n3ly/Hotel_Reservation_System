package model

import (
	"errors"
	"time"
)

// Message represents the 'Message' table in your database.
type Message struct {
	MessageID  string    `json:"messageID"`
	ChatID     string    `json:"chatID"`
	Content    string    `json:"content"` // Assuming JSON data as a string; adjust according to your needs
	CreateTime time.Time `json:"createTime"`
	SenderID   string    `json:"senderID"`
	ReceiverID string    `json:"receiverID"`
}

func (m *Message) Validate() error {
	if m.MessageID == "" {
		return errors.New("MessageID is required")
	}
	if m.Content == "" {
		return errors.New("content is required")
	}

	if m.SenderID == "" {
		return errors.New("senderID is required")
	}
	return nil
}

func (m *Message) IsFrom(senderID string) bool {
	return m.SenderID == senderID
}
