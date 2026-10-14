package application

import (
	"context"
	"fmt"
	"uuid"

	"storeapi/src/internal/domain"
	"storeapi/src/internal/ports"
)

type ImageService struct{ repo ports.ImageRepository }

func NewImageService(repo ports.ImageRepository) *ImageService { return &ImageService{repo: repo} }

func (s *ImageService) Create(ctx context.Context, image domain.Image) (domain.Image, error) {
	if len(image.Data) == 0 {
		return domain.Image{}, fmt.Errorf("%w: non-empty image is required", domain.ErrValidation)
	}
	return s.repo.Create(ctx, image)
}

func (s *ImageService) Replace(ctx context.Context, id uuid.UUID, data []byte) (domain.Image, error) {
	if len(data) == 0 {
		return domain.Image{}, fmt.Errorf("%w: image must not be empty", domain.ErrValidation)
	}
	return s.repo.Replace(ctx, id, data)
}

func (s *ImageService) ListAvailable(ctx context.Context) ([]domain.Image, error) {
	return s.repo.ListAvailable(ctx)
}

func (s *ImageService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ImageService) GetByProductID(ctx context.Context, id uuid.UUID) (domain.Image, error) {
	return s.repo.GetByProductID(ctx, id)
}

func (s *ImageService) GetByID(ctx context.Context, id uuid.UUID) (domain.Image, error) {
	return s.repo.GetByID(ctx, id)
}
