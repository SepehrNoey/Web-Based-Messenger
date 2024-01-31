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

type ContactHandler struct {
	cntRepo   contactrepo.Repository
	accRepo   accountrepo.Repository
	jwtConfig auth.JWTConfig
}

func NewContactHandler(cntRepo contactrepo.Repository, accRepo accountrepo.Repository, jwtConfig auth.JWTConfig) *ContactHandler {
	return &ContactHandler{
		cntRepo:   cntRepo,
		accRepo:   accRepo,
		jwtConfig: jwtConfig,
	}
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

	for i, cntAcc := range contactAccs {
		dto := clientdto.ContactDTO{}
		dto.UserID = &id
		dto.ContactID = &cntAcc.ID
		dto.ContactName = &contacts[i].ContactName
		dto.ContactUsername = &cntAcc.Username
		dto.Bio = &cntAcc.Bio

		// now setting fields according to the contact privacy settings
		dtoContacts := ch.cntRepo.Get(c.Request().Context(), contactrepo.GetCommand{UserID: &cntAcc.ID, ContactID: &id})
		if len(dtoContacts) > 1 {
			return echo.ErrInternalServerError
		}

		if cntAcc.ShowImg == model.All {
			dto.ImagePath = &cntAcc.ImagePath
		} else if cntAcc.ShowImg == model.ContactsOnly {
			if len(dtoContacts) == 1 {
				dto.ImagePath = &cntAcc.ImagePath
			}
		}

		if cntAcc.ShowLastVisit == model.All {
			dto.LastVisit = &cntAcc.LastVisit
		} else if cntAcc.ShowLastVisit == model.ContactsOnly {
			if len(dtoContacts) == 1 {
				dto.LastVisit = &cntAcc.LastVisit
			}
		}

		if cntAcc.ShowPhone == model.All {
			dto.Phone = &cntAcc.Phone
		} else if cntAcc.ShowPhone == model.ContactsOnly {
			if len(dtoContacts) == 1 {
				dto.Phone = &cntAcc.Phone
			}
		}

		contactDTOs = append(contactDTOs, dto)
	}

	return c.JSON(http.StatusOK, contactDTOs)

}
