package model

import "time"

type Account struct {
	ID        uint64    `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	FirstName string    `json:"firstname,omitempty"`
	LastName  string    `json:"lastname,omitempty"`
	Phone     string    `json:"phone_number,omitempty"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	ImagePath string    `json:"img_path,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Status    string    `json:"status,omitempty"`
	LastVisit string    `json:"last_visit,omitempty"`
}
