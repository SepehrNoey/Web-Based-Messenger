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

// }
