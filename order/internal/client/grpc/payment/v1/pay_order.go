package v1

import (
	"context"
	"fmt"

	paymentv1 "github.com/Bladforceone/rocket/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	resp, err := c.payment.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentv1.PaymentMethod(paymentv1.PaymentMethod_value[paymentMethod]),
	})
	if err != nil {
		return "", fmt.Errorf("payOrder: %w", err)
	}

	return resp.GetTransactionUuid(), nil
}
