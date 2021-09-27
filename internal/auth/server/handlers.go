package server

import (
	authRepo "github.com/grigagod/compresso/internal/auth/repo"

	authUseCase "github.com/grigagod/compresso/internal/auth/usecase"

	authHttp "github.com/grigagod/compresso/internal/auth/delivery/http"
)

func (s *AuthServer) MapHandlers() {
	aRepo := authRepo.NewAuthRepository(s.Db)
	aUseCase := authUseCase.NewAuthUseCase(s.authCfg, aRepo)
	aHandlers := authHttp.NewAuthHandlers(s.authCfg, aUseCase)

	authHttp.MapAuthRoutes(s.Router, aHandlers)
}
