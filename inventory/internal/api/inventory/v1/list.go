package inventoryv1

import (
	"context"

	"github.com/Bladforceone/rocket/inventory/internal/converter"
	inventory "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventory.ListPartsRequest) (*inventory.ListPartsResponse, error) {
	list, err := a.serv.List(ctx, converter.ToServicePartFilter(req.GetFilter()))
	if err != nil {
		return nil, err
	}

	return &inventory.ListPartsResponse{Parts: converter.ToAPIParts(list)}, nil
}
