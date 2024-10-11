package usecase

import (
	"context"
	"time"

	imagemetadata "github.com/askomar/apod/internal/image-metadata"
	"github.com/askomar/apod/internal/image-metadata/entities"
)

type usecase struct {
	repo imagemetadata.Repository
}

func NewUseCase(repo imagemetadata.Repository) imagemetadata.Usecase {
	return &usecase{repo: repo}
}

func (uc *usecase) CreateImageMetadata(ctx context.Context, metadata entities.ImageMetadata) error {
	return uc.repo.CreateImageMetadata(ctx, metadata)
}

func (uc *usecase) GetImageMetadata(ctx context.Context, date time.Time) (entities.ImageMetadata, error) {
	image, err := uc.repo.GetImageMetadataByDate(ctx, date)
	if err != nil {
		return image, err
	}
	return image, nil
}
func (uc *usecase) GetImagesMetadata(ctx context.Context) ([]entities.ImageMetadata, error) {
	metadata, err := uc.repo.GetImagesMetadata(ctx)
	if err != nil {
		return metadata, err
	}
	return metadata, nil
}
