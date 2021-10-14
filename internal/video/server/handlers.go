package server

import (
	"github.com/grigagod/compresso/internal/middleware"
	videoRepo "github.com/grigagod/compresso/internal/video/repo"

	videoUseCase "github.com/grigagod/compresso/internal/video/usecase"

	videoHttp "github.com/grigagod/compresso/internal/video/delivery/http"
)

func (s *VideoServer) MapHandlers() {
	vRepo := videoRepo.NewVideoRepository(s.Db)
	vUseCase := videoUseCase.NewVideoUseCase(&s.cfg.QueueWriteConfig, vRepo, s.Storage, s.Publisher)
	vHandlers := videoHttp.NewVideoHandlers(vUseCase)

	videoRouter := s.Router.
		With(middleware.Logger(s.Log)).
		With(middleware.JWTAuth(s.cfg.JwtSecretKey))
	videoHttp.MapVideoRoutes(videoRouter, vHandlers)
}
