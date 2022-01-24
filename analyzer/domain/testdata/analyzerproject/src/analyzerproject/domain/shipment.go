package domain

type (
	ShipmentCode string

	ShippingMethod struct {
		Name string
	}

	Shipment struct {
		ShippingMethod // want "EntityはExportedなフィールドを定義することはできません。"
		code           ShipmentCode
		address        Address
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
		Exported  string // want "ValueObjectはExportedなフィールドを定義することはできません。"
	}
)

func (ShipmentCode) ImplAsIdentifier()  {}
func (ShipmentCode) ImplAsValueObject() {}

func NewShipment(code ShipmentCode, address Address) (Shipment, error) {
	return Shipment{
		code:    code,
		address: address,
	}, nil
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

func (a *Address) SetLastName(lastname string) { // want "値オブジェクトのメソッドは値レシーバである必要があります。"
	a.lastname = lastname
}

func (a Address) SetFirstName(firstname string) {
	a.firstname = firstname
}

func (a Address) ChangeFirstname(firstname string) Address {
	o := a
	o.firstname = firstname
	return o
}
