package payment

import (
	"github.com/Bladforceone/rocket/payment/internal/service"
)

type serv struct{}

func NewService() service.PaymentService {
	return &serv{}
}
