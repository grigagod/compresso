package main

import (
	"log"
	"os"

	authCfg "github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/auth/server"
	"github.com/grigagod/compresso/internal/config"
	"github.com/grigagod/compresso/pkg/db/postgres"
	"github.com/grigagod/compresso/pkg/logger"

	_ "github.com/grigagod/compresso/docs/auth" // load API Docs files (Swagger)
)

// @title Auth service
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email podkidysh2002@gmail.com
func main() {
	logger, err := logger.NewWrappedLogger(logger.GetLoggerConfig())
	if err != nil {
		log.Fatal(err)
	}

	var cfg authCfg.Config
	err = config.LoadConfig(&cfg, config.GetConfigPath("auth", os.Getenv("config-auth")))
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.NewPsqlDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User,
		cfg.DB.DbName, cfg.DB.Password, cfg.DB.Driver)
	if err != nil {
		logger.Fatal("Postgres connection failed:", err)
	}
	defer db.Close()

	s := server.NewAuthServer(&cfg.Auth, db, logger)
	s.MapHandlers()
	s.ListenAndServe(&cfg.HTTP)
}
