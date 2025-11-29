package v1

import (
	"context"

	payment "github.com/Bladforceone/rocket/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(context.Context, *payment.PayOrderRequest) (*payment.PayOrderResponse, error) {
	transactionUUID := a.service.Pay()
	return &payment.PayOrderResponse{TransactionUuid: transactionUUID}, nil
}
