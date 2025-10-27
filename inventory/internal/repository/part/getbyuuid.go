package part

import (
	"context"

	"github.com/Bladforceone/rocket/inventory/internal/model"
	"github.com/Bladforceone/rocket/inventory/internal/repository/converter"
)

func (r *repository) GetByUUID(ctx context.Context, uuid string) (*model.Part, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	v, ok := r.data[uuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}

	return converter.ToServicePart(&v), nil
}
