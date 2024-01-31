package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/accountrepo"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/auth"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/clientdto"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/request"
	"github.com/labstack/echo/v4"
)

var lastRegisteredAccountID = 0

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

	if err := request.Bind(&req, c); err != nil {
		return echo.ErrBadRequest
	}

	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accountsSameUsername := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		Username: req.Username,
	})
	if len(accountsSameUsername) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrUsernameDuplicate.Error())
	}

	accountsSamePhone := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		Phone: req.Phone,
	})
	if len(accountsSamePhone) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrPhoneDuplicate.Error())
	}

	accountsSameImgPath := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		ImagePath: req.ImagePath,
	})
	if len(accountsSameImgPath) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrImagePathDuplicate.Error())
	}

	if err := ah.repo.Create(c.Request().Context(), model.Account{
		ID:            uint64(lastRegisteredAccountID + 1),
		FirstName:     *req.Firstname,
		LastName:      *req.Lastname,
		Phone:         *req.Phone,
		Username:      *req.Username,
		Password:      *req.Password,
		ImagePath:     *req.ImagePath,
		Bio:           *req.Bio,
		LastVisit:     time.Now(),
		ShowPhone:     model.ContactsOnly,
		ShowImg:       model.All,
		ShowLastVisit: model.All,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	lastRegisteredAccountID++
	return c.JSON(http.StatusCreated, fmt.Sprintf("id: %v", lastRegisteredAccountID))
}

func (ah *AccountHandler) LoginByUsername(c echo.Context) error {
	var req request.AccountLoginByUsername

	if err := request.Bind(&req, c); err != nil {
		return err
	}

	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accounts := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		Username: req.Username,
	})
	if len(accounts) > 1 {
		return echo.ErrInternalServerError
	} else if len(accounts) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrUserNotFound.Error())
	}

	account := accounts[0]
	if account.Password != *req.Password {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrInvalidCredentials.Error())
	}

	token, err := ah.jwtConfig.CreateToken(ah.addClaims(&account))
	if err != nil {
		return echo.ErrInternalServerError
	}

	c.Response().Header().Set(echo.HeaderAuthorization, auth.GetAuthHeaderValue(token))
	return c.NoContent(http.StatusOK)

}

func (ah *AccountHandler) LoginByPhone(c echo.Context) error {
	var req request.AccountLoginByPhone

	if err := request.Bind(&req, c); err != nil {
		return echo.ErrBadRequest
	}

	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accounts := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		Phone: req.Phone,
	})
	if len(accounts) > 1 {
		return echo.ErrInternalServerError
	} else if len(accounts) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrUserNotFound.Error())
	}

	account := accounts[0]
	if account.Password != *req.Password {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrInvalidCredentials.Error())
	}

	token, err := ah.jwtConfig.CreateToken(ah.addClaims(&account))
	if err != nil {
		return echo.ErrInternalServerError
	}

	c.Response().Header().Set(echo.HeaderAuthorization, auth.GetAuthHeaderValue(token))
	return c.NoContent(http.StatusOK)

}

func (ah *AccountHandler) chooseLogin(c echo.Context) error {
	var req request.AccountLoginWholeData

	if err := request.Bind(&req, c); err != nil {
		return err
	}

	if *req.Phone == "" {
		return ah.LoginByUsername(c)
	} else if *req.Username == "" {
		return ah.LoginByPhone(c)
	} else {
		return echo.ErrBadRequest
	}
}

func (ah *AccountHandler) GetUserInfo(c echo.Context) error {
	var req request.TokenAndID
	if err := request.BindT(&req, c); err != nil {
		return err
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var claims map[string]interface{}
	var err error
	claims, err = ah.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
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
	if claims["id"] != account.ID { // notice that still anybody that have the jwt token, can modify/access the data of the real user
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	return c.JSON(http.StatusOK, account)

}

func (ah *AccountHandler) UpdateUserInfo(c echo.Context) error {
	var req request.AccountUpdate
	if err := request.BindT(&req, c); err != nil {
		return err
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var claims map[string]interface{}
	var err error
	claims, err = ah.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
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
	// first updating unique fields
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

	if req.ImagePath != nil {
		existingAccs := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{ImagePath: req.ImagePath})
		if len(existingAccs) > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, model.ErrImagePathDuplicate.Error())
		}

		updatedModel.ImagePath = *req.ImagePath
	} else {
		updatedModel.ImagePath = account.ImagePath
	}

	// now updating non-unique fields
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

	if req.Password != nil {
		updatedModel.Password = *req.Password
	} else {
		updatedModel.Password = account.Password
	}

	if req.Bio != nil {
		updatedModel.Bio = *req.Bio
	} else {
		updatedModel.Bio = account.Bio
	}

	if req.ShowImg != nil {
		updatedModel.ShowImg = *req.ShowImg
	} else {
		updatedModel.ShowImg = account.ShowImg
	}

	if req.ShowLastVisit != nil {
		updatedModel.ShowLastVisit = *req.ShowLastVisit
	} else {
		updatedModel.ShowLastVisit = account.ShowLastVisit
	}

	if req.ShowPhone != nil {
		updatedModel.ShowPhone = *req.ShowPhone
	} else {
		updatedModel.ShowPhone = account.ShowPhone
	}

	if err := ah.repo.Update(c.Request().Context(), accountrepo.GetCommand{ID: &account.ID}, updatedModel); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, updatedModel)
}

func (ah *AccountHandler) Delete(c echo.Context) error {
	var req request.TokenAndID
	if err := request.BindT(&req, c); err != nil {
		return err
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var claims map[string]interface{}
	var err error
	claims, err = ah.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
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

	if err := ah.repo.Delete(c.Request().Context(), accountrepo.GetCommand{ID: &account.ID}); err != nil {
		return echo.ErrInternalServerError
	}

	// and here we may need to do other things due to deletion of user account!

	return c.NoContent(http.StatusOK)
}

func (ah *AccountHandler) Search(c echo.Context) error {
	var req request.AccountSearch
	if err := request.BindT(&req, c); err != nil {
		return err
	}

	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := ah.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	accounts := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{
		ID:        req.ID,
		Username:  req.Username,
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Phone:     req.Phone,
	})

	// creating contactDTO of these accounts (to avoid sending all account info)
	id, _ := claims["id"].(uint64)
	clientAcc := ah.repo.Get(c.Request().Context(), accountrepo.GetCommand{ID: &id})
	var contactDTOs []clientdto.ContactDTO
	for _, acc := range accounts {
		dto, err := GetContactDTOPrivacyConsidered(c, &acc, &clientAcc[0].ID)
		if err != nil {
			return err
		}
		contactDTOs = append(contactDTOs, *dto)
	}

	return c.JSON(http.StatusOK, contactDTOs)
}

func (ah *AccountHandler) RegisterMethods(g *echo.Group) {
	g.POST("register", ah.Register)
	g.POST("login", ah.chooseLogin)
	g.GET("users/:id", ah.GetUserInfo)
	g.PATCH("users/:id", ah.UpdateUserInfo)
	g.DELETE("users/:id", ah.Delete)
	g.GET("users", ah.Search)
}
