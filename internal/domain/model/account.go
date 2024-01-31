package model

import (
	"errors"
	"time"
)

type Account struct {
	ID        uint64    `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	FirstName string    `json:"firstname,omitempty"`
	LastName  string    `json:"lastname,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	ImagePath string    `json:"img_path,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	LastVisit time.Time `json:"last_visit,omitempty"`

	ShowPhone     ShowType `json:"show_phone,omitempty"`
	ShowImg       ShowType `json:"show_img,omitempty"`
	ShowLastVisit ShowType `json:"show_last_visit,omitempty"`
}

var ErrUsernameDuplicate = errors.New("username already exists")
var ErrPhoneDuplicate = errors.New("phone number already exists")
var ErrImagePathDuplicate = errors.New("image path already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid credentials") // for incorrect password, ...
