package main

import (
	"log"

	"github.com/grigagod/compresso/internal/storage"
	"github.com/grigagod/compresso/internal/video/config"
	"github.com/grigagod/compresso/internal/video/server"
	"github.com/grigagod/compresso/pkg/db/aws"
	"github.com/grigagod/compresso/pkg/db/postgres"
	"github.com/grigagod/compresso/pkg/logger"
	"github.com/grigagod/compresso/pkg/rmq"
)

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

	storageCfg, err := storage.GetConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	rmqCfg, err := rmq.GetConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	apiCfg, err := config.GetAPIConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.NewPsqlDB(dbCfg.URL, dbCfg.Driver, dbCfg.MaxOpenConns,
		dbCfg.MaxIdleConns, dbCfg.ConnMaxLifetime, dbCfg.ConnMaxIdleTime)
	if err != nil {
		logger.Fatal("Postgres connection failed:", err)
	}
	defer db.Close()

	s3client, err := aws.NewClientWithSharedCredentials("./.aws/credentials", "test")
	if err != nil {
		logger.Fatal("AWS S3 session failed:", err)
	}

	storage := storage.NewAWSStorage(storageCfg, s3client)

	ch, err := rmq.NewChannelFromConfig(rmqCfg)
	if err != nil {
		logger.Fatal("RMQ connection failed:", err)
	}
	defer ch.Close()

	pub := rmq.NewPublisher(ch)

	s := server.NewVideoServer(apiCfg, db, storage, pub, logger)
	s.MapHandlers()
	s.ListenAndServe(httpCfg)
}
