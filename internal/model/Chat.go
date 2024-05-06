package model

import "time"

type Chat struct {
	ChatID     string    `json:"chatID"`
	CreateTime time.Time `json:"createTime"`
	Data       []Message `json:"data"`
}
