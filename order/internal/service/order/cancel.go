package order

import (
	"context"
	"fmt"

	"github.com/Bladforceone/rocket/order/internal/model"
)

func (s *serv) CancelOrder(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("failed to get order: %v", err)
	}

	switch order.Status {
	// Заказ не оплачен
	case model.OrderStatusPendingPayment:
		order.Status = model.OrderStatusCancelled
		err = s.orderRepo.Update(ctx, order)
		if err != nil {
			return fmt.Errorf("failed to update order: %v", err)
		}

		return nil

	// Заказ уже оплачен и отмене не подлежит
	case model.OrderStatusPaid:
		return fmt.Errorf("order is paid")

	// Заказ уже отменён
	default:
		return nil
	}
}
