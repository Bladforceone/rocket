package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentAPI "github.com/Bladforceone/rocket/payment/internal/api/payment/v1"
	paymentServ "github.com/Bladforceone/rocket/payment/internal/service/payment"
	"github.com/Bladforceone/rocket/shared/pkg/proto/payment/v1"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()

	serv := paymentServ.NewService()

	api := paymentAPI.NewAPI(serv)

	paymentv1.RegisterPaymentServiceServer(s, api)
	reflection.Register(s)

	go func() {
		log.Println("Starting gRPC server on :50052")
		if err = s.Serve(lis); err != nil {
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
