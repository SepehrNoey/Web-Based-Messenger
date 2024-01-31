package handler

import (
	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/repository/contactrepo"
	"github.com/SepehrNoey/Web-Based-Messenger/internal/infra/http/auth"
)

type ContactHandler struct {
	repo      contactrepo.Repository
	jwtConfig auth.JWTConfig
}

func NewContactHandler(repo contactrepo.Repository, jwtConfig auth.JWTConfig) *ContactHandler {
	return &ContactHandler{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

// func (ch *ContactHandler) Get(c echo.Context) error {
// 	var req request.TokenAndID

// 	if err := request.BindT(&req, c); err != nil {
// 		return err
// 	}

// 	var claims map[string]interface{}
// 	var err error
// 	if claims, err = ch.jwtConfig.ValidateToken(*req.Token); err != nil {
// 		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
// 	}

// 	id, _ := claims["id"].(uint64)
// 	contacts := ch.repo.Get(c.Request().Context(), contactrepo.GetCommand{
// 		UserID: &id,
// 	})

// 	// shouldn't the contact model have the Phone field?

// }
