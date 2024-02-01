package contactsql

import (
	"context"
	"errors"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/contactrepo"
	"gorm.io/gorm"
)

type ContactDTO struct {
	model.Contact
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

func (r *Repository) Create(ctx context.Context, contact model.Contact) error {
	dto := ContactDTO{Contact: contact}
	result := r.db.WithContext(ctx).Create(&dto)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return model.ErrUserIdContactIdDuplicate
		}

		return result.Error
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, cmd contactrepo.GetCommand) []model.Contact {
	var contactDTOs []ContactDTO

	var dto ContactDTO
	var conditions []string

	if cmd.UserID != nil {
		dto.UserID = *cmd.UserID
		conditions = append(conditions, "UserID")
	}
	if cmd.ContactID != nil {
		dto.ContactID = *cmd.ContactID
		conditions = append(conditions, "ContactID")
	}
	if cmd.ContactName != nil {
		dto.ContactName = *cmd.ContactName
		conditions = append(conditions, "ContactName")
	}

	if len(conditions) == 0 {
		if err := r.db.WithContext(ctx).Find(&contactDTOs); err.Error != nil {
			return nil
		}
	} else {
		if err := r.db.WithContext(ctx).Where(&dto, conditions).Find(&contactDTOs); err.Error != nil {
			return nil
		}
	}

	contacts := make([]model.Contact, len(contactDTOs))
	for i, dto := range contactDTOs {
		contacts[i] = dto.Contact
	}

	return contacts
}

func (r *Repository) Update(ctx context.Context, cmd contactrepo.GetCommand, model model.Contact) error {
	dto := ContactDTO{Contact: model}
	result := r.db.WithContext(ctx).Save(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, cmd contactrepo.GetCommand) error {
	result := r.db.WithContext(ctx).Delete(&ContactDTO{Contact: model.Contact{UserID: *cmd.UserID, ContactID: *cmd.ContactID}})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
