package imagemetadata

import (
	"context"
	"time"

	"github.com/askomar/apod/internal/image-metadata/entities"
)

type Repository interface {
	CreateImageMetadata(context.Context, entities.ImageMetadata) error
	GetImageMetadataByDate(context.Context, time.Time) (entities.ImageMetadata, error)
	GetImagesMetadata(context.Context) ([]entities.ImageMetadata, error)
}
