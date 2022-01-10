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
	order := &domain.Order{}
	if err := o.orderRepository.Save(order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}
	return nil
}

func (o *OrderUsecase) BadPlaceOrder2() error {
	address := domain.Address{}
	ship, err := domain.NewShipment(o.shipmentCodeGenerator.GenerateShipmentCode(), address)
	if err != nil {
		return fmt.Errorf("failed to create shipment: %w", err)
	}

	number := ""
	order, err := domain.NewOrder(domain.OrderID(number), domain.OrderNumber(number), []domain.Shipment{ship})
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
