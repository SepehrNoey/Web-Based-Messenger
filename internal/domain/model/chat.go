package model

import "time"

type Chat struct {
	ID        uint64    `json:"id,omitempty"`
	Members   []uint64  `json:"members,omitempty"`
	Messages  []Message `json:"messages,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
