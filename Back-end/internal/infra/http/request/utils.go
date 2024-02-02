package request

import (
	"fmt"
	"strconv"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ReqStructWithToken interface {
	GetToken() string
	SetToken(token string)
}

// binds the values in body, query
func Bind(reqSt interface{}, c echo.Context) error {
	if err := c.Bind(reqSt); err != nil {
		return echo.ErrBadRequest
	}

	return nil
}

// binds the values in body, query and token in header (but not path parameters)
func BindT(reqSt ReqStructWithToken, c echo.Context) error {
	if err := c.Bind(reqSt); err != nil {
		return echo.ErrBadRequest
	}

	// bind jwt tokne in header

	tokenHeader := c.Request().Header.Get("Authorization")
	token, err := auth.ExtractTokenOfHeader(tokenHeader)
	if err != nil {
		return err
	}
	reqSt.SetToken(token)

	return nil
}

func Validate(reqSt interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(reqSt); err != nil {
		return fmt.Errorf("field validations failed: %w", err)
	}

	return nil
}

func GetUintParam(c echo.Context, name string) (uint64, error) {
	numStr := c.Param(name)
	if numStr == "" {
		return 0, echo.ErrBadRequest
	}

	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest
	}

	return num, nil
}

type TokenOnly struct {
	Token string `json:"-" validate:"required"`
}

func (t *TokenOnly) GetToken() string {
	return t.Token
}

func (t *TokenOnly) SetToken(token string) {
	t.Token = token
}
