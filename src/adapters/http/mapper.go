package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"uuid"

	"storeapi/src/internal/domain"
)

func mapAddress(d AddressDTO) domain.Address {
	return domain.Address{Country: d.Country, City: d.City, Street: d.Street}
}

func mapClientRequest(r ClientCreateRequest) (domain.Client, error) {
	birthday, err := time.Parse("2006-01-02", r.Birthday)
	if err != nil {
		return domain.Client{}, fmt.Errorf("birthday must use YYYY-MM-DD")
	}
	reg := time.Time{}
	if r.RegistrationDate != "" {
		reg, err = time.Parse(time.RFC3339, r.RegistrationDate)
		if err != nil {
			return domain.Client{}, fmt.Errorf("registration_date must be RFC3339")
		}
	}
	return domain.Client{Name: r.ClientName, Surname: r.ClientSurname, Birthday: birthday, Gender: r.Gender, RegistrationAt: reg, Address: mapAddress(r.Address)}, nil
}

func mapClient(d domain.Client) ClientResponse {
	return ClientResponse{ID: d.ID.String(), ClientName: d.Name, ClientSurname: d.Surname, Birthday: d.Birthday.Format("2006-01-02"), Gender: d.Gender, RegistrationDate: d.RegistrationAt.UTC().Format(time.RFC3339), Address: AddressDTO{Country: d.Address.Country, City: d.Address.City, Street: d.Address.Street}}
}

func mapSupplierRequest(r SupplierCreateRequest) domain.Supplier {
	return domain.Supplier{Name: r.Name, PhoneNumber: r.PhoneNumber, Address: mapAddress(r.Address)}
}

func mapSupplier(d domain.Supplier) SupplierResponse {
	return SupplierResponse{ID: d.ID.String(), Name: d.Name, PhoneNumber: d.PhoneNumber, Address: AddressDTO{Country: d.Address.Country, City: d.Address.City, Street: d.Address.Street}}
}

func mapProductRequest(r ProductCreateRequest) (domain.Product, error) {
	supplierID, err := uuid.Parse(r.SupplierID)
	if err != nil {
		return domain.Product{}, fmt.Errorf("supplier_id: %w", err)
	}
	var imageID *uuid.UUID
	if r.ImageID != nil {
		// v, e := uuid.Parse(*r.ImageID)
		v, e := uuid.Parse(*r.ImageID)
		if e != nil {
			return domain.Product{}, fmt.Errorf("image_id: %w", e)
		}
		imageID = &v
	}
	return domain.Product{Name: r.Name, Category: r.Category, Price: r.Price, AvailableStock: r.AvailableStock, SupplierID: supplierID, ImageID: imageID}, nil
}

func mapAddPhotoRequest(d AddProductPhotoRequest) (domain.Image, error) {
	imageID, err := uuid.Parse(d.ImageID)
	if err != nil {
		return domain.Image{}, fmt.Errorf("invalid image ID")
	}
	return domain.Image{ID: imageID}, nil
}

func mapImage(d domain.Image) ImageResponse {
	return ImageResponse{ID: d.ID.String()}
}

func mapProduct(d domain.Product) ProductResponse {
	var imageID *string
	if d.ImageID != nil {
		v := d.ImageID.String()
		imageID = &v
	}
	return ProductResponse{ID: d.ID.String(), Name: d.Name, Category: d.Category, Price: d.Price, AvailableStock: d.AvailableStock, LastUpdateDate: d.LastUpdateAt.UTC().Format(time.RFC3339), SupplierID: d.SupplierID.String(), ImageID: imageID}
}

func decodeImage(raw []byte) ([]byte, error) {
	var encoded string
	if err := json.Unmarshal(raw, &encoded); err != nil {
		return nil, fmt.Errorf("image must be a base64 string")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("image is not valid base64")
	}
	return data, nil
}
