package model

import (
	"errors"
	"time"
)

type Message struct {
	ID        uint64    `json:"message_id,omitempty" gorm:"primaryKey"`
	ChatID    uint64    `json:"chat_id,omitempty"`
	SenderID  uint64    `json:"sender_id,omitempty"`
	Content   string    `json:"content,omitempty"`
	Seen      bool      `json:"seen,omitempty"`       // added for seen features
	CreatedAt time.Time `json:"created_at,omitempty"` // sending time
	UpdatedAt time.Time `json:"updated_at,omitempty"` // added this for edit a message
}

type MessageConfig int

const (
	MaxMsgSize MessageConfig = 300 * 2 // persian characters are 2 bytes each
)

var ErrMessageNotFound = errors.New("message not found")
