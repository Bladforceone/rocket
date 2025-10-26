package inventory

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryv1 "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

type Client struct {
	inventoryv1.InventoryServiceClient
}

func NewInventoryClient() (*Client, error) {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := inventoryv1.NewInventoryServiceClient(conn)

	return &Client{client}, nil
}
