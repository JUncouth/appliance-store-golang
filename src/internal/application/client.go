package application

import (
	"context"
	"fmt"
	"strings"
	"uuid"

	"storeapi/src/internal/domain"
	"storeapi/src/internal/ports"
)

type ClientService struct{ repo ports.ClientRepository }

func NewClientService(repo ports.ClientRepository) *ClientService { return &ClientService{repo: repo} }

func (s *ClientService) Create(ctx context.Context, client domain.Client) (domain.Client, error) {
	if strings.TrimSpace(client.Name) == "" || strings.TrimSpace(client.Surname) == "" {
		return domain.Client{}, fmt.Errorf("%w: name and surname are required", domain.ErrValidation)
	}
	return s.repo.Create(ctx, client)
}

func (s *ClientService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ClientService) FindByName(ctx context.Context, firstName, lastName string) ([]domain.Client, error) {
	if strings.TrimSpace(firstName) == "" || strings.TrimSpace(lastName) == "" {
		return nil, fmt.Errorf("%w: first_name and last_name are required", domain.ErrValidation)
	}
	return s.repo.FindByName(ctx, firstName, lastName)
}

func (s *ClientService) List(ctx context.Context, limit, offset *int) ([]domain.Client, error) {
	if limit != nil && (*limit <= 0 || *limit > 1000) {
		return nil, fmt.Errorf("%w: limit must be between 1 and 1000", domain.ErrValidation)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("%w: offset must be non-negative", domain.ErrValidation)
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *ClientService) ChangeAddress(ctx context.Context, id uuid.UUID, address domain.Address) (domain.Client, error) {
	if strings.TrimSpace(address.Country) == "" || strings.TrimSpace(address.City) == "" || strings.TrimSpace(address.Street) == "" {
		return domain.Client{}, fmt.Errorf("%w: country, city and street are required", domain.ErrValidation)
	}
	return s.repo.ChangeAddress(ctx, id, address)
}
