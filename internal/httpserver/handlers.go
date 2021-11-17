package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func MapPprofHandler(router *chi.Mux) {
	router.Mount("/debug", middleware.Profiler())
}

// MapSwaggerHandler maps swagger handler for OpenAPI spec.
func MapSwaggerHandler(router *chi.Mux, url string) {
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(url), // The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
}
