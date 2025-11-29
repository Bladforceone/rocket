package model

import "fmt"

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

func (s OrderStatus) String() string {
	return string(s)
}

func ParseOrderStatus(s string) (OrderStatus, error) {
	switch s {
	case OrderStatusPendingPayment.String():
		return OrderStatusPendingPayment, nil
	case OrderStatusPaid.String():
		return OrderStatusPaid, nil
	case OrderStatusCancelled.String():
		return OrderStatusCancelled, nil
	}

	return "", fmt.Errorf("invalid payment status: %s", s)
}
