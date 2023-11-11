package usecase

import "github.com/RuanScherer/codepix-go/domain/model"

type TransactionUseCase struct {
	TransactionRepository model.TransactionRepository
	PixRepository         model.PixKeyRepository
}

func (t *TransactionUseCase) Register(accountID string, amount float64, pixKeyto string, pixKeyKindTo string, description string) (*model.Transaction, error) {
	account, err := t.PixRepository.FindAccount(accountID)
	if err != nil {
		return nil, err
	}

	pixKey, err := t.PixRepository.FindKeyByKind(pixKeyto, pixKeyKindTo)
	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, amount, pixKey, description)
	if err != nil {
		return nil, err
	}

	t.TransactionRepository.Register(transaction)
	if transaction.Base.ID != "" {
		return transaction, nil
	}

	return nil, err
}

func (t *TransactionUseCase) Confirm(transactionID string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	err = transaction.Confirm()
	if err != nil {
		return nil, err
	}

	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionUseCase) Complete(transactionID string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	err = transaction.Complete()
	if err != nil {
		return nil, err
	}

	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionUseCase) Error(transactionID string, reason string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	err = transaction.Cancel(reason)
	if err != nil {
		return nil, err
	}

	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
