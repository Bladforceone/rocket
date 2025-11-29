package v1

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcClient "github.com/Bladforceone/rocket/order/internal/client/grpc"
	payment "github.com/Bladforceone/rocket/shared/pkg/proto/payment/v1"
)

type client struct {
	payment payment.PaymentServiceClient
}

func NewClient() (grpcClient.PaymentClient, error) {
	conn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("payment inventory constructor: %w", err)
	}

	paymentClient := payment.NewPaymentServiceClient(conn)

	return &client{payment: paymentClient}, nil
}
