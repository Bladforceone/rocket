package storage

import (
	"errors"
	"log"
	"sync"
)

var ErrNotFound = errors.New("order not found")

type OrderStorage struct {
	mutex  sync.RWMutex
	orders map[string]Order
}
type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID string
	PaymentMethod   string
	Status          string
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]Order),
	}
}

func (os *OrderStorage) Set(order *Order) string {
	os.mutex.Lock()
	defer os.mutex.Unlock()

	os.orders[order.OrderUUID] = *order

	return order.OrderUUID
}

func (os *OrderStorage) Get(orderUUID string) (*Order, error) {
	os.mutex.RLock()
	defer os.mutex.RUnlock()

	ord, exists := os.orders[orderUUID]

	if !exists {
		return nil, ErrNotFound
	}

	return &ord, nil
}

func (os *OrderStorage) UpdateStatus(orderUUID, status string) error {
	os.mutex.Lock()
	defer os.mutex.Unlock()

	order, exists := os.orders[orderUUID]
	if !exists {
		return ErrNotFound
	}

	order.Status = status
	os.orders[orderUUID] = order

	return nil
}

func (os *OrderStorage) Pay(orderUUID, transactionUUID, paymentMethod string) error {
	os.mutex.Lock()
	defer os.mutex.Unlock()

	log.Println(os.orders)
	order, exists := os.orders[orderUUID]
	if !exists {
		return ErrNotFound
	}

	order.TransactionUUID = transactionUUID
	order.PaymentMethod = paymentMethod
	order.Status = "PAID"
	os.orders[orderUUID] = order

	return nil
}
