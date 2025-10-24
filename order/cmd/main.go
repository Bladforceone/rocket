package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	order "github.com/Bladforceone/rocket/shared/pkg/openapi/order/v1"
)

func main() {
	r := chi.NewRouter()

	orderStorage := NewOrderStorage()

	orderHandler := NewOrderHandler(orderStorage)

	orderServer, err := order.NewServer(orderHandler)
	if err != nil {
		panic(err)
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", "8080"),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Starting server on :8080")
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Error starting server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) OrderHandler {
	return OrderHandler{storage: storage}
}

func (OrderHandler) CancelOrder(ctx context.Context, params order.CancelOrderParams) (order.CancelOrderRes, error) {
	return &order.CancelOrderNoContent{}, nil
}

func (OrderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (order.CreateOrderRes, error) {
	// TODO implement me
	panic("implement me")
}

func (oh OrderHandler) GetOrderByUUID(ctx context.Context, params order.GetOrderByUUIDParams) (order.GetOrderByUUIDRes, error) {
	// TODO implement me
	panic("implement me")
}

func (OrderHandler) PayOrder(ctx context.Context, req *order.PayOrderRequest, params order.PayOrderParams) (order.PayOrderRes, error) {
	// TODO implement me
	panic("implement me")
}

func (OrderHandler) NewError(ctx context.Context, err error) *order.GenericErrorStatusCode {
	// TODO implement me
	panic("implement me")
}

type OrderStorage struct {
	mutex  sync.RWMutex
	orders map[string]Order
}
type Order struct {
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float32
	TransactionUUID string
	PaymentMethod   string
	Status          string
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]Order),
	}
}

func (os *OrderStorage) Set(order Order) string {
	os.mutex.Lock()
	defer os.mutex.Unlock()

	orderUUID := uuid.New().String()

	os.orders[orderUUID] = order

	return orderUUID
}

func (os *OrderStorage) Get(orderUUID string) (Order, bool) {
	os.mutex.RLock()
	defer os.mutex.RUnlock()
	ord, exists := os.orders[orderUUID]
	return ord, exists
}
