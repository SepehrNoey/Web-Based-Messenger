package chatsql

import (
	"context"
	"errors"
	"sort"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/chatrepo"
	"gorm.io/gorm"
)

type ChatDTO struct {
	model.Chat
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, chat model.Chat) error {
	dto := ChatDTO{Chat: chat}
	// we sort the members to make it easier to search in them later
	sort.Slice(dto.Members, func(i int, j int) bool {
		return dto.Members[i] < dto.Members[j]
	})

	result := r.db.WithContext(ctx).Create(&dto)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return model.ErrChatIdDuplicae
		}

		return result.Error
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, cmd chatrepo.GetCommand) []model.Chat {
	var chatDTOs []ChatDTO

	var dto ChatDTO
	var conditions []string

	if cmd.ID != nil {
		dto.ID = *cmd.ID
		conditions = append(conditions, "ID")
	}
	if cmd.Members != nil {
		sort.Slice(*cmd.Members, func(i int, j int) bool {
			return (*cmd.Members)[i] < (*cmd.Members)[j]
		})
		dto.Members = *cmd.Members
		conditions = append(conditions, "Members")
	}

	if len(conditions) == 0 {
		if err := r.db.WithContext(ctx).Find(&chatDTOs); err.Error != nil {
			return nil
		}
	} else {
		if err := r.db.WithContext(ctx).Where(&dto, conditions).Find(&chatDTOs); err.Error != nil {
			return nil
		}
	}

	chats := make([]model.Chat, len(chatDTOs))
	for i, dto := range chatDTOs {
		chats[i] = dto.Chat
	}

	return chats

}

func (r *Repository) Update(ctx context.Context, cmd chatrepo.GetCommand, model model.Chat) error {
	dto := ChatDTO{Chat: model}
	sort.Slice(dto.Members, func(i int, j int) bool {
		return dto.Members[i] < dto.Members[j]
	})
	result := r.db.WithContext(ctx).Save(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, cmd chatrepo.GetCommand) error {
	result := r.db.WithContext(ctx).Delete(&ChatDTO{Chat: model.Chat{ID: *cmd.ID}})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
