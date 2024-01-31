package accountrepo

import (
	"context"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
)

type GetCommand struct {
	ID        *uint64
	Username  *string
	Password  *string
	FirstName *string
	LastName  *string
	Phone     *string
	ImagePath *string
}

type Repository interface {
	Create(ctx context.Context, model model.Account) error
	Get(ctx context.Context, cmd GetCommand) []model.Account
	Update(ctx context.Context, cmd GetCommand, model model.Account) error
	Delete(ctx context.Context, cmd GetCommand) error
}
