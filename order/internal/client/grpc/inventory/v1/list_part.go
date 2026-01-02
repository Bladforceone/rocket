package v1

import (
	"context"

	"github.com/Bladforceone/rocket/order/internal/client/converter"
	modelService "github.com/Bladforceone/rocket/order/internal/model"
)

func (c *client) ListPart(ctx context.Context, filter *modelService.PartFilter) ([]modelService.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, converter.FilterToListRequest(filter))
	if err != nil {
		return nil, err
	}

	return converter.PartsToService(parts), nil
}
