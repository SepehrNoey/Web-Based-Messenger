package contactrepo

import (
	"context"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
)

type GetCommand struct {
	UserID        *uint64
	ContactID     *uint64
	ContactName   *string
	ShowPhone     *string
	ShowImg       *string
	ShowLastVisit *string
	ImgPath       *string
}

type Repository interface {
	Create(ctx context.Context, model model.Contact) error
	Get(ctx context.Context, cmd GetCommand) []model.Contact
	Update(ctx context.Context, cmd GetCommand) error
	Delete(ctx context.Context, cmd GetCommand) error
}
