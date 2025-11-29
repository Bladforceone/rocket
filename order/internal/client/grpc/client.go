package grpc

import (
	"context"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}

type InventoryClient interface {
	ListPart(ctx context.Context, order *modelService.Order) ([]float64, error)
}
