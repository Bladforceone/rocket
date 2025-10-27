package inventoryv1

import (
	"github.com/Bladforceone/rocket/inventory/internal/service"
	inventory "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

type api struct {
	serv service.PartService
	inventory.UnimplementedInventoryServiceServer
}

func NewAPI(serv service.PartService) inventory.InventoryServiceServer {
	return &api{
		serv: serv,
	}
}
