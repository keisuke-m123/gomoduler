package usecase

import (
	"analyzerproject/domain"
	"fmt"
)

type OrderUsecase struct {
	orderRepository       domain.OrderRepository
	orderIDGenerator      domain.OrderIDGenerator
	orderNumberGenerator  domain.OrderNumberGenerator
	shipmentCodeGenerator domain.ShipmentCodeGenerator
}

func (o *OrderUsecase) BadPlaceOrder() error {
	order := &domain.Order{} // want "Entityを実装した構造体はEntityが存在するパッケージ以外からcomposite literalで生成することはできません。"
	if err := o.orderRepository.Save(order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}
	return nil
}

func (o *OrderUsecase) BadPlaceOrder2() error {
	address := domain.Address{} // want "ValueObjectを実装した構造体はValueObjectが存在するパッケージ以外からcomposite literalで生成することはできません。"
	ship, err := domain.NewShipment(o.shipmentCodeGenerator.GenerateShipmentCode(), address)
	if err != nil {
		return fmt.Errorf("failed to create shipment: %w", err)
	}

	number := ""
	orderID := domain.OrderID(number)         // want "ValueObjectを実装した構造体はValueObjectが存在するパッケージ以外から直接生成することはできません。"
	orderNumber := domain.OrderNumber(number) // want "ValueObjectを実装した構造体はValueObjectが存在するパッケージ以外から直接生成することはできません。"
	order, err := domain.NewOrder(orderID, orderNumber, []domain.Shipment{ship})
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	if err := o.orderRepository.Save(&order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}
	return nil
}

func (o *OrderUsecase) PlaceOrder() error {
	shipment, err := domain.NewShipment(o.shipmentCodeGenerator.GenerateShipmentCode(), domain.NewAddress(
		"firstname",
		"lastname",
		"state",
		"city",
		"zipcode",
	))
	if err != nil {
		return fmt.Errorf("failed to create shipment: %w", err)
	}

	order, err := domain.NewOrder(
		o.orderIDGenerator.GenerateOrderID(),
		o.orderNumberGenerator.GenerateOrderNumber(),
		[]domain.Shipment{shipment},
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	if err := o.orderRepository.Save(&order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	return nil
}
