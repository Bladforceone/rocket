package v1

import (
	"context"

	generatedOrderV1 "github.com/Bladforceone/rocket/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *generatedOrderV1.PayOrderRequest, params generatedOrderV1.PayOrderParams) (generatedOrderV1.PayOrderRes, error) {
	transactionUUID, err := a.service.PayOrder(ctx, string(req.PaymentMethod.Value), params.OrderUUID)
	if err != nil {
		return nil, err
	}

	return &generatedOrderV1.PayOrderResponse{TransactionUUID: generatedOrderV1.NewOptString(transactionUUID)}, nil
}
