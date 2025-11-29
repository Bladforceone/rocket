package order

import (
	"context"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
	"github.com/Bladforceone/rocket/order/internal/repository/converter"
)

func (r *repo) Create(ctx context.Context, order *modelService.Order) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.data[order.OrderUUID] = converter.OrderToRepository(order)

	return nil
}
