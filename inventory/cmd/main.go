package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	desc "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	partStorage := NewPartStorage()
	FillTestData(partStorage)

	inventoryServer := &InventoryServer{
		storage: partStorage,
	}

	s := grpc.NewServer()

	desc.RegisterInventoryServiceServer(s, inventoryServer)

	reflection.Register(s)

	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}

type InventoryServer struct {
	storage *PartStorage
	desc.UnimplementedInventoryServiceServer
}

func (s *InventoryServer) ListParts(ctx context.Context, request *desc.ListPartsRequest) (*desc.ListPartsResponse, error) {
	s.storage.mtx.Lock()
	defer s.storage.mtx.Unlock()

	var result []*desc.Part
	for uuid, part := range s.storage.parts {
		if !partMatchesFilter(uuid, part, request.Filter) {
			continue
		}

		result = append(result, convertPartToDesc(uuid, part))
	}

	return &desc.ListPartsResponse{Parts: result}, nil
}

func (s *InventoryServer) GetPart(ctx context.Context, request *desc.GetPartRequest) (*desc.GetPartResponse, error) {
	s.storage.mtx.Lock()
	part := s.storage.parts[request.GetUuid()]
	s.storage.mtx.Unlock()

	metadata := convertMetadata(part.Metadata)

	return &desc.GetPartResponse{
		Part: &desc.Part{
			Uuid:          request.GetUuid(),
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      desc.Category(part.Category),
			Dimensions: &desc.Dimensions{
				Length: part.Dimensions.Length,
				Width:  part.Dimensions.Width,
				Height: part.Dimensions.Height,
				Weight: part.Dimensions.Weight,
			},
			Manufacturer: &desc.Manufacturer{
				Name:    part.Manufacturer.Name,
				Country: part.Manufacturer.Country,
				Website: part.Manufacturer.Website,
			},
			Tags:      part.Tags,
			Metadata:  metadata,
			CreatedAt: timestamppb.New(part.CreatedAt),
			UpdatedAt: timestamppb.New(part.UpdatedAt),
		},
	}, nil
}

type PartStorage struct {
	mtx   sync.RWMutex
	parts map[string]Part
}

func NewPartStorage() *PartStorage {
	return &PartStorage{
		parts: make(map[string]Part),
	}
}

