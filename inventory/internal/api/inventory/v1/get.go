package inventoryv1

import (
	"context"

	"github.com/Bladforceone/rocket/inventory/internal/converter"
	inventory "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventory.GetPartRequest) (*inventory.GetPartResponse, error) {
	part, err := a.serv.Get(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventory.GetPartResponse{Part: converter.ToAPIPart(part)}, nil
}
