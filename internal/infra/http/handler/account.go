package handler

import (
	"net/http"
	"time"

	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/repository/accountrepo"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/infra/http/auth"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/infra/http/request"
	"github.com/labstack/echo/v4"
)

var lastRegisteredID = 0

type AccountHandler struct {
	repo      accountrepo.Repository
	jwtConfig auth.JWTConfig
}

func NewAccountHandler(repo accountrepo.Repository, jwtConfig auth.JWTConfig) *AccountHandler {
	return &AccountHandler{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

func (ah *AccountHandler) addClaims(ac *model.Account) map[string]interface{} {
	claims := make(map[string]interface{})
	claims["username"] = ac.Username
	claims["id"] = ac.ID
	claims["phone"] = ac.Phone

	return claims
}

func (ah *AccountHandler) Register(c echo.Context) error {
	var req request.AccountCreate

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accountsSameUsername := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		Username: &req.Username,
	})
	if len(accountsSameUsername) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrUsernameDuplicate.Error())
	}

	accountsSamePhone := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		Phone: &req.Phone,
	})
	if len(accountsSamePhone) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrPhoneDuplicate.Error())
	}

	accountsSameImgPath := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		ImagePath: &req.ImagePath,
	})
	if len(accountsSameImgPath) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrImagePathDuplicate.Error())
	}

	if err := ah.repo.Create(c.Request().Context(), model.Account{
		ID:        uint64(lastRegisteredID + 1),
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Phone:     req.Phone,
		Username:  req.Username,
		Password:  req.Password,
		ImagePath: req.ImagePath,
		Bio:       req.Bio,
		// won't give status field
		LastVisit: time.Now(),
	}); err != nil {
		return echo.ErrInternalServerError
	}

	lastRegisteredID++
	return c.NoContent(http.StatusCreated)
}

func (ah *AccountHandler) LoginByUsername(c echo.Context) error {
	var req request.AccountLoginByUsername

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accounts := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		Username: &req.Username,
	})
	if len(accounts) > 1 {
		return echo.ErrInternalServerError
	} else if len(accounts) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrUserNotFound.Error())
	}

	account := accounts[0]
	if account.Password != req.Password {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrInvalidCredentials.Error())
	}

	token, err := ah.jwtConfig.CreateToken(ah.addClaims(&account))
	if err != nil {
		return echo.ErrInternalServerError
	}

	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+token)
	return c.NoContent(http.StatusOK)

}

func (ah *AccountHandler) LoginByPhone(c echo.Context) error {
	var req request.AccountLoginByPhone

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accounts := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		Phone: &req.Phone,
	})
	if len(accounts) > 1 {
		return echo.ErrInternalServerError
	} else if len(accounts) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrUserNotFound.Error())
	}

	account := accounts[0]
	if account.Password != req.Password {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrInvalidCredentials.Error())
	}

	token, err := ah.jwtConfig.CreateToken(ah.addClaims(&account))
	if err != nil {
		return echo.ErrInternalServerError
	}

	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+token)
	return c.NoContent(http.StatusOK)

}
