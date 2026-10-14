package postgres

import (
	"time"
	"uuid"
)

type AddressDAO struct {
	ID      uuid.UUID
	Country string
	City    string
	Street  string
}

type ClientDAO struct {
	ID             uuid.UUID
	Name           string
	Surname        string
	Birthday       time.Time
	Gender         string
	RegistrationAt time.Time
	Address        AddressDAO
}

type SupplierDAO struct {
	ID          uuid.UUID
	Name        string
	Address     AddressDAO
	PhoneNumber string
}

type ProductDAO struct {
	ID             uuid.UUID
	Name           string
	Category       string
	Price          float64
	AvailableStock int
	LastUpdateAt   time.Time
	SupplierID     uuid.UUID
	ImageID        *uuid.UUID
}

type ImageDAO struct {
	ID        uuid.UUID
	Data      []byte
}
