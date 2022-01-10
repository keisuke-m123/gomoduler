package infrastructure

import "analyzerproject/domain"

type (
	OrderIDGenerator      struct{}
	OrderNumberGenerator  struct{}
	ShipmentCodeGenerator struct{}
)

func (*OrderIDGenerator) ImplAsValueObjectGenerator()      {}
func (*OrderNumberGenerator) ImplAsValueObjectGenerator()  {}
func (*ShipmentCodeGenerator) ImplAsValueObjectGenerator() {}

func (*OrderIDGenerator) GenerateOrderID() domain.OrderID {
	id := ""
	return domain.OrderID(id)
}

func (*OrderNumberGenerator) GenerateOrderNumber() domain.OrderNumber {
	n := ""
	return domain.OrderNumber(n)
}

func (*ShipmentCodeGenerator) GenerateShipmentCode() domain.ShipmentCode {
	code := ""
	return domain.ShipmentCode(code)
}
