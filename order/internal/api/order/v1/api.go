package v1

import (
	"context"
	"net/http"

	"github.com/Bladforceone/rocket/order/internal/service"
	generatedOrderV1 "github.com/Bladforceone/rocket/shared/pkg/openapi/order/v1"
)

type api struct {
	service service.OrderService
	generatedOrderV1.UnimplementedHandler
}

func NewAPI(serv service.OrderService) *api {
	return &api{service: serv}
}

func (a *api) NewError(ctx context.Context, err error) *generatedOrderV1.GenericErrorStatusCode {
	return &generatedOrderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: generatedOrderV1.GenericError{
			Code:    generatedOrderV1.NewOptInt(http.StatusInternalServerError),
			Message: generatedOrderV1.NewOptString(err.Error()),
		},
	}
}
