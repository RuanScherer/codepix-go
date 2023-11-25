package factory

import (
	"github.com/RuanScherer/codepix-go/application/usecase"
	"github.com/RuanScherer/codepix-go/infrastructure/repository"
	"github.com/jinzhu/gorm"
)

func TransactionUseCaseFactory(db *gorm.DB) usecase.TransactionUseCase {
	pixRepository := repository.PixKeyDBRepository{DB: db}
	transactionRepository := repository.TransactionDBRepository{DB: db}

	transactionUseCase := usecase.TransactionUseCase{
		TransactionRepository: transactionRepository,
		PixRepository:         pixRepository,
	}
	return transactionUseCase
}
