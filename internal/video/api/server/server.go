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

type APIServer struct {
	cfg     *config.API
	router  *chi.Mux
	db      *sqlx.DB
	storage storage.Storage
	channel *rmq.Channel
	log     logger.Logger
}

func NewAPIServer(cfg *config.API, db *sqlx.DB, storage storage.Storage, ch *rmq.Channel, log logger.Logger) *APIServer {
	return &APIServer{
		cfg:     cfg,
		router:  chi.NewMux(),
		db:      db,
		storage: storage,
		channel: ch,
		log:     log,
	}
}

func (s *APIServer) ListenAndServe(cfg *httpserver.Config) {
	httpserver.ListenAndServe(cfg, s.router, s.log)
}
