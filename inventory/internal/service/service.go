package service

import (
	"context"

	"github.com/Bladforceone/rocket/inventory/internal/model"
)

type PartService interface {
	Get(ctx context.Context, uuid string) (*model.Part, error)
	List(ctx context.Context, filter *model.PartFilter) ([]*model.Part, error)
}
