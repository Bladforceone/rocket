package order

import (
	"context"
	"fmt"

	"github.com/Bladforceone/rocket/order/internal/model"
)

func (s *serv) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	return order, nil
}
