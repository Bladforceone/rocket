package model

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID string
	PaymentMethod   PaymentMethod
	Status          OrderStatus
}

type Part struct {
	UUID  string
	Price float64
}
