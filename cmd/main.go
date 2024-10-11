package main

import (
	"context"
	"os"

	"github.com/askomar/apod/config"
	_ "github.com/askomar/apod/docs"
	"github.com/askomar/apod/internal/app"
	log "github.com/sirupsen/logrus"
)

// @title		APOD API
// @version		1.0
// @description	Сервис предоставляет возможность получения "a picture of the day" и его метаданные
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("Unable load app configuration")
	}

	if err := app.NewApp(context.Background(), cfg).Run(); err != nil {
		os.Exit(1)
	}
}