type Part struct {
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Category int32

const (
	CategoryUnknown Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// FillTestData заполняет PartStorage тестовыми частями
func FillTestData(storage *PartStorage) {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()

	now := time.Now()

	parts := []Part{
		{
			Name:          "Main Engine",
			Description:   "Primary rocket engine",
			Price:         1200000.50,
			StockQuantity: 5,
			Category:      CategoryEngine,
			Dimensions:    Dimensions{Length: 3.5, Width: 2.0, Height: 2.0, Weight: 1500},
			Manufacturer:  Manufacturer{Name: "RocketCorp", Country: "USA", Website: "https://rocketcorp.com"},
			Tags:          []string{"engine", "main"},
			Metadata: map[string]interface{}{
				"power":  int64(7600),
				"active": true,
			},
			CreatedAt: now.Add(-24 * time.Hour),
			UpdatedAt: now,
		},
		{
			Name:          "Fuel Tank",
			Description:   "Liquid fuel tank",
			Price:         500000.00,
			StockQuantity: 10,
			Category:      CategoryFuel,
			Dimensions:    Dimensions{Length: 4.0, Width: 2.5, Height: 2.5, Weight: 1200},
			Manufacturer:  Manufacturer{Name: "FuelMakers", Country: "Germany", Website: "https://fuelmakers.de"},
			Tags:          []string{"fuel", "tank"},
			Metadata: map[string]interface{}{
				"capacity": int64(15000),
				"material": "aluminum",
			},
			CreatedAt: now.Add(-48 * time.Hour),
			UpdatedAt: now,
		},
		{
			Name:          "Porthole Window",
			Description:   "Window for spacecraft cabin",
			Price:         7500.25,
			StockQuantity: 20,
			Category:      CategoryPorthole,
			Dimensions:    Dimensions{Length: 0.5, Width: 0.5, Height: 0.1, Weight: 20},
			Manufacturer:  Manufacturer{Name: "SpaceGlass", Country: "Japan", Website: "https://spaceglass.jp"},
			Tags:          []string{"window", "porthole"},
			Metadata: map[string]interface{}{
				"material":  "reinforced glass",
				"radiation": true,
			},
			CreatedAt: now.Add(-72 * time.Hour),
			UpdatedAt: now,
		},
		{
			Name:          "Wing Stabilizer",
			Description:   "Stabilizing wing for spacecraft",
			Price:         45000.00,
			StockQuantity: 8,
			Category:      CategoryWing,
			Dimensions:    Dimensions{Length: 5.0, Width: 0.8, Height: 0.3, Weight: 300},
			Manufacturer:  Manufacturer{Name: "AeroParts", Country: "USA", Website: "https://aeroparts.com"},
			Tags:          []string{"wing", "stabilizer"},
			Metadata: map[string]interface{}{
				"material": "carbon fiber",
				"flexible": false,
			},
			CreatedAt: now.Add(-96 * time.Hour),
			UpdatedAt: now,
		},
	}

	for _, p := range parts {
		id := uuid.New()
		log.Printf("Adding part: %s with UUID: %s", p.Name, id.String())
		storage.parts[id.String()] = p
	}
}

func partMatchesFilter(uuid string, part Part, filter *desc.PartFilter) bool {
	if len(filter.Uuids) > 0 && !uuidInList(uuid, filter.Uuids) {
		return false
	}
	if len(filter.Names) > 0 && !nameInList(part.Name, filter.Names) {
		return false
	}
	if len(filter.Categories) > 0 && !categoryInList(desc.Category(part.Category), filter.Categories) {
		return false
	}
	if len(filter.ManufacturerCountries) > 0 && !countryInList(part.Manufacturer.Country, filter.ManufacturerCountries) {
		return false
	}
	if len(filter.Tags) > 0 && !tagsMatch(part.Tags, filter.Tags) {
		return false
	}
	return true
}

func uuidInList(uuid string, list []*wrapperspb.StringValue) bool {
	for _, u := range list {
		if u.Value == uuid {
			return true
		}
	}
	return false
}

func nameInList(name string, list []*wrapperspb.StringValue) bool {
	for _, n := range list {
		if n.Value == name {
			return true
		}
	}
	return false
}

func countryInList(country string, list []*wrapperspb.StringValue) bool {
	for _, c := range list {
		if c.Value == country {
			return true
		}
	}
	return false
}

func tagsMatch(partTags []string, filterTags []*wrapperspb.StringValue) bool {
	for _, ft := range filterTags {
		for _, pt := range partTags {
			if ft.Value == pt {
				return true
			}
		}
	}
	return false
}

func categoryInList(cat desc.Category, list []desc.Category) bool {
	for _, c := range list {
		if cat == c {
			return true
		}
	}
	return false
}

func convertPartToDesc(uuid string, part Part) *desc.Part {
	return &desc.Part{
		Uuid:          uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      desc.Category(part.Category),
		Dimensions: &desc.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &desc.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  convertMetadata(part.Metadata),
		CreatedAt: timestamppb.New(part.CreatedAt),
		UpdatedAt: timestamppb.New(part.UpdatedAt),
	}
}

func convertMetadata(src map[string]any) map[string]*desc.Value {
	result := make(map[string]*desc.Value)
	for k, v := range src {
		switch val := v.(type) {
		case string:
			result[k] = &desc.Value{Kind: &desc.Value_StringValue{StringValue: val}}
		case int64:
			result[k] = &desc.Value{Kind: &desc.Value_IntValue{IntValue: val}}
		case float64:
			result[k] = &desc.Value{Kind: &desc.Value_DoubleValue{DoubleValue: val}}
		case bool:
			result[k] = &desc.Value{Kind: &desc.Value_BoolValue{BoolValue: val}}
		default:
			log.Printf("unknown metadata type: %T", val)
		}
	}
	return result
}
