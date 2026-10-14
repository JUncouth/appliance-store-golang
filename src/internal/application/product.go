package application

import (
	"context"
	"fmt"
	"strings"
	"uuid"

	"storeapi/src/internal/domain"
	"storeapi/src/internal/ports"
)

type ProductService struct{ repo ports.ProductRepository }

func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	if strings.TrimSpace(product.Name) == "" || strings.TrimSpace(product.Category) == "" {
		return domain.Product{}, fmt.Errorf("%w: name and category are required", domain.ErrValidation)
	}
	if product.Price < 0 || product.AvailableStock < 0 {
		return domain.Product{}, fmt.Errorf("%w: price and stock must be non-negative", domain.ErrValidation)
	}
	if product.SupplierID.String() == "" {
		return domain.Product{}, fmt.Errorf("%w: supplier_id is required", domain.ErrValidation)
	}
	return s.repo.Create(ctx, product)
}

func (s *ProductService) AddProductPhoto(ctx context.Context, id uuid.UUID, imageID uuid.UUID) (domain.Product, error) {
	return s.repo.AddProductPhoto(ctx, id, imageID)
}

func (s *ProductService) DecreaseStock(ctx context.Context, id uuid.UUID, quantity int) (domain.Product, error) {
	if quantity <= 0 {
		return domain.Product{}, fmt.Errorf("%w: quantity must be greater than zero", domain.ErrValidation)
	}
	return s.repo.DecreaseStock(ctx, id, quantity)
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) ListAvailable(ctx context.Context) ([]domain.Product, error) {
	return s.repo.ListAvailable(ctx)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
