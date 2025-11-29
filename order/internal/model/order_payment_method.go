package model

import "fmt"

type PaymentMethod int

var mapPaymentMethodString = map[string]PaymentMethod{
	"UNKNOWN":        PaymentMethodUnknown,
	"CARD":           PaymentMethodCard,
	"SBP":            PaymentMethodSBP,
	"CREDIT_CARD":    PaymentMethodCreditCard,
	"INVESTOR_MONEY": PaymentMethodInvestorMoney,
}

var mapPaymentMethod = map[PaymentMethod]string{
	PaymentMethodUnknown:       "UNKNOWN",
	PaymentMethodCard:          "CARD",
	PaymentMethodSBP:           "SBP",
	PaymentMethodCreditCard:    "CREDIT_CARD",
	PaymentMethodInvestorMoney: "INVESTOR_MONEY",
}

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

func (s PaymentMethod) String() string {
	return mapPaymentMethod[s]
}

func ParsePaymentMethod(paymentMethod string) (PaymentMethod, error) {
	method, ok := mapPaymentMethodString[paymentMethod]
	if !ok {
		return method, fmt.Errorf("invalid payment method: %s", paymentMethod)
	}

	return method, nil
}
