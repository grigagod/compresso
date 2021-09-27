package httpserver

import (
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func MapPprofHandler(router *chi.Mux) {
	router.HandleFunc("/pprof/", pprof.Index)
}

func MapSwaggerHandler(router *chi.Mux, url string) {
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(url), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
}
