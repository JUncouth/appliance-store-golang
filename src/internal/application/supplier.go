package application

import (
	"context"
	"fmt"
	"strings"
	"uuid"

	"storeapi/src/internal/domain"
	"storeapi/src/internal/ports"
)

type SupplierService struct{ repo ports.SupplierRepository }

func NewSupplierService(repo ports.SupplierRepository) *SupplierService {
	return &SupplierService{repo: repo}
}

func (s *SupplierService) Create(ctx context.Context, supplier domain.Supplier) (domain.Supplier, error) {
	if strings.TrimSpace(supplier.Name) == "" || strings.TrimSpace(supplier.PhoneNumber) == "" {
		return domain.Supplier{}, fmt.Errorf("%w: name and phone_number are required", domain.ErrValidation)
	}
	if strings.TrimSpace(supplier.Address.Country) == "" || strings.TrimSpace(supplier.Address.City) == "" || strings.TrimSpace(supplier.Address.Street) == "" {
		return domain.Supplier{}, fmt.Errorf("%w: address is required", domain.ErrValidation)
	}
	return s.repo.Create(ctx, supplier)
}

func (s *SupplierService) ChangeAddress(ctx context.Context, id uuid.UUID, address domain.Address) (domain.Supplier, error) {
	if strings.TrimSpace(address.Country) == "" || strings.TrimSpace(address.City) == "" || strings.TrimSpace(address.Street) == "" {
		return domain.Supplier{}, fmt.Errorf("%w: country, city and street are required", domain.ErrValidation)
	}
	return s.repo.ChangeAddress(ctx, id, address)
}

func (s *SupplierService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
func (s *SupplierService) List(ctx context.Context) ([]domain.Supplier, error) {
	return s.repo.List(ctx)
}
func (s *SupplierService) GetByID(ctx context.Context, id uuid.UUID) (domain.Supplier, error) {
	return s.repo.GetByID(ctx, id)
}
