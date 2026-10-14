package domain

import (
	"time"
	"uuid"
)

type Address struct {
	ID      uuid.UUID
	Country string
	City    string
	Street  string
}

type Client struct {
	ID             uuid.UUID
	Name           string
	Surname        string
	Birthday       time.Time
	Gender         string
	RegistrationAt time.Time
	Address        Address
}

type Supplier struct {
	ID          uuid.UUID
	Name        string
	Address     Address
	PhoneNumber string
}

type Product struct {
	ID             uuid.UUID
	Name           string
	Category       string
	Price          float64
	AvailableStock int
	LastUpdateAt   time.Time
	SupplierID     uuid.UUID
	ImageID        *uuid.UUID
}

type Image struct {
	ID        uuid.UUID
	Data      []byte
}
