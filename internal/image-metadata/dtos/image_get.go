package dtos

import (
	"github.com/askomar/apod/internal/image-metadata/entities"
)

type (
	ImageMetadataResponse struct {
		Title       string `json:"title"`
		Explanation string `json:"explanation"`
		Date        string `json:"date"`
		Copyright   string `json:"copyright"`
		Url         string `json:"url"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)

func NewImageMetadataResponse(image entities.ImageMetadata) *ImageMetadataResponse {
	return &ImageMetadataResponse{
		Title:       image.Title,
		Explanation: image.Explanation,
		Date:        image.Date,
		Copyright:   image.Copyright,
		Url:         image.Filename,
	}
}
