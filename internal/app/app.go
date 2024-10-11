package app

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/askomar/apod/config"
	"github.com/askomar/apod/internal/middleware"
	"github.com/askomar/apod/pkg/apperror"
	"github.com/askomar/apod/pkg/database"
	log "github.com/sirupsen/logrus"
)

type App struct {
	db     *sql.DB
	server *http.Server
	cfg    config.Config
}

func NewApp(ctx context.Context, cfg config.Config) *App {
	return &App{
		db:     apperror.Must(database.NewDatabase(cfg.Database)),
		server: middleware.NewServer(cfg),
		cfg:    cfg,
	}
}

func (app *App) Run() error {
	if err := app.startService(); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		log.Info("Shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.cfg.ShutdownServerTimeoutSec)*time.Second)
		defer cancel()

		if err := app.db.Close(); err != nil {
			log.WithError(err).Warn("Unable close database connection")
		}
		if err := app.server.Shutdown(ctx); err != nil {
			log.WithError(err).Warn("Unable shutdown http server")
		}
	}()

	return app.server.ListenAndServe()
}
