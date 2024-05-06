package model

import "time"

type Chat struct {
	ChatID     string    `json:"chatID"`
	SenderID   string    `json:"senderID"`
	ReceiverID string    `json:"receiverID"`
	CreateTime time.Time `json:"createTime"`
	Messages   []Message `json:"data"`
}
