package order

import (
	"context"
	"errors"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
	"github.com/Bladforceone/rocket/order/internal/repository/converter"
)

func (r *repo) Get(ctx context.Context, uuid string) (*modelService.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	order, ok := r.data[uuid]
	if !ok {
		return nil, errors.New("order not found")
	}

	return converter.OrderToService(order), nil
}
