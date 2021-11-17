package main

import (
	"log"

	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/auth/server"
	"github.com/grigagod/compresso/pkg/db/postgres"
	"github.com/grigagod/compresso/pkg/logger"

	_ "github.com/grigagod/compresso/docs/auth" // load API Docs files (Swagger)
)

// @title Auth service
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email podkidysh2002@gmail.com.
func main() {
	logger, err := logger.NewWrappedLogger(logger.GetLoggerConfig())
	if err != nil {
		log.Fatal(err)
	}

	httpCfg, err := config.GetHTTPConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	dbCfg, err := postgres.GetConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	authCfg, err := config.GetAuthConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.NewPsqlDB(dbCfg.URL, dbCfg.Driver, dbCfg.MaxOpenConns,
		dbCfg.MaxIdleConns, dbCfg.ConnMaxLifetime, dbCfg.ConnMaxIdleTime)
	if err != nil {
		logger.Fatal("Postgres connection failed:", err)
	}
	defer db.Close()

	s := server.NewAuthServer(authCfg, db, logger)
	s.MapHandlers()
	s.ListenAndServe(httpCfg)
}
