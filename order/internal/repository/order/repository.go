package order

import (
	"sync"

	"github.com/Bladforceone/rocket/order/internal/repository"
	"github.com/Bladforceone/rocket/order/internal/repository/model"
)

type repo struct {
	mtx  sync.RWMutex
	data map[string]*model.Order
}

func NewRepository() repository.OrderRepository {
	return &repo{
		data: make(map[string]*model.Order),
	}
}
