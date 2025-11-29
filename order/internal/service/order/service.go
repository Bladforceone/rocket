package order

import (
	grpcClient "github.com/Bladforceone/rocket/order/internal/client/grpc"
	"github.com/Bladforceone/rocket/order/internal/repository"
)

type serv struct {
	orderRepo       repository.OrderRepository
	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient
}

func NewService() *serv {
	return &serv{}
}
