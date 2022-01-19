package domain

import (
	"fmt"
)

type (
	OrderID string

	OrderNumber string

	Order struct { // want "identifierを実装するフィールドは1つのみ含む必要があります。"
		id        OrderID
		subId     OrderID
		number    OrderNumber
		shipments []Shipment
		Exported  string // want "EntityはExportedなフィールドを定義することはできません。"
	}

	OrderRepository interface {
		Save(order *Order) error
	}
)

func (OrderID) ImplAsIdentifier()  {}
func (OrderID) ImplAsValueObject() {}

func (OrderNumber) ImplAsValueObject() {}

func NewOrder(id OrderID, number OrderNumber, ships []Shipment) (Order, error) {
	o := Order{
		id:        id,
		number:    number,
		shipments: ships,
	}

	if err := o.checkInvariants(); err != nil {
		return Order{}, err
	}
	return o, nil
}
func (Order) ImplAsEntity() {}

func (o Order) checkInvariants() error {
	if len(o.shipments) == 0 {
		return fmt.Errorf("配送は１つ以上必要です。")
	}
	return nil
}
