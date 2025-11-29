package order

import (
	"context"
	"errors"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
	"github.com/Bladforceone/rocket/order/internal/repository/converter"
)

func (r *repo) Update(ctx context.Context, order *modelService.Order) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	updateData := converter.OrderToRepository(order)

	v, ok := r.data[updateData.OrderUUID]
	if !ok {
		return errors.New("order not found")
	}

	if updateData.TransactionUUID != nil {
		v.TransactionUUID = updateData.TransactionUUID
	}
	if updateData.PaymentMethod != nil {
		v.PaymentMethod = updateData.PaymentMethod
	}
	if updateData.Status != "" {
		v.Status = updateData.Status
	}

	r.data[order.OrderUUID] = v

	return nil
}
