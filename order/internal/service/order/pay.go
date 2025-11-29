package order

import (
	"context"
	"fmt"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
)

func (s *serv) PayOrder(ctx context.Context, paymentMethod, orderUUID string) (string, error) {
	// Находим заказ
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return "", fmt.Errorf("failed to get order: %v", err)
	}

	// Проверяем метод
	method, err := modelService.ParsePaymentMethod(paymentMethod)
	if err != nil {
		return "", fmt.Errorf("failed to parse payment method: %v", err)
	}

	// Оплачиваем заказ
	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.OrderUUID, order.UserUUID, method.String())
	if err != nil {
		return "", fmt.Errorf("failed to pay order: %v", err)
	}

	// Обновляем статус заказа
	order.Status = modelService.OrderStatusPaid
	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return "", fmt.Errorf("failed to update order: %v", err)
	}

	return transactionUUID, nil
}
