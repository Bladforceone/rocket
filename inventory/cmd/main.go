package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bladforceone/rocket/inventory/internal/repository/part"
	desc "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	_ = part.NewRepository()

	inventoryServer := &InventoryServer{}

	s := grpc.NewServer()

	desc.RegisterInventoryServiceServer(s, inventoryServer)

	reflection.Register(s)

	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}

type InventoryServer struct {
	desc.UnimplementedInventoryServiceServer
}

func (s *InventoryServer) ListParts(ctx context.Context, request *desc.ListPartsRequest) (*desc.ListPartsResponse, error) {
	//TODO
	panic("implement me")
}

func (s *InventoryServer) GetPart(ctx context.Context, request *desc.GetPartRequest) (*desc.GetPartResponse, error) {
	//TODO
	panic("implement me")
}
