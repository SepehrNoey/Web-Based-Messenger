package model

import "time"

type Contact struct {
	UserID      uint64    `json:"user_id,omitempty"`
	ContactID   uint64    `json:"contact_id,omitempty"`
	ContactName string    `json:"contact_name,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
