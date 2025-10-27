package part

import (
	"context"

	"github.com/Bladforceone/rocket/inventory/internal/model"
)

func (s *serv) List(ctx context.Context, filter *model.PartFilter) ([]*model.Part, error) {
	list, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return list, nil
}
