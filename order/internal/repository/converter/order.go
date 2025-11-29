package converter

import (
	modelService "github.com/Bladforceone/rocket/order/internal/model"
	"github.com/Bladforceone/rocket/order/internal/repository/model"
)

func OrderToService(order *model.Order) *modelService.Order {
	return &modelService.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: *order.TransactionUUID,
		PaymentMethod:   convertPaymentMethod(order),
		Status:          modelService.OrderStatus(order.Status),
	}
}

// ASK Нормальная практика доверять БД?
func convertPaymentMethod(order *model.Order) modelService.PaymentMethod {
	// Игнорируем ошибку (Доверяем БД)
	method, _ := modelService.ParsePaymentMethod(*order.PaymentMethod)
	return method
}

func OrderToRepository(order *modelService.Order) *model.Order {
	paymentMethod := order.PaymentMethod.String()
	return &model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: &order.TransactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          order.Status.String(),
	}
}
