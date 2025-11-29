package repository

import (
	"context"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *modelService.Order) error
	Get(ctx context.Context, uuid string) (*modelService.Order, error)
	Update(ctx context.Context, order *modelService.Order) error
}
