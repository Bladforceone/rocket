package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Bladforceone/rocket/order/internal/inventory"
	"github.com/Bladforceone/rocket/order/internal/payment"
	"github.com/Bladforceone/rocket/order/internal/storage"
	orderdesc "github.com/Bladforceone/rocket/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Bladforceone/rocket/shared/pkg/proto/payment/v1"
)

const (
	host        = "0.0.0.0"
	port        = "8080"
	readTimeout = 5 * time.Second
)

func main() {
	r := chi.NewRouter()

	orderStorage := storage.NewOrderStorage()

	orderInventoryClient, err := inventory.NewInventoryClient()
	if err != nil {
		panic(err)
	}

	orderPaymentClient, err := payment.NewPaymentClient()
	if err != nil {
		panic(err)
	}

	orderHandler := NewOrderHandler(orderStorage, orderInventoryClient, orderPaymentClient)

	orderServer, err := orderdesc.NewServer(orderHandler)
	if err != nil {
		panic(err)
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(20 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           r,
		ReadHeaderTimeout: readTimeout,
	}

	go func() {
		log.Printf("Starting server on %s:%s\n", host, port)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Error starting server: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")
}

type OrderHandler struct {
	storage   *storage.OrderStorage
	inventory *inventory.Client
	payment   *payment.Client
}

func NewOrderHandler(storage *storage.OrderStorage, inventory *inventory.Client, payment *payment.Client) OrderHandler {
	return OrderHandler{
		storage:   storage,
		inventory: inventory,
		payment:   payment,
	}
}

func (oh OrderHandler) CancelOrder(ctx context.Context, params orderdesc.CancelOrderParams) (orderdesc.CancelOrderRes, error) {
	orderData, err := oh.storage.Get(params.OrderUUID)
	if err != nil {
		return &orderdesc.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	if orderData.Status == "PAID" {
		return &orderdesc.ConflictError{
			Code:    409,
			Message: "The order has already been paid for and cannot be cancelled",
		}, nil
	}

	orderData.Status = "CANCELLED"

	_ = oh.storage.Set(orderData)

	return &orderdesc.CancelOrderNoContent{}, nil
}

func (oh OrderHandler) CreateOrder(ctx context.Context, req *orderdesc.CreateOrderRequest) (orderdesc.CreateOrderRes, error) {
	countPart := len(req.PartUuids)
	uuids := make([]*wrapperspb.StringValue, 0, countPart)
	for _, partUUID := range req.PartUuids {
		uuids = append(uuids, &wrapperspb.StringValue{Value: partUUID})
	}

	part, err := oh.inventory.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartFilter{
			Uuids: uuids,
		},
	})
	if err != nil {
		return &orderdesc.BadGatewayError{
			Code:    http.StatusBadGateway,
			Message: "error fetching parts from inventory service",
		}, err
	}

	if len(part.Parts) != countPart {
		return &orderdesc.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "one or more parts not found",
		}, nil
	}

	var totalPrice float64
	partUUIDs := make([]string, 0, countPart)
	for _, p := range part.Parts {
		totalPrice += p.Price
		partUUIDs = append(partUUIDs, p.Uuid)
	}

	orderUUID := oh.storage.Set(&storage.Order{
		OrderUUID:  uuid.New().String(),
		UserUUID:   req.UserUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     "PENDING_PAYMENT",
	})

	return &orderdesc.CreateOrderResponse{
		OrderUUID:  orderdesc.NewOptString(orderUUID),
		TotalPrice: orderdesc.NewOptFloat64(totalPrice),
	}, nil
}

func (oh OrderHandler) GetOrderByUUID(ctx context.Context, params orderdesc.GetOrderByUUIDParams) (orderdesc.GetOrderByUUIDRes, error) {
	orderData, err := oh.storage.Get(params.OrderUUID)
	if err != nil {
		return &orderdesc.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "orderdesc with UUID '" + params.OrderUUID + "' not found",
		}, nil
	}

	transactionUUID := orderdesc.OptString{}
	paymentMethod := orderdesc.OptPaymentMethod{}

	if orderData.Status != "PENDING_PAYMENT" {
		transactionUUID.Value = orderData.TransactionUUID
		transactionUUID.Set = true

		paymentMethod.Value = orderdesc.PaymentMethod(orderData.PaymentMethod)
		paymentMethod.Set = true
	}

	return &orderdesc.GetOrderResponse{
		OrderUUID:       params.OrderUUID,
		UserUUID:        orderData.UserUUID,
		PartUuids:       nil,
		TotalPrice:      orderData.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderdesc.OrderStatus(orderData.Status),
	}, nil
}

func (oh OrderHandler) PayOrder(ctx context.Context, req *orderdesc.PayOrderRequest, params orderdesc.PayOrderParams) (orderdesc.PayOrderRes, error) {
	payMethod := req.GetPaymentMethod().Value
	orderUUID := params.OrderUUID
	log.Print(string(payMethod), "\n")
	log.Printf("orderUUID: %s", orderUUID)
	payOrder, err := oh.payment.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: paymentv1.PaymentMethod(paymentv1.PaymentMethod_value[string(payMethod)]),
	})
	if err != nil {
		log.Print("tut")
		return nil, err
	}

	err = oh.storage.Pay(orderUUID, payOrder.TransactionUuid, string(payMethod))
	if err != nil {
		log.Print("tam")
		return nil, err
	}

	return &orderdesc.PayOrderResponse{TransactionUUID: orderdesc.NewOptString(orderUUID)}, nil
}

func (OrderHandler) NewError(ctx context.Context, err error) *orderdesc.GenericErrorStatusCode {
	return &orderdesc.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderdesc.GenericError{
			Code:    orderdesc.NewOptInt(http.StatusInternalServerError),
			Message: orderdesc.NewOptString(err.Error()),
		},
	}
}
