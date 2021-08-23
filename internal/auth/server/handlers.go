package server

import (
	"github.com/gorilla/mux"
	authRepo "github.com/grigagod/compresso/internal/auth/repo"

	authUseCase "github.com/grigagod/compresso/internal/auth/usecase"

	authHttp "github.com/grigagod/compresso/internal/auth/delivery/http"
)

func (s *Server) MapHandlers(router *mux.Router) {
	aRepo := authRepo.NewAuthRepository(s.db)
	aUseCase := authUseCase.NewAuthUseCase(s.cfg, aRepo)
	aHandlers := authHttp.NewAuthHandlers(s.cfg, aUseCase)

	authHttp.MapAuthRoutes(router, aHandlers)

}
