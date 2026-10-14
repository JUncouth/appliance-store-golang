package postgres

import (
	"storeapi/src/internal/domain"
	"uuid"
)

func addressToDomain(d AddressDAO) domain.Address {
	return domain.Address{ID: d.ID, Country: d.Country, City: d.City, Street: d.Street}
}
func clientToDomain(d ClientDAO) domain.Client {
	return domain.Client{ID: d.ID, Name: d.Name, Surname: d.Surname, Birthday: d.Birthday, Gender: d.Gender, RegistrationAt: d.RegistrationAt, Address: addressToDomain(d.Address)}
}
func supplierToDomain(d SupplierDAO) domain.Supplier {
	return domain.Supplier{ID: d.ID, Name: d.Name, Address: addressToDomain(d.Address), PhoneNumber: d.PhoneNumber}
}
func productToDomain(d ProductDAO) domain.Product {
	var imageID *uuid.UUID
	if d.ImageID != nil {
		v := *d.ImageID
		imageID = &v
	}
	return domain.Product{ID: d.ID, Name: d.Name, Category: d.Category, Price: d.Price, AvailableStock: d.AvailableStock, LastUpdateAt: d.LastUpdateAt, SupplierID: d.SupplierID, ImageID: imageID}
}
func imageToDomain(d ImageDAO) domain.Image {
	return domain.Image{ID: d.ID, Data: d.Data}
}
