package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/repository/accountrepo"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/infra/http/auth"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/infra/http/clientdto"
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

func (ah *AccountHandler) chooseLogin(c echo.Context) error {
	var req request.AccountLoginWholeData

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if req.Phone == "" {
		return ah.LoginByUsername(c)
	} else if req.Username == "" {
		return ah.LoginByPhone(c)
	} else {
		return echo.ErrBadRequest
	}
}

func (ah *AccountHandler) GetUserInfo(c echo.Context) error {
	var req request.AccountRequestById
	var claims map[string]interface{}

	// bind fields path params
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	// bind fields in header (jwt token)
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &req); err != nil {
		return echo.ErrBadRequest
	} else {
		parts := strings.Split(req.Token, " ")
		if len(parts) != 2 {
			return echo.ErrBadRequest
		} else if parts[0] != "Bearer" {
			return echo.ErrBadRequest
		} else {
			req.Token = parts[1]
			if claims, err = ah.jwtConfig.ValidateToken(req.Token); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
		}
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accounts := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		ID: &req.ID,
	})
	if len(accounts) > 1 {
		return echo.ErrInternalServerError
	} else if len(accounts) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrUserNotFound.Error())
	}

	account := accounts[0]
	if claims["id"] != account.ID { // notice that still anybody that have the jwt token, can modify/access the data of the real user
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	accountDTO := clientdto.Account{ // we don't send some fields like password, etc...
		ID:        account.ID,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Phone:     account.Phone,
		Username:  account.Username,
		ImagePath: account.ImagePath,
		Bio:       account.Bio,
		Status:    account.Status,
		LastVisit: account.LastVisit,
	}
	return c.JSON(http.StatusOK, accountDTO)

}

func (ah *AccountHandler) UpdateUserInfo(c echo.Context) error {
	var req request.AccountUpdate
	var claims map[string]interface{}

	// bind fields path params and body
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	// bind fields in header (jwt token)
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &req); err != nil {
		return echo.ErrBadRequest
	} else {
		parts := strings.Split(*req.Token, " ")
		if len(parts) != 2 {
			return echo.ErrBadRequest
		} else if parts[0] != "Bearer" {
			return echo.ErrBadRequest
		} else {
			*req.Token = parts[1]
			if claims, err = ah.jwtConfig.ValidateToken(*req.Token); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
		}
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accounts := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		ID: req.ID,
	})
	if len(accounts) > 1 {
		return echo.ErrInternalServerError
	} else if len(accounts) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrUserNotFound.Error())
	}

	account := accounts[0]
	if claims["id"] != account.ID {
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	updatedModel := model.Account{ID: *req.ID, LastVisit: time.Now()}
	if req.Firstname != nil {
		updatedModel.FirstName = *req.Firstname
	} else {
		updatedModel.FirstName = account.FirstName
	}

	if req.Lastname != nil {
		updatedModel.LastName = *req.Lastname
	} else {
		updatedModel.LastName = account.LastName
	}

	if req.Phone != nil {
		existingAccs := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{Phone: req.Phone})
		if len(existingAccs) > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, model.ErrPhoneDuplicate.Error())
		}

		updatedModel.Phone = *req.Phone
	} else {
		updatedModel.Phone = account.Phone
	}

	if req.Username != nil {
		existingAccs := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{Username: req.Username})
		if len(existingAccs) > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, model.ErrUsernameDuplicate.Error())
		}

		updatedModel.Username = *req.Username
	} else {
		updatedModel.Username = account.Username
	}

	if req.Password != nil {
		updatedModel.Password = *req.Password
	} else {
		updatedModel.Password = account.Password
	}

	if req.ImagePath != nil {
		existingAccs := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{ImagePath: req.ImagePath})
		if len(existingAccs) > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, model.ErrImagePathDuplicate.Error())
		}

		updatedModel.ImagePath = *req.ImagePath
	} else {
		updatedModel.ImagePath = account.ImagePath
	}

	if req.Bio != nil {
		updatedModel.Bio = *req.Bio
	} else {
		updatedModel.Bio = account.Bio
	}

	if err := ah.repo.Update(c.Request().Context(), accountrepo.GetCommand{ID: &account.ID}, updatedModel); err != nil {
		return echo.ErrInternalServerError
	}

	accountDTO := clientdto.Account{
		ID:        updatedModel.ID,
		FirstName: updatedModel.FirstName,
		LastName:  updatedModel.LastName,
		Phone:     updatedModel.Phone,
		Username:  updatedModel.Username,
		ImagePath: updatedModel.ImagePath,
		Bio:       updatedModel.Bio,
		Status:    updatedModel.Status,
		LastVisit: updatedModel.LastVisit,
	}
	return c.JSON(http.StatusOK, accountDTO)
}

func (ah *AccountHandler) Delete(c echo.Context) error {
	var req request.AccountRequestById
	var claims map[string]interface{}

	// bind fields path params and body
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	// bind fields in header (jwt token)
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &req); err != nil {
		return echo.ErrBadRequest
	} else {
		parts := strings.Split(req.Token, " ")
		if len(parts) != 2 {
			return echo.ErrBadRequest
		} else if parts[0] != "Bearer" {
			return echo.ErrBadRequest
		} else {
			req.Token = parts[1]
			if claims, err = ah.jwtConfig.ValidateToken(req.Token); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
		}
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accounts := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		ID: &req.ID,
	})
	if len(accounts) > 1 {
		return echo.ErrInternalServerError
	} else if len(accounts) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrUserNotFound.Error())
	}

	account := accounts[0]
	if claims["id"] != account.ID {
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	if err := ah.repo.Delete(c.Request().Context(), accountrepo.GetCommand{ID: &account.ID}); err != nil {
		return echo.ErrInternalServerError
	}

	// and here we may need to do other things due to deletion of user account!

	return c.NoContent(http.StatusOK)
}

func (ah *AccountHandler) Search(c echo.Context) error {
	var req request.AccountSearch

	// bind fields query params
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	// bind fields in header (jwt token)
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &req); err != nil {
		return echo.ErrBadRequest
	} else {
		parts := strings.Split(*req.Token, " ")
		if len(parts) != 2 {
			return echo.ErrBadRequest
		} else if parts[0] != "Bearer" {
			return echo.ErrBadRequest
		} else {
			*req.Token = parts[1]
			if _, err = ah.jwtConfig.ValidateToken(*req.Token); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
		}
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accounts := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		ID:        req.ID,
		Username:  req.Username,
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Phone:     req.Phone,
	})

	// here we can again create accountDTO array, but our current
	// accounts array is a subset of accountDTO, so we just return it
	return c.JSON(http.StatusOK, accounts)
}

func (ah *AccountHandler) RegisterMethods(g *echo.Group) {
	g.POST("register", ah.Register)
	g.POST("login", ah.chooseLogin)
	g.GET("users/:id", ah.GetUserInfo)
	g.PATCH("users/:id", ah.UpdateUserInfo)
	g.DELETE("users/:id", ah.Delete)
	g.GET("users", ah.Search)
}
