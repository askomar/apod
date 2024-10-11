package delivery

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *handlers) ImageRoutes(mux *http.ServeMux) error {
	mux.HandleFunc("GET /api/apod/{date}", h.GetImage)
	mux.HandleFunc("GET /api/apod/", h.GetImages)
	mux.HandleFunc("/", h.UnimplementedRouting)
	mux.HandleFunc("GET /specs", h.handleSwaggerFile)
	mux.HandleFunc("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/specs"),
	))
	return nil
}
