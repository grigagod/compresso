package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/grigagod/compresso/internal/httpserver"
	"github.com/grigagod/compresso/internal/storage"
	"github.com/grigagod/compresso/internal/video/config"
	"github.com/grigagod/compresso/pkg/logger"
	"github.com/grigagod/compresso/pkg/rmq"
	"github.com/jmoiron/sqlx"
)

type VideoServer struct {
	cfg       *config.API
	Router    *chi.Mux
	DB        *sqlx.DB
	Storage   storage.Storage
	Publisher *rmq.Publisher
	Log       logger.Logger
}

func NewVideoServer(cfg *config.API, db *sqlx.DB, storage storage.Storage, publisher *rmq.Publisher, log logger.Logger) *VideoServer {
	return &VideoServer{
		cfg:       cfg,
		Router:    chi.NewMux(),
		DB:        db,
		Storage:   storage,
		Publisher: publisher,
		Log:       log,
	}
}

func (s *VideoServer) ListenAndServe(cfg *httpserver.Config) {
	httpserver.ListenAndServe(cfg, s.Router, s.Log)
}
