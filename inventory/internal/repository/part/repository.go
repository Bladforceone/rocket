package part

import (
	"sync"

	repo "github.com/Bladforceone/rocket/inventory/internal/repository"
	"github.com/Bladforceone/rocket/inventory/internal/repository/model"
)

type repository struct {
	mtx  sync.RWMutex
	data map[string]model.Part
}

func NewRepository() repo.InventoryRepository {
	return &repository{data: make(map[string]model.Part)}
}
