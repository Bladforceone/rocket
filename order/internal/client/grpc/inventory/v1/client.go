package v1

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcClient "github.com/Bladforceone/rocket/order/internal/client/grpc"
	inventory "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

type client struct {
	inventory inventory.InventoryServiceClient
}

func NewInventoryClient() (grpcClient.InventoryClient, error) {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("NewInventoryClient: %w", err)
	}

	inventoryClient := inventory.NewInventoryServiceClient(conn)

	return &client{inventory: inventoryClient}, nil
}
