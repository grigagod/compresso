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
	authCfg *config.Config
	*server.Server
}

func NewAuthServer(authCfg *config.Config, db *sqlx.DB) *AuthServer {
	return &AuthServer{
		authCfg: authCfg,
		Server: &server.Server{
			Cfg:    authCfg.Server,
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
