package payment

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentv1 "github.com/Bladforceone/rocket/shared/pkg/proto/payment/v1"
)

type Client struct {
	paymentv1.PaymentServiceClient
}

func NewPaymentClient() (*Client, error) {
	conn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	paymentClient := paymentv1.NewPaymentServiceClient(conn)

	return &Client{paymentClient}, nil
}
