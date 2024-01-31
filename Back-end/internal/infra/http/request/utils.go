package request

import (
	"fmt"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ReqStructWithToken interface {
	GetToken() string
	SetToken(token string)
}

// binds the values in body, query and path param
func Bind(reqSt interface{}, c echo.Context) error {
	if err := c.Bind(reqSt); err != nil {
		return echo.ErrBadRequest
	}

	return nil
}

// binds the values in body, query and path param and also the token given
// in the header of request
func BindT(reqSt ReqStructWithToken, c echo.Context) error {
	if err := c.Bind(reqSt); err != nil {
		return echo.ErrBadRequest
	}

	// bind jwt tokne in header
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, reqSt); err != nil {
		return echo.ErrBadRequest
	} else {
		token, err := auth.ExtractTokenOfHeader(reqSt.GetToken())
		if err != nil {
			return err
		}
		reqSt.SetToken(token)
	}

	return nil
}

func Validate(reqSt interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(reqSt); err != nil {
		return fmt.Errorf("field validations failed: %w", err)
	}

	return nil
}

type TokenAndID struct {
	ID    *uint64 `param:"id,omitempty" validate:"number,required"`
	Token *string `header:"Authorization,omitempty" validate:"required"`
}

func (tid *TokenAndID) GetToken() string {
	return *tid.Token
}

func (tid *TokenAndID) SetToken(token string) {
	*tid.Token = token
}

type TokenOnly struct {
	Token *string `header:"Authorization,omitempty" validate:"required"`
}

func (t *TokenOnly) GetToken() string {
	return *t.Token
}

func (t *TokenOnly) SetToken(token string) {
	*t.Token = token
}
