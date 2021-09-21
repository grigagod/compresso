package server

import (
	"net/http"

	_ "net/http/pprof"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) MapDefaultHandlers() {
	// enable pprof
	s.Router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	s.Router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(s.Cfg.SwaggerUrl),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
}
