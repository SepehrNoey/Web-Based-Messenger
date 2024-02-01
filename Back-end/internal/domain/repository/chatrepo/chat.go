package chatrepo

import (
	"context"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
)

type GetCommand struct {
	ID      *uint64
	Members *[]uint64 // in search by members, the chat must contain all members, but order isn't important
}

type Repository interface {
	Create(ctx context.Context, model model.Chat) error
	Get(ctx context.Context, cmd GetCommand) []model.Chat
	Update(ctx context.Context, cmd GetCommand, model model.Chat) error
	Delete(ctx context.Context, cmd GetCommand) error
}
