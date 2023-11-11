package grpc

import (
	"context"

	"github.com/RuanScherer/codepix-go/application/grpc/pb"
	"github.com/RuanScherer/codepix-go/application/usecase"
)

type PixGRPCService struct {
	PixUseCase usecase.PixUseCase
	pb.UnimplementedPixServiceServer
}

func NewPixGRPCService(usecase usecase.PixUseCase) *PixGRPCService {
	return &PixGRPCService{
		PixUseCase: usecase,
	}
}

func (p PixGRPCService) RegisterPixKey(ctx context.Context, in *pb.PixKeyRegistration) (*pb.PixKeyCreatedResult, error) {
	key, err := p.PixUseCase.RegisterKey(in.Key, in.Kind, in.AccountID)
	if err != nil {
		return &pb.PixKeyCreatedResult{
			Status: "not created",
			Error:  err.Error(),
		}, err
	}

	return &pb.PixKeyCreatedResult{
		Status: "created",
		Id:     key.ID,
	}, nil
}

func (p PixGRPCService) Find(ctx context.Context, in *pb.PixKey) (*pb.PixKeyInfo, error) {
	pixkey, err := p.PixUseCase.FindKey(in.Key, in.Kind)
	if err != nil {
		return &pb.PixKeyInfo{}, err
	}

	return &pb.PixKeyInfo{
		Id:   pixkey.ID,
		Kind: pixkey.Kind,
		Key:  pixkey.Key,
		Account: &pb.Account{
			AccountID:     pixkey.Account.ID,
			AccountNumber: pixkey.Account.Number,
			BankID:        pixkey.Account.Bank.ID,
			BankName:      pixkey.Account.Bank.Name,
			OwnerName:     pixkey.Account.OwnerName,
			CreatedAt:     pixkey.Account.CreatedAt.String(),
		},
		CreatedAt: pixkey.CreatedAt.String(),
	}, nil
}
