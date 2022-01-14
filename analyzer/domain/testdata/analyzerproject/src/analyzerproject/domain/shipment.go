package domain

type (
	ShipmentCode string

	Shipment struct {
		code ShipmentCode
	}

	ShipmentList struct {
		list []Shipment
	}

	Address struct {
		firstname string
		lastname  string
		state     string
		city      string
		street    string
	}
)

func (ShipmentCode) ImplAsIdentifier()  {}
func (ShipmentCode) ImplAsValueObject() {}

func NewShipment(code ShipmentCode, address Address) (Shipment, error) {
	return Shipment{code: code}, nil
}
func (Shipment) ImplAsEntity() {}

func (s Shipment) checkInvariants() error {
	return nil
}

func NewAddress(firstname, lastname, state, city, street string) Address {
	return Address{
		firstname: firstname,
		lastname:  lastname,
		state:     state,
		city:      city,
		street:    street,
	}
}
func (Address) ImplAsValueObject() {}

func (a *Address) SetLastName(lastName string) {
	a.lastname = lastName
}

func (a Address) ChangeFirstname(firstname string) Address {
	o := a
	o.firstname = firstname
	return o
}
