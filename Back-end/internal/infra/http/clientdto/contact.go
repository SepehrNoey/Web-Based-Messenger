package clientdto

import "time"

type ContactDTO struct {
	UserID          *uint64    `json:"user_id,omitempty"`
	ContactID       *uint64    `json:"contact_id,omitempty"`
	ContactName     *string    `json:"contact_name,omitempty"`
	ContactUsername *string    `json:"contact_username,omitempty"`
	Phone           *string    `json:"phone,omitempty"`
	ImagePath       *string    `json:"img_path,omitempty"`
	Bio             *string    `json:"bio,omitempty"`
	LastVisit       *time.Time `json:"last_visit,omitempty"`
}
