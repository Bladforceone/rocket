package v1

import (
	"context"

	generatedOrderV1 "github.com/Bladforceone/rocket/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params generatedOrderV1.CancelOrderParams) (generatedOrderV1.CancelOrderRes, error) {
	err := a.service.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		return nil, err
	}

	return &generatedOrderV1.CancelOrderNoContent{}, nil
}
