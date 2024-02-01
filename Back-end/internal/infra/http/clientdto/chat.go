package clientdto

import (
	"time"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
)

type ChatWithContentDTO struct {
	ID        *uint64          `json:"chat_id,omitempty"`
	Members   *[]uint64        `json:"members,omitempty"`
	CreatedAt *time.Time       `json:"created_at,omitempty"`
	UpdatedAt *time.Time       `json:"updated_at,omitempty"`
	Messages  *[]model.Message `json:"messages,omitempty"`
}
