package v1

import (
	def "github.com/Bladforceone/rocket/order/internal/client/grpc"
	generatedPaymentV1 "github.com/Bladforceone/rocket/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	payment generatedPaymentV1.PaymentServiceClient
}

func NewClient(paymentClient generatedPaymentV1.PaymentServiceClient) *client {
	return &client{payment: paymentClient}
}
