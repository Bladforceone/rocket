package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	api "github.com/Bladforceone/rocket/inventory/internal/api/inventory/v1"
	partServ "github.com/Bladforceone/rocket/inventory/internal/service/part"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partRepo "github.com/Bladforceone/rocket/inventory/internal/repository/part"
	desc "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	repo := partRepo.NewRepository()
	serv := partServ.NewService(repo)
	inApi := api.NewAPI(serv)

	s := grpc.NewServer()

	desc.RegisterInventoryServiceServer(s, inApi)

	reflection.Register(s)

	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}
