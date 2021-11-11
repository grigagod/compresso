package server

import (
	"github.com/grigagod/compresso/internal/middleware"
	"github.com/grigagod/compresso/internal/utils"
	videoRepo "github.com/grigagod/compresso/internal/video/repo"

	apiHttp "github.com/grigagod/compresso/internal/video/api/delivery/http"
	apiUseCase "github.com/grigagod/compresso/internal/video/api/usecase"

	videoPublisher "github.com/grigagod/compresso/internal/video/publisher"
)

func (s *APIServer) MapHandlers() {
	vRepo := videoRepo.NewVideoRepository(s.db)
	vPublisher := videoPublisher.NewVideoPublisher(s.cfg.QueueWriteConfig, s.channel)
	aUseCase := apiUseCase.NewAPIUseCase(vRepo, s.storage, vPublisher)
	aHandlers := apiHttp.NewAPIHandlers(aUseCase)

	apiRouter := s.router.
		With(middleware.Logger(s.log)).
		With(middleware.JWTAuth(s.cfg.JwtSecretKey)).
		With(middleware.ContentType(utils.AllowedContentTypes...))
	apiHttp.MapVideoRoutes(apiRouter, aHandlers)
}
