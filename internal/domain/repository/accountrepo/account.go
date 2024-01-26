package accountrepo

import (
	"context"

	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/model"
)

type GetCommand struct {
	ID        *uint64
	Username  *string
	Password  *string
	FirstName *string
	LastName  *string
	Phone     *string
}

type Repository interface {
	Register(ctx context.Context, model model.Account) error
	Login(ctx context.Context, cmd GetCommand) (string, error) // token and error
	Get(ctx context.Context, cmd GetCommand) []model.Account
	Update(ctx context.Context, cmd GetCommand) error
	Delete(ctx context.Context, cmd GetCommand) error
}
