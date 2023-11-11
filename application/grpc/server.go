package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/RuanScherer/codepix-go/application/grpc/pb"
	"github.com/RuanScherer/codepix-go/application/usecase"
	"github.com/RuanScherer/codepix-go/infrastructure/repository"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPCServer(db *gorm.DB, port int) {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	pixKeyRepository := repository.PixKeyDBRepository{DB: db}
	pixUseCase := usecase.PixUseCase{PixKeyRepository: pixKeyRepository}
	pixGRPCService := NewPixGRPCService(pixUseCase)
	pb.RegisterPixServiceServer(grpcServer, pixGRPCService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}

	log.Printf("grpc server is running on port %d", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}
}
