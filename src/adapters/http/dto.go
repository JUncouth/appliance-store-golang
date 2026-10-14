package httpapi

import "encoding/json"

type AddressDTO struct {
	Country string `json:"country" example:"Russia"`
	City    string `json:"city" example:"Moscow"`
	Street  string `json:"street" example:"Leninskiy Avenue, 30"`
}

type ClientCreateRequest struct {
	ClientName       string     `json:"client_name" example:"Mike"`
	ClientSurname    string     `json:"client_surname" example:"Wazovski"`
	Birthday         string     `json:"birthday" example:"2000-01-01"`
	Gender           string     `json:"gender" example:"male"`
	RegistrationDate string     `json:"registration_date,omitempty" example:"2026-01-15T10:00:00Z"`
	Address          AddressDTO `json:"address"`
}

type ClientResponse struct {
	ID               string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ClientName       string     `json:"client_name" example:"Mike"`
	ClientSurname    string     `json:"client_surname" example:"Wazovski"`
	Birthday         string     `json:"birthday" example:"2000-01-01"`
	Gender           string     `json:"gender" example:"male"`
	RegistrationDate string     `json:"registration_date" example:"2026-01-15T10:00:00Z"`
	Address          AddressDTO `json:"address"`
}

type SupplierCreateRequest struct {
	Name        string     `json:"name" example:"DNS"`
	Address     AddressDTO `json:"address"`
	PhoneNumber string     `json:"phone_number" example:"+79959992299"`
}

type SupplierResponse struct {
	ID          string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string     `json:"name" example:"DNS"`
	Address     AddressDTO `json:"address"`
	PhoneNumber string     `json:"phone_number" example:"+79959992299"`
}

type ProductCreateRequest struct {
	Name           string  `json:"name" example:"Microwave Oven"`
	Category       string  `json:"category" example:"Kitchen Appliances"`
	Price          float64 `json:"price" example:"129.99"`
	AvailableStock int     `json:"available_stock" example:"25"`
	SupplierID     string  `json:"supplier_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ImageID        *string `json:"image_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
}

type ProductResponse struct {
	ID             string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Name           string  `json:"name" example:"Microwave Oven"`
	Category       string  `json:"category" example:"Kitchen Appliances"`
	Price          float64 `json:"price" example:"129.99"`
	AvailableStock int     `json:"available_stock" example:"25"`
	LastUpdateDate string  `json:"last_update_date" example:"2026-01-15T10:00:00Z"`
	SupplierID     string  `json:"supplier_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ImageID        *string `json:"image_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
}

type DecreaseStockRequest struct {
	Quantity int `json:"quantity" example:"2"`
}

type ImageCreateRequest struct {
	Image json.RawMessage `json:"image" swaggertype:"string" example:"iVBORw0KGgoAAAANSUhEUg..."`
}

type ImageReplaceRequest struct {
	Image json.RawMessage `json:"image" swaggertype:"string" example:"iVBORw0KGgoAAAANSUhEUg..."`
}

type AddProductPhotoRequest struct {
	ImageID string `json:"image_id" example:"550e8400-e29b-41d4-a716-446655440001"`
}

type ImageResponse struct {
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"resource not found"`
}
