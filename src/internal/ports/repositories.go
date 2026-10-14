package ports

import (
	"context"
	"uuid"

	"storeapi/src/internal/domain"
)

type ClientRepository interface {
	Create(ctx context.Context, client domain.Client) (domain.Client, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByName(ctx context.Context, firstName, lastName string) ([]domain.Client, error)
	List(ctx context.Context, limit, offset *int) ([]domain.Client, error)
	ChangeAddress(ctx context.Context, id uuid.UUID, address domain.Address) (domain.Client, error)
}

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) (domain.Product, error)
	AddProductPhoto(ctx context.Context, id uuid.UUID, imageID uuid.UUID) (domain.Product, error)
	DecreaseStock(ctx context.Context, id uuid.UUID, quantity int) (domain.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	ListAvailable(ctx context.Context) ([]domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SupplierRepository interface {
	Create(ctx context.Context, supplier domain.Supplier) (domain.Supplier, error)
	ChangeAddress(ctx context.Context, id uuid.UUID, address domain.Address) (domain.Supplier, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]domain.Supplier, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Supplier, error)
}

type ImageRepository interface {
	Create(ctx context.Context, image domain.Image) (domain.Image, error)
	Replace(ctx context.Context, id uuid.UUID, data []byte) (domain.Image, error)
	ListAvailable(ctx context.Context) ([]domain.Image, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByProductID(ctx context.Context, productID uuid.UUID) (domain.Image, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Image, error)
}
