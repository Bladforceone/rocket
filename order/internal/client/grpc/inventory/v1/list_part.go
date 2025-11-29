package v1

import (
	"context"

	"github.com/Bladforceone/rocket/order/internal/client/converter"
	modelService "github.com/Bladforceone/rocket/order/internal/model"
)

func (c *client) ListPart(ctx context.Context, order *modelService.Order) ([]float64, error) {
	parts, err := c.inventory.ListParts(ctx, converter.OrderToListRequest(order))
	if err != nil {
		return nil, err
	}

	return converter.PartsPriceToService(parts), nil
}
