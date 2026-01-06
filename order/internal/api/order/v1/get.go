package v1

import (
	"context"

	"github.com/Bladforceone/rocket/order/internal/converter"
	generatedOrderV1 "github.com/Bladforceone/rocket/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params generatedOrderV1.GetOrderByUUIDParams) (generatedOrderV1.GetOrderByUUIDRes, error) {
	order, err := a.service.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		return nil, err
	}

	return converter.OrderToGetOrderByUUIDRes(order), nil
}
