package service

import (
	"context"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
	"github.com/Bladforceone/rocket/order/internal/repository/model"
)

type OrderService interface {
	CancelOrder(ctx context.Context, orderUUID string) error
	CreateOrder(ctx context.Context, order *modelService.Order) (*modelService.Order, error)
	GetOrder(ctx context.Context, orderUUID string) (*model.Order, error)
	PayOrder(ctx context.Context, paymentMethod, orderUUID string) (string, error)
}
