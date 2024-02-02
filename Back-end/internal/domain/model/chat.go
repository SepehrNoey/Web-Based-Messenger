package model

import (
	"errors"
	"time"
)

type Chat struct {
	ID        uint64    `json:"chat_id,omitempty" gorm:"primaryKey"`
	Members   []uint64  `json:"members,omitempty" gorm:"type:bigint[]"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

var ErrChatDuplicate = errors.New("chat already exists")
var ErrChatNotFound = errors.New("chat not found")
var ErrChatIdDuplicae = errors.New("chat id already exists")
