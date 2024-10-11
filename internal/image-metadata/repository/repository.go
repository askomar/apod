package repository

import (
	"context"
	"database/sql"
	"time"

	app "github.com/askomar/apod/internal/image-metadata"
	"github.com/askomar/apod/internal/image-metadata/entities"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) app.Repository {
	return &repository{db: db}
}

func (repo *repository) CreateImageMetadata(ctx context.Context, image entities.ImageMetadata) error {
	_, err := repo.db.Exec("INSERT INTO metadata (title, explanation, date, copyright, filename ) VALUES ($1, $2, $3, $4, $5)", image.Title, image.Explanation, image.Date, image.Copyright, image.Filename)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repository) GetImageMetadataByDate(ctx context.Context, date time.Time) (entities.ImageMetadata, error) {
	image := entities.ImageMetadata{}
	err := repo.db.QueryRow("SELECT title, explanation, date, copyright, filename FROM metadata WHERE date = $1", date).Scan(&image.Title, &image.Explanation, &image.Date, &image.Copyright, &image.Filename)
	if err != nil {
		return image, err
	}
	return image, nil
}

func (repo *repository) GetImagesMetadata(ctx context.Context) (images []entities.ImageMetadata, err error) {
	images = []entities.ImageMetadata{}
	rows, err := repo.db.Query("SELECT title, explanation, date, copyright, filename FROM metadata")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		image := entities.ImageMetadata{}
		err := rows.Scan(&image.Title, &image.Explanation, &image.Date, &image.Copyright, &image.Filename)
		if err != nil {
			return images, err
		}
		images = append(images, image)
	}
	return images, nil
}
