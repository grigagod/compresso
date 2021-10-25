package server

import (
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/httpserver"
	"github.com/grigagod/compresso/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type AuthServer struct {
	authCfg *config.Auth
	Router  *chi.Mux
	DB      *sqlx.DB
	Log     logger.Logger
}

func NewAuthServer(cfg *config.Auth, db *sqlx.DB, log logger.Logger) *AuthServer {
	return &AuthServer{
		authCfg: cfg,
		Router:  chi.NewMux(),
		DB:      db,
		Log:     log,
	}
}

func (s *AuthServer) ListenAndServe(cfg *httpserver.Config) {
	httpserver.ListenAndServe(cfg, s.Router, s.Log)
}
