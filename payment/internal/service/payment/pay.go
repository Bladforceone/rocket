package payment

import (
	"log"

	"github.com/google/uuid"
)

func (s *serv) Pay() string {
	transactionUUID := uuid.New()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID.String())

	return transactionUUID.String()
}
