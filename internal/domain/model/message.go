package model

import "time"

type Message struct {
	ID         uint64    `json:"id,omitempty"`
	ChatID     uint64    `json:"chat_id,omitempty"`
	SenderID   uint64    `json:"sender_id,omitempty"`
	ReceiverID uint64    `json:"receiver_id,omitempty"`
	Content    string    `json:"content,omitempty"`
	Seen       bool      `json:"seen,omitempty"`       // added for seen features
	CreatedAt  time.Time `json:"created_at,omitempty"` // sending time
	UpdatedAt  time.Time `json:"updated_at,omitempty"` // added this for edit a message
}
