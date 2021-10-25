package server

import (
	authRepo "github.com/grigagod/compresso/internal/auth/repo"
	"github.com/grigagod/compresso/internal/middleware"

	authUseCase "github.com/grigagod/compresso/internal/auth/usecase"

	authHttp "github.com/grigagod/compresso/internal/auth/delivery/http"
)

func (s *AuthServer) MapHandlers() {
	aRepo := authRepo.NewAuthRepository(s.DB)
	aUseCase := authUseCase.NewAuthUseCase(s.authCfg, aRepo)
	aHandlers := authHttp.NewAuthHandlers(s.authCfg, aUseCase)

	s.Router.Use(middleware.Logger(s.Log))
	authHttp.MapAuthRoutes(s.Router, aHandlers)
}
