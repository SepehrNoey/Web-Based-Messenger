package messagerepo

import (
	"context"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
)

type GetCommand struct {
	ID       *uint64
	ChatID   *uint64
	SenderID *uint64
}

type Repository interface {
	Create(ctx context.Context, model model.Message) error
	Get(ctx context.Context, cmd GetCommand) []model.Message
	Update(ctx context.Context, cmd GetCommand, model model.Message) error
	Delete(ctx context.Context, cmd GetCommand) error
}
