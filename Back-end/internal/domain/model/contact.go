package model

import (
	"errors"
	"time"
)

type Contact struct {
	ContactID   uint64    `json:"contact_id,omitempty" gorm:"primaryKey"`
	UserID      uint64    `json:"user_id,omitempty" gorm:"primaryKey"`
	ContactName string    `json:"contact_name,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

var ErrUserIdContactIdDuplicate = errors.New("pair of (contact_id, user_id) already exists")
