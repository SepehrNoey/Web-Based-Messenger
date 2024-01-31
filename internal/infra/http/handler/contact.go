package handler

import (
	"net/http"

	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/repository/accountrepo"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/repository/contactrepo"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/infra/http/auth"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/infra/http/clientdto"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/infra/http/request"
	"github.com/labstack/echo/v4"
)

var handler ContactHandler

func IsContactOf(c echo.Context, cntID uint64, userID uint64) (bool, *model.Contact, error) {
	contacts := handler.cntRepo.Get(c.Request().Context(), contactrepo.GetCommand{
		UserID:    &userID,
		ContactID: &cntID,
	})

	if len(contacts) > 1 {
		return false, nil, echo.ErrInternalServerError
	} else if len(contacts) == 0 {
		return false, nil, nil
	} else {
		return true, &contacts[0], nil
	}
}

func GetContactDTOPrivacyConsidered(c echo.Context, cntAcc *model.Account, userID *uint64) (*clientdto.ContactDTO, error) {
	dto := clientdto.ContactDTO{}
	isContact, contact, err := IsContactOf(c, cntAcc.ID, *userID)
	if err != nil {
		return nil, err
	}

	dto.UserID = userID
	dto.ContactID = &cntAcc.ID
	dto.ContactName = &contact.ContactName
	dto.ContactUsername = &cntAcc.Username
	dto.Bio = &cntAcc.Bio

	// now setting fields according to the contact privacy settings
	if cntAcc.ShowImg == model.All {
		dto.ImagePath = &cntAcc.ImagePath
	} else if cntAcc.ShowImg == model.ContactsOnly {
		if isContact {
			dto.ImagePath = &cntAcc.ImagePath
		}
	}

	if cntAcc.ShowLastVisit == model.All {
		dto.LastVisit = &cntAcc.LastVisit
	} else if cntAcc.ShowLastVisit == model.ContactsOnly {
		if isContact {
			dto.LastVisit = &cntAcc.LastVisit
		}
	}

	if cntAcc.ShowPhone == model.All {
		dto.Phone = &cntAcc.Phone
	} else if cntAcc.ShowPhone == model.ContactsOnly {
		if isContact {
			dto.Phone = &cntAcc.Phone
		}
	}

	return &dto, nil
}

type ContactHandler struct {
	cntRepo   contactrepo.Repository
	accRepo   accountrepo.Repository
	jwtConfig auth.JWTConfig
}

func NewContactHandler(cntRepo contactrepo.Repository, accRepo accountrepo.Repository, jwtConfig auth.JWTConfig) *ContactHandler {
	handler = ContactHandler{
		cntRepo:   cntRepo,
		accRepo:   accRepo,
		jwtConfig: jwtConfig,
	}

	return &handler
}

func (ch *ContactHandler) Get(c echo.Context) error {
	var req request.TokenAndID

	if err := request.BindT(&req, c); err != nil {
		return err
	}

	var claims map[string]interface{}
	var err error
	if claims, err = ch.jwtConfig.ValidateToken(*req.Token); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	id, _ := claims["id"].(uint64)
	if id != *req.ID {
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	contacts := ch.cntRepo.Get(c.Request().Context(), contactrepo.GetCommand{
		UserID: &id,
	})

	var contactAccs []model.Account
	for _, cnt := range contacts {
		accs := ch.accRepo.Get(c.Request().Context(), accountrepo.GetCommand{ID: &cnt.ContactID})
		if len(accs) != 1 {
			return echo.ErrInternalServerError
		}

		contactAccs = append(contactAccs, accs[0])
	}

	var contactDTOs []clientdto.ContactDTO
	for _, cntAcc := range contactAccs {
		dto, err := GetContactDTOPrivacyConsidered(c, &cntAcc, &id)
		if err != nil {
			return err
		}

		contactDTOs = append(contactDTOs, *dto)
	}

	return c.JSON(http.StatusOK, contactDTOs)

}

func (ch *ContactHandler) Create(c echo.Context) error {
	var req request.ContactCreate

	if err := request.BindT(&req, c); err != nil {
		return err
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := ch.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userID, _ := claims["id"].(uint64)
	if userID != *req.UserID {
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	if err := ch.cntRepo.Create(c.Request().Context(), model.Contact{
		UserID:      *req.UserID,
		ContactID:   *req.ContactID,
		ContactName: *req.ContactName,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}

func (ch *ContactHandler) Delete(c echo.Context) error {
	var req request.ContactDelete

	if err := request.BindT(&req, c); err != nil {
		return err
	}
	if err := request.Validate(&req); err != nil {
		return err
	}

	claims, err := ch.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userID, _ := claims["id"].(uint64)
	if userID != *req.UserID {
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	if err := ch.cntRepo.Delete(c.Request().Context(), contactrepo.GetCommand{
		UserID:    req.UserID,
		ContactID: req.ContactID,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (ch *ContactHandler) RegisterMethods(g *echo.Group) {
	g.GET("users/:id/contacts", ch.Get)
	g.POST("users/:id/contacts", ch.Create)
	g.DELETE("users/:id/contacts/:contact_id", ch.Delete)
}
