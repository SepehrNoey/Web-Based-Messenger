package contactrepo

import (
	"context"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
)

type GetCommand struct {
	UserID      *uint64
	ContactID   *uint64
	ContactName *string
}

// notice that to get (or update, etc) a unique
// record, at least both of the UserID and ContactID must be given
type Repository interface {
	Create(ctx context.Context, model model.Contact) error
	Get(ctx context.Context, cmd GetCommand) []model.Contact
	Update(ctx context.Context, cmd GetCommand, model model.Contact) error
	Delete(ctx context.Context, cmd GetCommand) error
}
