package repository

import (
	"fmt"

	"github.com/RuanScherer/codepix-go/domain/model"
	"github.com/jinzhu/gorm"
)

type PixKeyDBRepository struct {
	DB *gorm.DB
}

func (r PixKeyDBRepository) AddBank(bank *model.Bank) error {
	err := r.DB.Create(bank).Error
	return err
}

func (r PixKeyDBRepository) AddAccount(account *model.Account) error {
	err := r.DB.Create(account).Error
	return err
}

func (r PixKeyDBRepository) RegisterKey(pixKey *model.PixKey) (*model.PixKey, error) {
	err := r.DB.Create(pixKey).Error
	if err != nil {
		return nil, err
	}
	return pixKey, nil
}

func (r PixKeyDBRepository) FindKeyByKind(key string, kind string) (*model.PixKey, error) {
	var pixKey model.PixKey

	r.DB.Preload("Account.Bank").First(&pixKey, "kind = ? and key = ?", kind, key)

	if pixKey.ID == "" {
		return nil, fmt.Errorf("no key was found")
	}
	return &pixKey, nil
}

func (r PixKeyDBRepository) FindAccount(id string) (*model.Account, error) {
	var account model.Account

	r.DB.Preload("Bank").First(&account, "id = ?", id)

	if account.ID == "" {
		return nil, fmt.Errorf("no account was found")
	}
	return &account, nil
}

func (r PixKeyDBRepository) FindBank(id string) (*model.Bank, error) {
	var bank model.Bank

	r.DB.First(&bank, "id = ?", id)

	if bank.ID == "" {
		return nil, fmt.Errorf("no bank was found")
	}
	return &bank, nil
}
