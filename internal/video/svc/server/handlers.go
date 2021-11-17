package server

import (
	sRMQ "github.com/grigagod/compresso/internal/video/svc/delivery/rmq"
	svcUseCase "github.com/grigagod/compresso/internal/video/svc/usecase"
	"github.com/grigagod/compresso/pkg/rmq"
)

func (s *Server) MapHandlers() {
	router, _ := rmq.NewRouter()
	s.router = router
	sUseCase := svcUseCase.NewSVCUseCase(s.repo, s.storage, s.logger)
	sHandlers := sRMQ.NewSVCHandlers(sUseCase)

	sRMQ.MapHandlers(s.router, sHandlers)
}
