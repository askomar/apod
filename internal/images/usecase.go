package images

import (
	"context"
)

type Usecase interface {
	// get pre-signed url image
	GetImageURL(ctx context.Context, name string) (url string, err error)

	// add image to the storage with the specified name
	AddImage(ctx context.Context, name string, image []byte) error
}
