package v1

import (
	"github.com/Bladforceone/rocket/payment/internal/service"
	payment "github.com/Bladforceone/rocket/shared/pkg/proto/payment/v1"
)

type api struct {
	payment.UnimplementedPaymentServiceServer
	service service.PaymentService
}

func NewAPI(serv service.PaymentService) payment.PaymentServiceServer {
	return &api{
		service: serv,
	}
}
