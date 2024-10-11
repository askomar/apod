package imagemetadata

import (
	"context"
	"time"

	"github.com/askomar/apod/internal/image-metadata/entities"
)

type Usecase interface {
	CreateImageMetadata(context.Context, entities.ImageMetadata) error
	GetImageMetadata(ctx context.Context, date time.Time) (response entities.ImageMetadata, err error)
	GetImagesMetadata(ctx context.Context) (response []entities.ImageMetadata, err error)
}
