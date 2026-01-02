package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
)

func (s *serv) CreateOrder(ctx context.Context, order *modelService.Order) (*modelService.Order, error) {
	filter := &modelService.PartFilter{
		UUIDs: order.PartUUIDs,
	}
	// Получаем все детали
	parts, err := s.inventoryClient.ListPart(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Если найдены не все детали, возвращаем ошибку
	if len(parts) != len(order.PartUUIDs) {
		missingParts := foundMissingPart(order.PartUUIDs, parts)
		return nil, fmt.Errorf("missing parts: %v", missingParts)
	}

	// Считаем общую стоимость заказа
	var totalPrice float64
	for _, p := range parts {
		totalPrice += p.Price
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

// foundMissingPart Ищем детали которые небыли найдены
func foundMissingPart(orderUUIDs []string, parts []modelService.Part) []string {
	missing := make([]string, 0, len(orderUUIDs))
	partUUIDMap := make(map[string]bool)

	// Заполняем map UUID'ами найденных деталей
	for _, part := range parts {
		partUUIDMap[part.UUID] = true
	}

	// Проверяем каждый UUID из orderUUIDs на наличие в partUUIDMap
	for _, uuid := range orderUUIDs {
		if !partUUIDMap[uuid] {
			missing = append(missing, uuid)
		}
	}

	return missing
}
