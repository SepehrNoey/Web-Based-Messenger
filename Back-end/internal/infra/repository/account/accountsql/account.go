package accountsql

import (
	"context"
	"errors"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/accountrepo"
	"gorm.io/gorm"
)

type AccountDTO struct {
	model.Account
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

func (r *Repository) Create(ctx context.Context, account model.Account) error {
	dto := AccountDTO{Account: account}
	result := r.db.WithContext(ctx).Create(&dto)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return model.ErrIdDuplicate
		}

		return result.Error
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, cmd accountrepo.GetCommand) []model.Account {
	var accountDTOs []AccountDTO

	var dto AccountDTO
	var conditions []string

	if cmd.ID != nil {
		dto.ID = *cmd.ID
		conditions = append(conditions, "ID")
	}
	if cmd.Username != nil {
		dto.Username = *cmd.Username
		conditions = append(conditions, "Username")
	}
	if cmd.Password != nil {
		dto.Password = *cmd.Password
		conditions = append(conditions, "Password")
	}
	if cmd.FirstName != nil {
		dto.FirstName = *cmd.FirstName
		conditions = append(conditions, "FirstName")
	}
	if cmd.LastName != nil {
		dto.LastName = *cmd.LastName
		conditions = append(conditions, "LastName")
	}
	if cmd.Phone != nil {
		dto.Phone = *cmd.Phone
		conditions = append(conditions, "Phone")
	}
	if cmd.ImagePath != nil {
		dto.ImagePath = *cmd.ImagePath
		conditions = append(conditions, "ImagePath")
	}

	if len(conditions) == 0 {
		if err := r.db.WithContext(ctx).Find(&accountDTOs); err.Error != nil {
			return nil
		}
	} else {
		if err := r.db.WithContext(ctx).Where(&dto, conditions).Find(&accountDTOs); err.Error != nil {
			return nil
		}
	}

	accounts := make([]model.Account, len(accountDTOs))
	for i, dto := range accountDTOs {
		accounts[i] = dto.Account
	}

	return accounts

}

func (r *Repository) Update(ctx context.Context, cmd accountrepo.GetCommand, account model.Account) error {
	dto := AccountDTO{Account: account}
	result := r.db.WithContext(ctx).Save(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, cmd accountrepo.GetCommand) error {
	result := r.db.WithContext(ctx).Delete(&AccountDTO{Account: model.Account{ID: *cmd.ID}})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
