package clientdto

import "time"

type Account struct {
	ID        uint64    `json:"id,omitempty"`
	FirstName string    `json:"firstname,omitempty"`
	LastName  string    `json:"lastname,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Username  string    `json:"username,omitempty"`
	ImagePath string    `json:"img_path,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Status    string    `json:"status,omitempty"`
	LastVisit time.Time `json:"last_visit,omitempty"`
}
