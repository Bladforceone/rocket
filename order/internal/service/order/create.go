package order

import (
	"context"
	"errors"

	"github.com/google/uuid"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
)

func (s *serv) CreateOrder(ctx context.Context, order *modelService.Order) (*modelService.Order, error) {
	// Получаем цену всех деталей
	partsPrice, err := s.inventoryClient.ListPart(ctx, order)
	if err != nil {
		return nil, err
	}

	// Если найдены не все детали, возвращаем ошибку
	if len(partsPrice) != len(order.PartUUIDs) {
		return nil, errors.New("incorrect number of parts")
	}

	// Считаем общую стоимость заказа
	var totalPrice float64
	for _, price := range partsPrice {
		totalPrice += price
	}
	order.TotalPrice = totalPrice

	// Генерируем UUID заказа
	order.OrderUUID = uuid.New().String()

	// Выставляем статус ожидания оплаты
	order.Status = modelService.OrderStatusPendingPayment

	// Сохраняем заказ
	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
