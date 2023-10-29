package model

import (
	"time"

	"github.com/asaskevich/govalidator"

	uuid "github.com/satori/go.uuid"
)

type Bank struct {
	Base     `valid:"required"`
	Code     string     `json:"code" valid:"notnull"`
	Name     string     `json:"name" valid:"notnull"`
	Accounts []*Account `valid:"-"`
}

func (b *Bank) isValid() error {
	_, err := govalidator.ValidateStruct(b)
	return err
}

func NewBank(code string, name string) (*Bank, error) {
	bank := &Bank{
		Code: code,
		Name: name,
	}

	bank.ID = uuid.NewV4().String()
	bank.CreatedAt = time.Now()

	err := bank.isValid()
	if err != nil {
		return nil, err
	}
	return bank, nil
}
