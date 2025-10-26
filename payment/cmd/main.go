package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/Bladforceone/rocket/shared/pkg/proto/payment/v1"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	paymentv1.RegisterPaymentServiceServer(s, &PaymentServer{})

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

type PaymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func (PaymentServer) PayOrder(ctx context.Context, request *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	transactionUUID := uuid.New()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID.String())

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
