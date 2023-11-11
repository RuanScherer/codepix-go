package model

import (
	"time"

	"github.com/asaskevich/govalidator"

	uuid "github.com/satori/go.uuid"
)

type Bank struct {
	Base     `valid:"required"`
	Code     string     `json:"code" gorm:"type:varchar(20)" valid:"notnull"`
	Name     string     `json:"name" gorm:"type:varchar(255)" valid:"notnull"`
	Accounts []*Account `gorm:"ForeignKey:BankID" valid:"-"`
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
