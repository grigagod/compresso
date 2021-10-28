package server

import (
	"github.com/grigagod/compresso/internal/middleware"
	"github.com/grigagod/compresso/internal/utils"
	videoRepo "github.com/grigagod/compresso/internal/video/repo"

	videoHttp "github.com/grigagod/compresso/internal/video/delivery/http"
	videoUseCase "github.com/grigagod/compresso/internal/video/usecase"
)

func (s *VideoServer) MapHandlers() {
	vRepo := videoRepo.NewVideoRepository(s.DB)
	vUseCase := videoUseCase.NewVideoUseCase(s.cfg.QueueWriteConfig, vRepo, s.Storage, s.Publisher)
	vHandlers := videoHttp.NewVideoHandlers(vUseCase)

	videoRouter := s.Router.
		With(middleware.Logger(s.Log)).
		With(middleware.JWTAuth(s.cfg.JwtSecretKey)).
		With(middleware.ContentType(utils.AllowedContentTypes...))
	videoHttp.MapVideoRoutes(videoRouter, vHandlers)
}
