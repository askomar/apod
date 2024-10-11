package repository

import (
	"context"

	app "github.com/askomar/apod/internal/images"
	"github.com/askomar/apod/pkg/minio"
)

type repository struct {
	client minio.Client
}

func NewRepository(client minio.Client) app.Repository {
	return &repository{client: client}
}

func (repo *repository) GetImageURL(ctx context.Context, name string) (string, error) {
	url, err := repo.client.GetOne(name)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (repo *repository) AddImage(ctx context.Context, name string, image []byte) error {
	_, err := repo.client.CreateOne(minio.FileDataType{
		FileName: name,
		Data:     image,
	})
	return err
}
