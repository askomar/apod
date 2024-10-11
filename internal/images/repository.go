package images

import (
	"context"
)

type Repository interface {
	GetImageURL(ctx context.Context, name string) (url string, err error)
	AddImage(ctx context.Context, name string, image []byte) error
}
