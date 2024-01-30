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
		return fmt.Errorf("phone and password are required; %w", err)
	}

	return nil
}
