package converter

import (
	"github.com/Bladforceone/rocket/order/internal/model"
	generatedOrderV1 "github.com/Bladforceone/rocket/shared/pkg/openapi/order/v1"
)

func OrderToGetOrderByUUIDRes(o *model.Order) *generatedOrderV1.GetOrderResponse {
	return &generatedOrderV1.GetOrderResponse{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUuids:       o.PartUUIDs,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: generatedOrderV1.NewOptString(o.TransactionUUID),
		PaymentMethod:   generatedOrderV1.NewOptPaymentMethod(generatedOrderV1.PaymentMethod(o.PaymentMethod.String())),
		Status:          generatedOrderV1.OrderStatus(o.Status.String()),
	}
}
