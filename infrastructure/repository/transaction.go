package repository

import (
	"fmt"

	"github.com/RuanScherer/codepix-go/domain/model"
	"github.com/jinzhu/gorm"
)

type TransactionDBRepository struct {
	DB *gorm.DB
}

func (r TransactionDBRepository) Register(transaction *model.Transaction) error {
	err := r.DB.Create(transaction).Error
	return err
}

func (r TransactionDBRepository) Save(transaction *model.Transaction) error {
	err := r.DB.Save(transaction).Error
	return err
}

func (r TransactionDBRepository) Find(id string) (*model.Transaction, error) {
	var transaction model.Transaction

	r.DB.Preload("AccountFrom.Bank").First(&transaction, "id = ?", id)

	if transaction.ID == "" {
		return nil, fmt.Errorf("no transaction was found")
	}
	return &transaction, nil
}
