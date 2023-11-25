package model

import (
	"encoding/json"

	validator "github.com/asaskevich/govalidator"
)

type Transaction struct {
	ID               string  `json:"id" validate:"required,uuid4"`
	AccountID        string  `json:"accountId" validate:"required,uuid4"`
	Amount           float64 `json:"amount" validate:"required"`
	PixKeyTo         string  `json:"pixKeyTo" validate:"required"`
	PixKeyKindTo     string  `json:"pixKeyKindTo" validate:"required"`
	Description      string  `json:"description" validate:"required"`
	Status           string  `json:"status" validate:"required"`
	ErrorDescription string  `json:"errorDescription"`
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (t *Transaction) isValid() error {
	_, err := validator.ValidateStruct(t)
	return err
}

func (t *Transaction) ParseJSON(data []byte) error {
	err := json.Unmarshal(data, t)
	if err != nil {
		return err
	}

	err = t.isValid()
	if err != nil {
		return err
	}

	return nil
}

func (t *Transaction) ToJSON() ([]byte, error) {
	err := t.isValid()
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return result, nil
}
