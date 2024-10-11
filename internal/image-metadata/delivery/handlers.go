package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/askomar/apod/config"
	imageMetadata "github.com/askomar/apod/internal/image-metadata"
	"github.com/askomar/apod/internal/image-metadata/dtos"
	"github.com/askomar/apod/internal/images"
	log "github.com/sirupsen/logrus"
)

var (
	simpleDateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

type handlers struct {
	cfg config.Config
	ic  imageMetadata.Usecase
	icc images.Usecase
}

func NewHandlers(cfg config.Config, ic imageMetadata.Usecase, icc images.Usecase) *handlers {
	return &handlers{
		cfg: cfg,
		ic:  ic,
		icc: icc,
	}
}

// GetImage godoc
//
//	@Summary		Получение метаданных изображения
//	@Description	Позволяет получить метаданные изображения согласно переданной дате
//	@Tags			data
//	@Produce		json
//	@Param			date	path		string	true	"Дата в формате YYYY-MM-DD"
//	@Success		200		{object}	dtos.ImageMetadataResponse
//	@Failure		400		{object}	dtos.ErrorResponse	"Если в качестве параметра передана неправильная дата или она не в формате YYYY-MM-DD. Сообщение: invalid_date"
//	@Failure		404		{object}	dtos.ErrorResponse	"Если метаданные по переданной дате отсутствуют. Сообщение: not_found"
//	@Failure		500		{object}	dtos.ErrorResponse	"Если мы не можем получить изображение из хранилища или в случае других непредвиденных проблем"
//
// @Router	/api/apod/{date} [get]
func (h *handlers) GetImage(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("date")
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, date)

	if !simpleDateRegex.MatchString(date) || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
			Error: "invalid_date",
		})
		return
	}

	image, err := h.ic.GetImageMetadata(context.Background(), parsedDate)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
			Error: "not_found",
		})
		return
	}

	u, err := h.icc.GetImageURL(context.Background(), image.Filename)
	if err != nil {
		log.WithError(err).Errorf("Cannot get image url for %s file from storage", image.Filename)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
			Error: "internal_error",
		})
		return
	}
	uri, err := url.Parse(u)
	if err != nil {
		log.WithError(err).Errorf("Cannot parse '%s' url", u)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
			Error: "internal_error",
		})
		return
	}
	uri.Host = fmt.Sprintf("localhost:%s", strconv.Itoa(h.cfg.Minio.EndpointPort))
	image.Filename = uri.String()

	if err := json.NewEncoder(w).Encode(dtos.NewImageMetadataResponse(image)); err != nil {
		log.WithError(err).Errorf("Unable encode response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
			Error: "internal_error",
		})
		return
	}
}

// GetImages godoc
//
//	@Summary		Получение метаданных всех доступных изображений
//	@Description	Позволяет получить метаданные всех доступных изображений
//	@Tags			data
//	@Produce		json
//	@Success		200	{array}		dtos.ImageMetadataResponse
//	@Failure		500	{object}	dtos.ErrorResponse	"Если мы не можем получить изображение(я) из хранилища или в случае других непредвиденных проблем"
//
// @Router	/api/apod/ [get]
func (h *handlers) GetImages(w http.ResponseWriter, r *http.Request) {
	images, err := h.ic.GetImagesMetadata(context.Background())
	if err != nil {
		log.WithError(err).Errorf("Cannot get images metadata")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
			Error: "internal_error",
		})
		return
	}

	response := []dtos.ImageMetadataResponse{}
	for _, image := range images {
		u, err := h.icc.GetImageURL(context.Background(), image.Filename)
		if err != nil {
			log.WithError(err).Errorf("Cannot get image url from storage for '%s' file", image.Filename)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
				Error: "internal_error",
			})
			return
		}
		uri, err := url.Parse(u)
		if err != nil {
			log.WithError(err).Errorf("Cannot parse '%s' url", u)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
				Error: "internal_error",
			})
			return
		}
		uri.Host = fmt.Sprintf("localhost:%s", strconv.Itoa(h.cfg.Minio.EndpointPort))
		image.Filename = uri.String()
		response = append(response, *dtos.NewImageMetadataResponse(image))
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.WithError(err).Errorf("Unable encode response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
			Error: "internal_error",
		})
	}
}

func (h *handlers) UnimplementedRouting(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
		Error: "not_found",
	})
}

func (h *handlers) handleSwaggerFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/swagger.json")
}
