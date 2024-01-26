package chatrepo

import (
	"context"

	"github.com/SepehrNoey/Web-Based-Messenger/internal/domain/model"
)

type GetCommand struct {
	ID       *uint64
	Messages *[]model.Message
}

type Repository interface {
	Create(ctx context.Context, model model.Chat) error
	Get(ctx context.Context, cmd GetCommand) []model.Chat
	Delete(ctx context.Context, cmd GetCommand) error
	Update(ctx context.Context, cmd GetCommand) error
}
