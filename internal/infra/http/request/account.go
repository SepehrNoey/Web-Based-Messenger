package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type AccountCreate struct {
	Firstname string `json:"firstname,omitempty" validate:"required"`
	Phone     string `json:"phone,omitempty" validate:"e164,required"`
	Username  string `json:"username,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"min=8,required"`

	Lastname  string `json:"lastname,omitempty"`
	ImagePath string `json:"img_path,omitempty" validate:"url"`
	Bio       string `json:"bio,omitempty" validate:"max=100"`
}

func (ac AccountCreate) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(ac); err != nil {
		return fmt.Errorf("account create validation failed: %w", err)
	}

	return nil
}

type AccountLoginByUsername struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func (alu AccountLoginByUsername) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(alu); err != nil {
		return fmt.Errorf("username and password are required: %w", err)
	}

	return nil
}

type AccountLoginByPhone struct {
	Phone    string `json:"phone,omitempty" validate:"e164,required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func (alp AccountLoginByPhone) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(alp); err != nil {
		return fmt.Errorf("phone and password are required: %w", err)
	}

	return nil
}

type AccountLoginWholeData struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type AccountRequestById struct {
	ID    uint64 `param:"id,omitempty" validate:"number,required"`
	Token string `header:"Authorization,omitempty" validate:"required"`
}

func (agu AccountRequestById) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(agu); err != nil {
		return fmt.Errorf("field validations failed: %w", err)
	}

	return nil
}

type AccountUpdate struct {
	ID    *uint64 `param:"id,omitempty" validate:"number,required"`
	Token *string `header:"Authorization,omitempty" validate:"required"`

	Firstname *string `json:"firstname,omitempty"`
	Phone     *string `json:"phone,omitempty" validate:"e164"`
	Username  *string `json:"username,omitempty"`
	Password  *string `json:"password,omitempty" validate:"min=8"`
	Lastname  *string `json:"lastname,omitempty"`
	ImagePath *string `json:"img_path,omitempty" validate:"url"`
	Bio       *string `json:"bio,omitempty" validate:"max=100"`
}

func (au AccountUpdate) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(au); err != nil {
		return fmt.Errorf("account update validation failed: %w", err)
	}

	return nil
}

type AccountSearch struct {
	Token *string `header:"Authorization,omitempty" validate:"required"`

	ID        *uint64 `query:"id,omitempty" validate:"number"`
	Firstname *string `query:"firstname,omitempty"`
	Lastname  *string `query:"lastname,omitempty"`
	Phone     *string `query:"phone,omitempty" validate:"e164"`
	Username  *string `query:"username,omitempty"`
}

func (as AccountSearch) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(as); err != nil {
		return fmt.Errorf("account search validation failed: %w", err)
	}

	return nil
}
