package usecase

import (
	"context"

	"github.com/askomar/apod/internal/images"
)

type usecase struct {
	repo images.Repository
}

func NewUseCase(repo images.Repository) images.Usecase {
	return &usecase{repo: repo}
}

func (uc *usecase) GetImageURL(ctx context.Context, name string) (string, error) {
	url, err := uc.repo.GetImageURL(ctx, name)
	if err != nil {
		return "", err
	}
	return url, nil
}
func (uc *usecase) AddImage(ctx context.Context, name string, bytes []byte) error {
	err := uc.repo.AddImage(ctx, name, bytes)
	if err != nil {
		return err
	}
	return nil
}
