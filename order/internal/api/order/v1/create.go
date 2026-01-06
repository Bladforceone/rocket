package v1

import (
	"context"

	"github.com/Bladforceone/rocket/order/internal/model"
	generatedOrderV1 "github.com/Bladforceone/rocket/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *generatedOrderV1.CreateOrderRequest) (generatedOrderV1.CreateOrderRes, error) {
	order, err := a.service.CreateOrder(ctx, &model.Order{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	})
	if err != nil {
		return nil, err
	}

	return &generatedOrderV1.CreateOrderResponse{
			OrderUUID:  generatedOrderV1.NewOptString(order.OrderUUID),
			TotalPrice: generatedOrderV1.NewOptFloat64(order.TotalPrice)},
		nil
}
