package infrastructure

import "analyzerproject/domain"

type (
	OrderIDGenerator      struct{}
	OrderNumberGenerator  struct{}
	ShipmentCodeGenerator struct{}
)

func (*OrderNumberGenerator) ImplAsValueObjectGenerator()  {}
func (*ShipmentCodeGenerator) ImplAsValueObjectGenerator() {}

func (*OrderIDGenerator) GenerateOrderID() domain.OrderID {
	id := ""
	return domain.OrderID(id) // want "ValueObjectを実装した構造体はValueObjectが存在するパッケージ以外から直接生成することはできません。"
}

func (*OrderNumberGenerator) GenerateOrderNumber() domain.OrderNumber {
	n := ""
	return domain.OrderNumber(n)
}

func (*ShipmentCodeGenerator) GenerateShipmentCode() domain.ShipmentCode {
	code := ""
	return domain.ShipmentCode(code)
}
