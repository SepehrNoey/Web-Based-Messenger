package messagesql

import (
	"context"
	"errors"
	"sort"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/messagerepo"
	"gorm.io/gorm"
)

type MessageDTO struct {
	model.Message
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

func (r *Repository) Create(ctx context.Context, msg model.Message) error {
	dto := MessageDTO{Message: msg}
	result := r.db.WithContext(ctx).Create(&dto)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return model.ErrChatIdDuplicae
		}

		return result.Error
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, cmd messagerepo.GetCommand) []model.Message {
	var msgDTOs []MessageDTO

	var dto MessageDTO
	var conditions []string

	if cmd.ID != nil {
		dto.ID = *cmd.ID
		conditions = append(conditions, "ID")
	}
	if cmd.ChatID != nil {
		dto.ChatID = *cmd.ChatID
		conditions = append(conditions, "ChatID")
	}
	if cmd.SenderID != nil {
		dto.SenderID = *cmd.SenderID
		conditions = append(conditions, "SenderID")
	}

	if len(conditions) == 0 {
		if err := r.db.WithContext(ctx).Find(&msgDTOs); err.Error != nil {
			return nil
		}
	} else {
		if err := r.db.WithContext(ctx).Where(&dto, conditions).Find(&msgDTOs); err.Error != nil {
			return nil
		}
	}

	msgs := make([]model.Message, len(msgDTOs))
	for i, dto := range msgDTOs {
		msgs[i] = dto.Message
	}

	// notice that the returned list of messages must be sorted in increasing order
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID < msgs[j].ID
	})

	return msgs
}

func (r *Repository) Update(ctx context.Context, cmd messagerepo.GetCommand, model model.Message) error {
	dto := MessageDTO{Message: model}
	result := r.db.WithContext(ctx).Save(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, cmd messagerepo.GetCommand) error {
	result := r.db.WithContext(ctx).Delete(&MessageDTO{Message: model.Message{ID: *cmd.ID}})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
