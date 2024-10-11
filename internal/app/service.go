package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/askomar/apod/internal/image-metadata/delivery"
	"github.com/askomar/apod/internal/image-metadata/dtos"
	"github.com/askomar/apod/internal/image-metadata/entities"
	imageMetadataRepository "github.com/askomar/apod/internal/image-metadata/repository"
	imageRepository "github.com/askomar/apod/internal/images/repository"
	imageutils "github.com/askomar/apod/pkg/image"
	"github.com/askomar/apod/pkg/minio"
	"github.com/go-co-op/gocron/v2"
	log "github.com/sirupsen/logrus"

	imageMetadataUsecase "github.com/askomar/apod/internal/image-metadata/usecase"
	imageUsecase "github.com/askomar/apod/internal/images/usecase"

	"github.com/askomar/apod/internal/middleware"
)

func (app *App) startService() error {

	minioClient := minio.NewMinioClient()
	err := minioClient.InitMinio(app.cfg.Minio.EndpointHost+":"+strconv.Itoa(app.cfg.Minio.EndpointPort), app.cfg.Minio.User, app.cfg.Minio.Password, app.cfg.Minio.UseSSL, app.cfg.Minio.BucketName)
	if err != nil {
		log.WithError(err).Error("Unable init minio client")
		return err
	}
	imageRepo := imageRepository.NewRepository(minioClient)
	imageStorageUC := imageUsecase.NewUseCase(imageRepo)

	imageMetadataRepo := imageMetadataRepository.NewRepository(app.db)
	imageMetadataUC := imageMetadataUsecase.NewUseCase(imageMetadataRepo)

	handlers := delivery.NewHandlers(app.cfg, imageMetadataUC, imageStorageUC)
	mux := http.NewServeMux()
	handlers.ImageRoutes(mux)
	app.server.Handler = middleware.PanicMiddleware(middleware.AccessLogMiddleware(mux))

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.WithError(err).Error("Unable create scheduler instance")
		return err
	}
	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Errorf("recovered panic in scheduler: %v", err)
					}
				}()
				response, err := http.Get(fmt.Sprintf("%s?api_key=%v", app.cfg.ProviderEndpoint, app.cfg.ProviderApiKey))
				if err != nil {
					log.WithError(err).Errorf("Unable retrieve image metadata from '%s'", app.cfg.ProviderEndpoint)
					return
				}
				metadata := dtos.ImageMetadataResponse{}
				err = json.NewDecoder(response.Body).Decode(&metadata)
				if err != nil {
					log.WithError(err).Errorf("Unable decode response from '%s'", app.cfg.ProviderEndpoint)
					return
				}
				defer response.Body.Close()

				strs := strings.Split(metadata.Url, ".")
				filename := strings.Join([]string{metadata.Date, strs[len(strs)-1]}, ".")

				image, err := imageutils.LoadImageFromURL(metadata.Url)
				if err != nil {
					log.WithError(err).Errorf("Unable retrieve image by url: %v", metadata.Url)
					return
				}
				log.Infof("Image will be saved to the storage with name '%s'", filename)
				if err := imageStorageUC.AddImage(context.Background(), filename, image); err != nil {
					log.WithError(err).Errorf("Unable save '%s' image to the storage: ", filename)
					return
				}

				parsedData, err := time.Parse("2006-01-02", metadata.Date)
				if err != nil {
					log.WithError(err).Errorf("Unable parse '%s' date", metadata.Date)
					return
				}
				if _, err := imageMetadataUC.GetImageMetadata(context.Background(), parsedData); err == nil {
					log.WithError(err).Warnf("Image metadata already exists for %v date", parsedData)
				} else {
					if err := imageMetadataUC.CreateImageMetadata(context.Background(), entities.ImageMetadata{
						Title:       metadata.Title,
						Explanation: metadata.Explanation,
						Date:        metadata.Date,
						Copyright:   metadata.Copyright,
						Filename:    filename,
					}); err != nil {
						log.WithError(err).Errorf("Unable save image metadata to the storage: %s", filename)
					}
				}
			},
		), gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	if err != nil {
		log.WithError(err).Error("Unable add job to the scheduler")
		return err
	}
	scheduler.Start()
	return nil
}
