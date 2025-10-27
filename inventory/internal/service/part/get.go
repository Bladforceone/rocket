package part

import (
	"context"

	"github.com/Bladforceone/rocket/inventory/internal/model"
)

func (s *serv) Get(ctx context.Context, uuid string) (*model.Part, error) {
	data, err := s.repo.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return data, nil
}
