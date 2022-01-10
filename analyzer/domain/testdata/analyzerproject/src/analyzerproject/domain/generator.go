package domain

type (
	OrderIDGenerator interface {
		GenerateOrderID() OrderID
	}

	OrderNumberGenerator interface {
		GenerateOrderNumber() OrderNumber
	}

	ShipmentCodeGenerator interface {
		GenerateShipmentCode() ShipmentCode
	}
)
