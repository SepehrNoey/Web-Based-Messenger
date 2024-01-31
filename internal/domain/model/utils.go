package model

import "errors"

type ShowType string

const (
	All          ShowType = "All"
	Noone        ShowType = "Noone"
	ContactsOnly ShowType = "Contacts-Only"
)

type Status string

const (
	Online  Status = "Online"
	Offline Status = "Offline"
)

var ErrAccessForbidden = errors.New("access to private resources is forbidden")
