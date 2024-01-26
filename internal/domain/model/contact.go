package model

import "time"

type Contact struct {
	UserID      uint64 `json:"user_id,omitempty"`
	ContactID   uint64 `json:"contact_id,omitempty"`
	ContactName string `json:"contact_name,omitempty"`

	// used string for ShowSomething to have multiple options like "nobody", "anyone", "contacts-only", ...
	ShowPhone     string    `json:"show_phone,omitempty"`
	ShowImg       string    `json:"show_img,omitempty"`
	ShowLastVisit string    `json:"show_last_visit,omitempty"`
	Status        string    `json:"status,omitempty"`
	ImgPath       string    `json:"img_path,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
