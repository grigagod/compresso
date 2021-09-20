package server

import (
	"log"
	_ "net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/server"
	"github.com/jmoiron/sqlx"
)

type AuthServer struct {
	authCfg *config.Auth
	*server.Server
}

func NewAuthServer(cfg *config.Config, db *sqlx.DB) *AuthServer {
	return &AuthServer{
		authCfg: &cfg.Auth,
		Server: &server.Server{
			Cfg:    cfg.Server,
			Db:     db,
			Router: mux.NewRouter(),
			Logger: log.Default(),
		},
	}
}

func (s *AuthServer) Run() {
	s.MapHandlers()
	s.Server.Start()
}
